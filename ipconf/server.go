package ipconf

import (
	"github.com/gin-gonic/gin"
	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/ipconf/domain"
	"github.com/qingw1230/plato/ipconf/source"
)

func RunMain(path string) {
	config.Init(path)
	source.Init()
	domain.Init()
	r := gin.Default()
	r.GET("/ip/list", GetIPInfoList)
	r.Run(":6798")
}
