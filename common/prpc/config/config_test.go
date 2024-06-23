package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/qingw1230/plato/common/config"
)

func TestMain(m *testing.M) {
	config.Init("../../../im.yaml")
	os.Exit(m.Run())
}

func TestGetDiscovName(t *testing.T) {
	fmt.Println(GetDiscovName())
}

func TestGetDiscovEndpints(t *testing.T) {
	fmt.Println(GetDiscovEndpoints())
}

func TestGetTraceEnable(t *testing.T) {
	fmt.Println(GetTraceEnable())
}

func TestGetTraceCollectionURL(t *testing.T) {
	fmt.Println(GetTraceCollectionURL())
}

func TestGetTraceServiceName(t *testing.T) {
	fmt.Println(GetTraceServiceName())
}

func TestGetTraceSampler(t *testing.T) {
	fmt.Println(GetTraceSampler())
}
