package context

import (
	"github.com/andrskom/gorpcgen/protocol/logger"
	"github.com/andrskom/gorpcgen/protocol/models"
)

type Context struct {
	meta   *models.RequestMeta
	logger logger.Interface
}

func NewContext(meta *models.RequestMeta, logger logger.Interface) *Context {
	return &Context{
		meta:   meta,
		logger: logger,
	}
}

func (c *Context) GetMeta() *models.RequestMeta {
	return c.meta
}

func (c *Context) GetLogger() logger.Interface {
	return c.logger
}
