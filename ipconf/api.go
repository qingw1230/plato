package ipconf

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qingw1230/plato/ipconf/domain"
)

// Response 获取机器列表返回的响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Time    time.Time   `json:"time"`
	Data    interface{} `json:"data"`
}

func GetIPInfoList(ctx *gin.Context) {
	defer func() {
		// 不能因为一次请求失败服务就挂掉了
		if err := recover(); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"err": err})
		}
	}()

	eds := domain.Dispactch(ctx)
	ctx.JSON(http.StatusOK, packRes(top5Endpoints(eds)))
}
