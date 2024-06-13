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

	server.WatchService("/web/", func(key, value string) {}, func(key, value string) {})
	server.WatchService("/gRPC/", func(key, value string) {}, func(key, value string) {})
	for {
		select {
		case <-time.Tick(5 * time.Second):
		}
	}

}
