// Code generated by gorpcgen; DO NOT EDIT.

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	coreErrors "github.com/andrskom/gorpcgen/protocol/errors"
	coreModels "github.com/andrskom/gorpcgen/protocol/models"
	coreCtx "github.com/andrskom/gorpcgen/protocol/context"
	"github.com/sirupsen/logrus"
	"github.com/nats-io/go-nats"
	serverHandlers "github.com/andrskom/gorpcgen/testservice/handlers"
)

type HandlersInterface interface {
	Get(coreModels.RequestMeta, *serverHandlers.GetRequest) (serverHandlers.GetResponse, *coreErrors.Error)
}

type MetricsInterface interface {
	APIDurationFunc(string, string, func() (respType coreModels.ResponseType))
}

type Server struct {
	logger          *logrus.Entry
	nc              *nats.Conn
	ctx             context.Context
	cancelFunc      context.CancelFunc
	shutdownTimeout time.Duration
	shutdownChan    chan bool
	runningWG       sync.WaitGroup
	handlers        *serverHandlers.Test
	metrics         MetricsInterface
	
}

func NewServer(
	ctx context.Context,
	logger *logrus.Entry,
	nc *nats.Conn,
	handlers *serverHandlers.Test,
	metrics MetricsInterface,
	
) *Server {
	s := &Server{
		logger:          logger,
		nc:              nc,
		runningWG:       sync.WaitGroup{},
		shutdownTimeout: 3 * time.Second,
		shutdownChan:    make(chan bool),
		handlers:        handlers,
		metrics:         metrics,
		
	}
	s.ctx, s.cancelFunc = context.WithCancel(ctx)
	return s
}

// Listen run listening of nats server
// nolint:gocyclo
func (s *Server) Serve() error {
	defer close(s.shutdownChan)

	var wg sync.WaitGroup
	// Run listener Test.Get
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.listenTopic("Test", "Get", s.GetCall)
	}()
	wg.Wait()
	if err := s.nc.LastError(); err != nil {
		s.logger.WithFields(logrus.Fields{"err": err}).Error("Last error of nats connection")
		return err
	}
	s.runningWG.Wait()
	return nil
}

func (s *Server) GetCall(meta *coreModels.RequestMeta, data []byte) *coreModels.Response {
	var req serverHandlers.GetRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return coreModels.NewResponseErr("Unknown", coreErrors.New(coreErrors.CodeBadRequest, "Can't unmarshal request data", err))
	}
	l := s.logger.WithField("requestID", meta.RequestID).WithField("origin", meta.Origin)
	result, errModel := s.handlers.Get(coreCtx.NewContext(meta, l), &req)
	if errModel != nil {
		return coreModels.NewResponseErr(meta.RequestID, errModel)
	}
	respBytes, err := json.Marshal(result)
	if err != nil {
		return coreModels.NewResponseErr(meta.RequestID, coreErrors.New(coreErrors.CodeInternal, "Can't marshal response", err))
	}
	return coreModels.NewResponseOK(meta.RequestID, respBytes)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.cancelFunc()
	select {
	case <-time.After(s.shutdownTimeout):
		return fmt.Errorf("can't graceful shutdown, after %s", s.shutdownTimeout.String())
	case <-s.shutdownChan:
		s.logger.Info("All calls was correct finished")
	}
	return nil
}

// Reload operation implementation
func (s *Server) Reload() error {
	return fmt.Errorf("method reload not implemented")
}

// Maintenance operation implementation
func (s *Server) Maintenance() error {
	return fmt.Errorf("method reload not implemented")
}

func (s *Server) listenTopic(service string, method string, topicFunc func(meta *coreModels.RequestMeta, data []byte) *coreModels.Response) {
	topicName := fmt.Sprintf("%s.%s", service, method)
	listenChan := make(chan *nats.Msg, 1000)
	subscr, err := s.nc.QueueSubscribeSyncWithChan(
		topicName,
		topicName,
		listenChan,
	)
	if err != nil {
		s.logger.WithField("err", err).Error("Can't subscribe to subject")
	}
	s.logger.Debugf("Listen subject: '%s'", subscr.Subject)
	if err := s.nc.Flush(); err != nil {
		s.logger.WithFields(logrus.Fields{"err": err}).Error("Can't flush nats connection")
	}
	for {
		select {
		case msg := <-listenChan:
			s.runningWG.Add(1)
			go func(msg *nats.Msg) {
				s.logger.
					WithField("topic", topicName).
					WithField("data", string(msg.Data)).
					Debug("Request")
				defer func() {
					if r := recover(); r != nil {
						s.logger.WithField("topic", topicName).Errorf(
							"Panic in call for message: %+v\n panic: %+v",
							msg,
							r,
						)
					}
					s.runningWG.Done()
				}()
				s.metrics.APIDurationFunc(service, method, func() (respType coreModels.ResponseType) {
					
					var req coreModels.Request
					if err := json.Unmarshal(msg.Data, &req); err != nil {
						return s.sendResponse(
							msg.Reply,
							topicName,
							coreModels.ResponseTypeErr,
							coreModels.NewResponseErr(
								"Unknown",
								coreErrors.New(coreErrors.CodeBadRequest, "Can't unmarshal request data", err),
							),
						)
					}
					resp := topicFunc(&req.Meta, req.Data)
					return s.sendResponse(
						msg.Reply,
						topicName,
						resp.Type,
						resp,
					)
				})
			}(msg)
		case <-s.ctx.Done():
			err := subscr.Unsubscribe()
			if err != nil {
				s.logger.WithField("err", err).Error("Can't unsubscribe from nats")
			}
			return
		}
	}
}

func (s *Server) sendResponse(
	replyTopic string,
	topicName string,
	responseType coreModels.ResponseType,
	response *coreModels.Response,
)coreModels.ResponseType {
	bytes, err := json.Marshal(response)
	if err != nil {
		bytes = coreModels.NewJSONMarshallErrResponseBytes(response.Meta.RequestID)
	}

	

	if err := s.nc.Publish(replyTopic, bytes); err != nil {
		s.logger.
			WithField("topic", topicName).
			WithField("data", string(bytes)).
			WithField("err", err).
			Errorf("Can't publish response to topic %s", replyTopic)
		return coreModels.ResponseTypeErr
	}
	s.logger.
		WithField("topic", topicName).
		WithField("data", string(bytes)).
		Debug("Response")
	return responseType
}

//======================================================================================================================

