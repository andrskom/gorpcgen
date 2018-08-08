// +build teste2e

package testing

import (
	"testing"

	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/andrskom/gorpcgen/testservice/handlers"
	"github.com/andrskom/gorpcgen/testservice/http/gen/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
)

func Test_Service(t *testing.T) {
	a := assert.New(t)

	c := client.NewClient(
		logrus.WithField("client", "http"),
		http.DefaultClient,
		"http://localhost:4020/",
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
