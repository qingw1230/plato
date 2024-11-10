package config

import (
	"time"

	"github.com/spf13/viper"
)

func GetEndpointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

func GetTimeouForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

func GetServicePathForIPConf() string {
	return viper.GetString("ip_conf.service_path")
}
