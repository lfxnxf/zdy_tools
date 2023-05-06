package zdy_tools

import (
	"github.com/gin-gonic/gin"
	"github.com/lfxnxf/zdy_tools/inits"
	"github.com/lfxnxf/zdy_tools/zd_http/server"
	"net/http"
	"testing"
)

func Test_Main(t *testing.T) {
	inits.Init(
		inits.ConfigPath("./config.toml"),
		inits.Once(),
		inits.ConfigInstance(new(TomlConfig)),
	)
	s := inits.NewHttpServer(server.HttpServerConfig{
		ServiceName: "resource-center",
		Port:        8080,
		Mode:        gin.DebugMode,
		HttpsPort:   443,
		Crt:         "./ssl/tls.crt",
		Key:         "./ssl/tls.key",
	})

	s.GET("/test", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"hello": "word",
		})
	})

	go func() {
		if err := s.StartHttp(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}
