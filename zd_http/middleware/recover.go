package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lfxnxf/zdy_tools/logging"
	"go.uber.org/zap"
	"runtime"
)

func recoverSysMW() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				err := fmt.Errorf("errgroup: panic recovered: %s\n%s", r, buf)
				logging.Errorw("mw_sys_recover_happen", zap.Error(err))
				logging.CrashLog(err)
			}
		}()
		c.Next()
	}
}
