package zdy_tools

import (
	config2 "github.com/lfxnxf/zdy_tools/config"
)

type TomlConfig struct {
	config2.Config
}

func (c *TomlConfig) GetBase() config2.Config {
	return c.Config
}
