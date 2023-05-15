package rpc_server

import (
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"

	"github.com/lfxnxf/zdy_tools/zd_rpc/middleware"
)

type register func(server *grpc.Server)

type RpcServer struct {
	conf     RpcServerConfig
	options  []grpc.UnaryServerInterceptor
	register register
}

type RpcServerConfig struct {
	ServiceName string `yaml:"service_name"`
	Port        int64  `yaml:"port"`
}

func NewRpcServer(conf RpcServerConfig, register register) *RpcServer {
	return &RpcServer{
		conf:     conf,
		register: register,
	}
}

func (r *RpcServer) Use(options ...grpc.UnaryServerInterceptor) {
	r.options = append(r.options, options...)
}

func (r *RpcServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", r.conf.Port))
	if err != nil {
		return err
	}

	opts := middleware.GetServerOpts()

	opts = append(opts, r.options...)

	// 拦截器
	opt := grpc_middleware.WithUnaryServerChain(opts...)

	rpcServe := grpc.NewServer(opt)

	// 启动服务
	r.register(rpcServe)

	err = rpcServe.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}
