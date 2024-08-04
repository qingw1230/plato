package config

import (
	"fmt"
	"testing"

	"github.com/qingw1230/plato/common/config"
)

func TestMain(m *testing.M) {
	config.Init("../../../im.yaml")
	m.Run()
}

func TestGetDiscovName(t *testing.T) {
	fmt.Println(GetDiscovName())
}

func TestGetDiscovEndpoints(t *testing.T) {
	fmt.Println(GetDiscovEndpoints())
}

func TestGetTraceEnable(t *testing.T) {
	fmt.Println(GetTraceEnable())
}

func TestGetTraceCollectionUrl(t *testing.T) {
	fmt.Println(GetTraceEnable())
}

func TestGetTraceServiceName(t *testing.T) {
	fmt.Println(GetTraceServiceName())
}

func TestGetTraceSampler(t *testing.T) {
	fmt.Println(GetTraceSampler())
}
