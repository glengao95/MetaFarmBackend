package middleware

import (
	"MetaFarmBackend/component/errors"
	"MetaFarmBackend/component/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 创建全局错误恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 转换为自定义错误类型
				var apiError *errors.Error
				switch e := err.(type) {
				case *errors.Error:
					apiError = e
				case error:
					apiError = errors.New(http.StatusInternalServerError, GetMsg(c, http.StatusInternalServerError)).WithError(e)
				default:
					apiError = errors.New(http.StatusInternalServerError, GetMsg(c, http.StatusInternalServerError)).WithMessage("unknown error")
				}

				// 记录错误日志
				logger.Errorf("Panic recovered: %+v\nStack: %s", err, apiError.WithStack())

				// 统一错误响应格式
				Fail(c, apiError.Code, nil)
			}
		}()

		c.Next()
	}
}
