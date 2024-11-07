package discovery

import (
	"context"
	"testing"
	"time"
)

func TestServiceDiscovery(t *testing.T) {
	ctx := context.Background()
	server := NewServiceDiscovery(&ctx)
	defer server.Close()

	empty := func(key, value string) {}
	server.WatchService("/web", empty, empty)
	server.WatchService("/gRPC", empty, empty)
	for {
		select {
		case <-time.Tick(5 * time.Second):
		}
	}
}
