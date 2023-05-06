package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetOpts() []gin.HandlerFunc {
	// todo max_connects
	// todo prometheus
	return []gin.HandlerFunc{
		loggingAccess(), // 生成access_log
		setTrace(),      // 设置trace
		recoverSysMW(),  // recover
		crossDomain(),   // 跨域设置
	}
}
