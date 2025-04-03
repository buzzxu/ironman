package jwtt

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"testing"
)

func TestJwtConfig(t *testing.T) {

}

func TestSetErrorConfig(t *testing.T) {
	// 默认情况下使用内置错误处理
	_ = JwtConfig(middleware.DefaultSkipper)
	//e.Use(echojwt.WithConfig(jwtConfig))

	// 创建自定义错误处理配置
	customErrorConfig := ErrorConfig{
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"status":  "error",
				"code":    "AUTH_FAILED",
				"message": "认证失败：" + err.Error(),
			})
		},
		ErrorStatusCodes: map[string]int{
			ErrJWTMissing.Error():      http.StatusBadRequest,
			ErrJWTInvalid.Error():      http.StatusUnauthorized,
			ErrJWTUnauthorized.Error(): http.StatusForbidden,
		},
		ErrorMessages: map[string]string{
			ErrJWTMissing.Error():      "请提供有效的认证令牌",
			ErrJWTInvalid.Error():      "认证令牌无效或已过期",
			ErrJWTUnauthorized.Error(): "您没有权限访问此资源",
		},
	}

	// 设置错误配置
	SetErrorConfig(customErrorConfig)

	// 仅修改特定错误的消息
	SetErrorMessage(ErrJWTInvalid, "您的登录已过期，请重新登录")
	SetErrorStatusCode(ErrJWTMissing, http.StatusUnauthorized)

	// 仅自定义错误处理函数
	SetErrorHandler(func(c echo.Context, err error) error {
		// 记录日志
		log.Printf("JWT Authentication Error: %v", err)

		// 根据不同的错误类型返回适当的响应
		switch err.Error() {
		case ErrJWTMissing.Error():
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "认证令牌缺失"})
		case ErrJWTInvalid.Error():
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "令牌已过期"})
		default:
			return c.JSON(http.StatusForbidden, map[string]string{"message": "访问被拒绝"})
		}
	})
}
