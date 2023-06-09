package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/lfxnxf/zdy_tools/logging"
	"github.com/lfxnxf/zdy_tools/zd_http/http_ctx"
	"github.com/lfxnxf/zdy_tools/zd_http/middleware"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

type HttpServerConfig struct {
	ServiceName string `yaml:"service_name"`
	Port        int64  `yaml:"port"`
	Mode        string `yaml:"mode"`
	HttpsPort   int64  `yaml:"https_port"`
	Crt         string `yaml:"crt"`
	Key         string `yaml:"key"`
}

type HttpServer struct {
	*gin.Engine
	cfg         HttpServerConfig
	server      *http.Server
	httpsServer *http.Server
}

type HttpRoute struct {
	route   string
	handler http_ctx.HttpHandler
}

func NewHttpServer(cfg HttpServerConfig) *HttpServer {
	// 设置运行模式
	gin.SetMode(cfg.Mode)

	engine := gin.New()
	s := &HttpServer{
		Engine: engine,
		cfg:    cfg,
	}

	pprof.Register(engine) // 性能

	// 初始化中间件
	s.initPublicMiddleware()
	return s
}

func (s *HttpServer) StartHttp() error {
	s.server = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.cfg.Port),
		Handler:        s.Engine,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.server.ListenAndServe()
	if err != nil {
		logging.Errorw("start http server failed %v", zap.Error(err))
	}
	return err
}

func (s *HttpServer) StartHttps() error {
	s.httpsServer = &http.Server{
		Addr:           fmt.Sprintf(":%d", s.cfg.HttpsPort),
		Handler:        s.Engine,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.httpsServer.ListenAndServeTLS(s.cfg.Crt, s.cfg.Key)
	if err != nil {
		logging.Errorw("start http server failed %v", zap.Error(err))
	}
	return err
}

func (s *HttpServer) Shutdown(ctx context.Context) {
	err := s.server.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
	}

	if s.httpsServer != nil {
		err = s.httpsServer.Shutdown(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *HttpServer) initPublicMiddleware() {
	// 设置中间件
	s.Use(middleware.GetOpts()...)
}
