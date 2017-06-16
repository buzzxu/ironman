package ironman

import (
	"fmt"
	"testing"

	"github.com/buzzxu/ironman/conf"
)

func TestMongoDbConnectTest(t *testing.T) {
	conf.LoadConf()
	conf := conf.ServerConf.MongoDb
	fmt.Print("%v", &conf)

	MongoDbConnect()
	fmt.Println(&MongoDb)
}
