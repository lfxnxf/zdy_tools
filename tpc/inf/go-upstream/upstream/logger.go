package upstream

import (
	log "github.com/lfxnxf/zdy_tools/logging"
)

var (
	logging *log.Logger
)

func init() {
	logging = log.New()
}

func SetLogger(l *log.Logger) {
	logging = l
}
