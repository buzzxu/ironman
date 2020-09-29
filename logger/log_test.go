package logger

import (
	"errors"
	"github.com/buzzxu/ironman/conf"
	"testing"
)

func init() {
	conf.LoadDefaultConf()
	InitLogger()
}

func TestNewCompatibleLogger(t *testing.T) {

	Logger("order").Errorf("fdfdfdf %v", errors.New("测试异常"))
	Logger("order").Debug("我是debug")
	Logger("error").Info("323>>>")
}
