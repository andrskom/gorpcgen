package main

import (
	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/tcp/gen/server"
	"github.com/sirupsen/logrus"
)

func main() {
	s := server.NewServer(
		logrus.WithField("component", "tcpServer"),
		"localhost",
		"8081",
		&handlers.Test{},
		&metrics{},
	)
	logrus.Info("Start server")
	err := s.Serve()
	if err != nil {
		logrus.WithError(err).Fatal("Serve error")
	}
}

type metrics struct {
}

func (m *metrics) APIDurationFunc(service string, method string, do func() (respType models.ResponseType)) {
	do()
}
