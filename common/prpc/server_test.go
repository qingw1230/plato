package prpc

import (
	"context"
	"testing"

	"github.com/qingw1230/plato/common/config"

	"github.com/qingw1230/plato/common/prpc/example/helloservice"

	ptrace "github.com/qingw1230/plato/common/prpc/trace"
	"google.golang.org/grpc"
)

const (
	testIp   = "127.0.0.1"
	testPort = 8867
)

func TestNewPServer(t *testing.T) {
	config.Init("../../im.yaml")

	ptrace.StartAgent()
	defer ptrace.StopAgent()

	s := NewPServer(WithServiceName("im_server"), WithIP(testIp), WithPort(testPort), WithWeight(100))
	s.RegisterService(func(server *grpc.Server) {
		helloservice.RegisterGreeterServer(server, helloservice.HelloServer{})
	})
	s.Start(context.TODO())
}
