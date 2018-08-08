package main

import (
	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/http/gen/server"
	"github.com/sirupsen/logrus"
	"net/http"
)

const DefaultHTTPAddr = ":4020"

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	s := server.New(
		logrus.WithField("server", "http"),
		&handlers.Test{},
		&http.Server{Addr: DefaultHTTPAddr},
		&metrics{},
	)
	s.Serve()
}

type metrics struct {
}

func (m *metrics) APIDurationFunc(service string, method string, do func() (respType models.ResponseType)) {
	do()
}
