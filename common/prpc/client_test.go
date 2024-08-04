package prpc

import (
	"testing"

	"github.com/qingw1230/plato/common/config"

	ptrace "github.com/qingw1230/plato/common/prpc/trace"
	"github.com/stretchr/testify/assert"
)

func TestNewPClient(t *testing.T) {
	config.Init("../../im.yaml")
	ptrace.StartAgent()
	defer ptrace.StopAgent()

	_, err := NewPClient("im_server")
	assert.NoError(t, err)
}
