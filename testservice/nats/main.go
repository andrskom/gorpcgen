package main

import (
	"context"

	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/nats/gen/server"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
)

func main() {
	nc, err := nats.Connect("nats://nats:nats@localhost:4222")
	if err != nil {
		logrus.WithError(err).Fatal("nats connection error")
	}
	logrus.SetLevel(logrus.DebugLevel)

	s := server.NewServer(context.Background(), logrus.WithField("server", "nats"), nc, &handlers.Test{}, &metrics{})
	s.Serve()
}

type metrics struct {
}

func (m *metrics) APIDurationFunc(service string, method string, do func() (respType models.ResponseType)) {
	do()
}
