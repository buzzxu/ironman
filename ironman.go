package ironman

import (
	"context"
	"fmt"
	"github.com/buzzxu/ironman/conf"
	"github.com/buzzxu/ironman/echomw"
	"github.com/buzzxu/ironman/logger"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Server 启动服务
func Server(e *echo.Echo) {
	//初始化日志
	logger.InitLogger()
	log, err := logger.New("access", "info", true, true, false)
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
	//平滑关闭
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", conf.ServerConf.Port)); err != nil {
			logger.Error(err.Error())
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Fatale("Shutdown error", err)
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

func httpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "StatusInternalServerError"
		msg  interface{}
	)

	if e, ok := err.(*Error); ok {
		code = e.Code
		key = e.Key
		msg = e.Message
	} else if e, ok := err.(*echo.HTTPError); ok {
		code = e.Code
		key = http.StatusText(code)
		msg = e.Message
	} else if c.Echo().Debug {
		msg = err.Error()
	} else {
		key = http.StatusText(code)
		msg = err.Error()
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			err := c.JSON(code, NewError(code, key, msg))
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}
}

func (e *Error) Error() string {
	return e.Key + ": " + e.Message.(string)
}
