package ironman

import (
	"flag"
	"github.com/buzzxu/ironman/conf"
	"github.com/labstack/echo/v4"
	"runtime"
	"testing"
)

func init() {
	conf.LoadDefaultConf()
}

func TestServer(t *testing.T) {
	runtime.GOMAXPROCS(conf.ServerConf.MaxProc)
	flag.Parse()
	// 关闭redis
	defer RedisClose()
	// 关闭数据库
	defer DataSourceClose()
	echo := echo.New()
	Server(echo)
}
