package server

import (
	"github.com/gin-gonic/gin"
	"github.com/lfxnxf/zdy_tools/logging"
	"github.com/lfxnxf/zdy_tools/zd_http"
	"path/filepath"
	"testing"
)

func TestHttpServer_Start(t *testing.T) {

	accessLog := logging.NewLogging(filepath.Join("./", "access.log"))
	logging.DefaultKit = logging.NewKit(accessLog, nil, nil, nil, nil, nil)

	s := NewHttpServer(HttpServerConfig{
		ServiceName: "test",
		Port:        8888,
		Mode:        DebugMode,
	})

	s.GET("/test", func(c *gin.Context) {
		resp := struct {
			Name string `json:"name"`
			Age  int64  `json:"age"`
		}{
			Name: "张三",
			Age:  100,
		}
		zd_http.WriteJson(c, resp, nil)
	})

	if err := s.Start(); err != nil {
		logging.Fatalf("http server start failed, err %v", err)
	}
}
