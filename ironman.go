package ironman

import (
	"context"
	"errors"
	"fmt"
	"github.com/buzzxu/ironman/conf"
	"github.com/buzzxu/ironman/echomw"
	"github.com/buzzxu/ironman/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Server 启动服务
func Server(e *echo.Echo) {
	log, err := logger.New("access", "info", true, true, true)
	if err != nil {
		logger.Fatalf("logger[%s] 创建失败: %v", "access", err)
	}

	e.Use(echomw.ZapLogger(log.Unwrap()))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Secure())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		DisablePrintStack: true,
		DisableStackAll:   true,
		StackSize:         4 << 10,
	}))
	e.HTTPErrorHandler = httpErrorHandler
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server := echo.StartConfig{
		Address:         fmt.Sprintf(":%s", conf.ServerConf.Port),
		GracefulTimeout: 10 * time.Second,
	}
	if err := server.Start(ctx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}

// NewError 异常
func NewError(code int, key string, msg interface{}) *Error {
	return &Error{
		Code:    code,
		Key:     key,
		Message: msg,
	}
}

func httpErrorHandler(c *echo.Context, err error) {
	if response, _ := echo.UnwrapResponse(c.Response()); response != nil && response.Committed {
		return
	}
	code := http.StatusInternalServerError
	key := http.StatusText(code)
	var msg interface{} = err.Error()
	var appError *Error
	var httpError *echo.HTTPError
	switch {
	case errors.As(err, &appError):
		code, key, msg = appError.Code, appError.Key, appError.Message
	case errors.As(err, &httpError):
		code = httpError.Code
		key = http.StatusText(code)
		msg = httpError.Message
	default:
		if status := echo.StatusCode(err); status != 0 {
			code = status
			key = http.StatusText(code)
			msg = key
		}
	}
	var writeErr error
	if c.Request().Method == http.MethodHead {
		writeErr = c.NoContent(code)
	} else {
		writeErr = c.JSON(code, NewError(code, key, msg))
	}
	if writeErr != nil {
		c.Logger().Error("Failed to write error response", "error", writeErr)
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Key, e.Message)
}

// StatusCode lets Echo middleware resolve application error status codes.
func (e *Error) StatusCode() int { return e.Code }
