package context

import (
	"github.com/andrskom/gorpcgen/protocol/models"
	"github.com/sirupsen/logrus"
)

type Context struct {
	meta   *models.RequestMeta
	logger *logrus.Entry
}

func NewContext(meta *models.RequestMeta, logger *logrus.Entry) *Context {
	return &Context{
		meta:   meta,
		logger: logger,
	}
}

func (c *Context) GetMeta() *models.RequestMeta {
	return c.meta
}

func (c *Context) GetLogger() *logrus.Entry {
	return c.logger
}
