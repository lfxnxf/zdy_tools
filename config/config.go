package config

import (
	"github.com/lfxnxf/zdy_tools/resource/kafka"
	"github.com/lfxnxf/zdy_tools/resource/redis"
	"github.com/lfxnxf/zdy_tools/resource/sql"
	"github.com/lfxnxf/zdy_tools/trace"
	"github.com/lfxnxf/zdy_tools/zd_http/server"
	rpc_client "github.com/lfxnxf/zdy_tools/zd_rpc/client"
	rpc_server "github.com/lfxnxf/zdy_tools/zd_rpc/server"
)

type Instance interface {
	SetBase(Config)
}

type Config struct {
	Log           Log                        `yaml:"log"`
	Server        server.HttpServerConfig    `yaml:"server"`
	RpcServer     rpc_server.RpcServerConfig `yaml:"rpc_server"`
	RpcClient     []rpc_client.RpcClientConf `yaml:"rpc_client"`
	Telemetry     trace.Config               `yaml:"telemetry"`
	Database      []sql.GroupConfig          `yaml:"mysql"`
	Redis         []redis.Conf               `yaml:"redis"`
	KafkaProducer []kafka.ProducerConfig     `yaml:"kafka_producer_client"`
	KafkaConsumer []kafka.ConsumeConfig      `yaml:"kafka_consume"`
}

type Log struct {
	Level              string `yaml:"level"`
	Rotate             string `yaml:"rotate"`
	AccessRotate       string `yaml:"access_rotate"`
	AccessLog          string `yaml:"access_log"`
	BusinessLog        string `yaml:"business_log"`
	ServerLog          string `yaml:"server_log"`
	StatLog            string `yaml:"stat_log"`
	ErrorLog           string `yaml:"err_log"`
	LogPath            string `yaml:"log_path"`
	BalanceLogLevel    string `yaml:"balance_log_level"`
	GenLogLevel        string `yaml:"gen_log_level"`
	AccessLogOff       bool   `yaml:"access_log_off"`
	BusinessLogOff     bool   `yaml:"business_log_off"`
	RequestBodyLogOff  bool   `yaml:"request_log_off"`
	RespBodyLogMaxSize int    `yaml:"response_log_max_size"` // -1:不限制;默认1024字节;
	SuccessStatCode    []int  `yaml:"success_stat_code"`
	StorageDay         int64  `yaml:"storage_day"` // 日志保留时间
}
