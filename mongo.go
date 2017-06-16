package ironman

import (
	"github.com/buzzxu/ironman/conf"
	"gopkg.in/mgo.v2"
)

// MongoDb 数据库
var MongoDb *mgo.Database

//MongoDbConnect 连接Mongodb并获取Database
func MongoDbConnect() {
	conf := conf.ServerConf.MongoDb
	session, err := mgo.Dial(conf.Url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	MongoDb = session.DB(conf.DB)
	if err := MongoDb.Login(conf.User, conf.Password); err != nil {
		panic(err)
	}
}
