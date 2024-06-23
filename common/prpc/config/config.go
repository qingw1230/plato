package config

import "github.com/spf13/viper"

func GetDiscovName() string {
	return viper.GetString("prpc.discov.name")
}

func GetDiscovEndpoints() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

func GetTraceEnable() bool {
	return viper.GetBool("prpc.trace.enable")
}

func GetTraceCollectionURL() string {
	return viper.GetString("prpc.trace.url")
}

func GetTraceServiceName() string {
	return viper.GetString("prpc.trace.service_name")
}

func GetTraceSampler() float64 {
	return viper.GetFloat64("prpc.trace.sampler")
}
