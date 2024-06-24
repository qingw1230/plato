package prome

import (
	"github.com/prometheus/client_golang/prometheus"
)

func NewCountVec(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	counterVer := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(counterVer)
	return counterVer
}

func NewHistogramVec(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	histogramVec := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(histogramVec)
	return histogramVec
}
