package zdy_tools

import (
	config2 "github.com/lfxnxf/zdy_tools/config"
)

type Config struct {
	Test string `yaml:"test"`
	config2.Config
}

func (c *Config) SetBase(cfg config2.Config) {
	c.Config = cfg
}
