package domain

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
)

type IPConfConext struct {
	Ctx       *context.Context
	AppCtx    *app.RequestContext
	ClientCtx *ClientContext
}

type ClientContext struct {
	IP string `json:"ip"`
}

func BuildIPConfContext(c *context.Context, ctx *app.RequestContext) *IPConfConext {
	ipConfContext := &IPConfConext{
		Ctx:       c,
		AppCtx:    ctx,
		ClientCtx: &ClientContext{},
	}
	return ipConfContext
}
