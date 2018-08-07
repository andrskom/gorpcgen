// +build teste2e

package testing

import (
	"testing"

	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/nats/gen/client"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

func Test_Service(t *testing.T) {
	a := assert.New(t)
	nc, err := nats.Connect("nats://nats:nats@localhost:4222")
	a.NoError(err)
	c := client.NewClient(nc)
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
