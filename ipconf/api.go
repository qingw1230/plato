package ipconf

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/qingw1230/plato/ipconf/domain"
)

// Response 获取机器列表时返回的响应体
type Response struct {
	// Message 响应的简短描述
	Message string `json:"message"`
	// Code 响应码
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func GetIPInfoList(c context.Context, ctx *app.RequestContext) {
	defer func() {
		// 不能因为一次请求失败服务就挂掉了，因此需要 recover
		if err := recover(); err != nil {
			ctx.JSON(consts.StatusBadRequest, utils.H{"err": err})
		}
	}()

	ipConfCtx := domain.BuildIPConfContext(&c, ctx)
	eds := domain.Dispatch(ipConfCtx)
	ipConfCtx.AppCtx.JSON(consts.StatusOK, packRes(top5Endports(eds)))
}
