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

	MongoSession.DB("xca").C("base.xx").Insert(Map{"n": 1})
	fmt.Println()

}
