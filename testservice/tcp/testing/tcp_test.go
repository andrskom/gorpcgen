// +build teste2e

package testing

import (
	"testing"

	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/tcp/gen/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	a := assert.New(t)

	c := client.New(
		logrus.WithField("client", "tcp"),
		"localhost",
		"8081",
	)
	meta := models.RequestMeta{RequestID: "reqId", Origin: "client"}
	{
		// ok
		resp, errModel := c.Get(meta, &handlers.GetRequest{Type: handlers.TypeOK})
		a.Nil(errModel)
		a.NotNil(resp)
		a.Equal(handlers.HandlerOK, resp)
	}
	{
		// err
		resp, errModel := c.Get(meta, &handlers.GetRequest{Type: handlers.TypeErr})
		a.Nil(resp)
		a.NotNil(errModel)
		a.Equal(handlers.HandlerErr, errModel)
	}
}
