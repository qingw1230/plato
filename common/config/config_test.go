package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	path := "../../im.yaml"
	Init(path)
	fmt.Println(IsDebug())
	fmt.Println(GetServicePathForIPConf())
}
