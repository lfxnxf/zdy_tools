package rpc_client

import (
	"google.golang.org/grpc"

	"github.com/lfxnxf/zdy_tools/logging"
	"github.com/lfxnxf/zdy_tools/tools/syncx"
	"github.com/lfxnxf/zdy_tools/zd_rpc/middleware"
)

type RpcClient struct {
	conf         RpcClientConf
	conn         *grpc.ClientConn
	singleFlight syncx.SingleFlight
}

func NewRpcClient(c RpcClientConf) *RpcClient {
	return &RpcClient{
		conf:         c,
		singleFlight: syncx.NewSingleFlight(),
	}
}

type RpcClientConf struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

func (c *RpcClient) GetRpcConn(options ...grpc.DialOption) *grpc.ClientConn {
	if c.conn != nil {
		return c.conn
	}
	_, _ = c.singleFlight.Do(c.conf.Name, func() (interface{}, error) {
		opt := middleware.GetClientOpts()
		opt = append(opt, options...)
		conn, err := grpc.Dial(c.conf.Address, opt...)
		if err != nil {
			logging.Fatalf("did not connect: %v", err)
			return nil, err
		}
		c.conn = conn
		return conn, nil
	})
	return c.conn
}

func (c *RpcClient) Close() {
	_ = c.conn.Close()
}
