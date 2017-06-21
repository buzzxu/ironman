package ironman

import (
	"github.com/buzzxu/ironman/conf"
	"gopkg.in/mgo.v2"
	"fmt"
)






//MongoSession
var MongoSession *mgo.Session

//MongoDbConnect 连接Mongodb并获取Database
func MongoDbConnect() {
	MongoSession = createMongoDbSession()
}

func createMongoDbSession()*mgo.Session  {
	conf := conf.ServerConf.MongoDb
	session, err := mgo.Dial(conf.Url)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	fmt.Println(conf)
	credential := mgo.Credential{
		Username: conf.User,
		Password: conf.Password,
		Source:conf.Database,
	}
	if err := session.Login(&credential); err != nil {
		fmt.Print(err)
	}
	return session
}