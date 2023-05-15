package zdy_tools

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/lfxnxf/zdy_tools/inits"
	"github.com/lfxnxf/zdy_tools/zd_error"
	"github.com/lfxnxf/zdy_tools/zd_http"
)

func Test_Main(t *testing.T) {
	var cfg Config
	inits.Init(
		inits.ConfigPath("./config.yaml"),
		inits.Once(),
		inits.LoadLocalConfig(&cfg),
	)

	s := inits.NewHttpServer(cfg.Server)

	s.GET("/test", func(c *gin.Context) {
		zd_http.WriteJson(c, map[int]string{1: "word"}, zd_error.AddError("test error", "error"))
	})

	if err := s.StartHttp(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
