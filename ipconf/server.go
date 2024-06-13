package ipconf

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/qingw1230/plato/common/config"
	"github.com/qingw1230/plato/ipconf/domain"
	"github.com/qingw1230/plato/ipconf/source"
)

func RunMain(path string) {
	config.Init(path)
	source.Init()
	domain.Init()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIPInfoList)
	s.Spin()
}
