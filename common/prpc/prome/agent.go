package prome

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var once sync.Once

// StartAgent 开启 prometheus
func StartAgent(host string, port int) {
	go func() {
		once.Do(func() {
			// 请求根目录时，调用 http.Handler 接口的 ServeHTTP(ResponseWriter, *Request) 方法
			http.Handle("/", promhttp.Handler())
			addr := fmt.Sprintf("%s:%d", host, port)
			logger.Infof("Starting prometheus agent at %s", addr)
			if err := http.ListenAndServe(addr, nil); err != nil {
				logger.Error(err)
			}
		})
	}()
}
