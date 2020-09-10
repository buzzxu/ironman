package conf

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type (
	serverConf struct {
		MaxProc    int                    `yaml:"maxProc"`
		Port       string                 `yaml:"port"`
		Log        string                 `yaml:"log"`
		Props      map[string]interface{} `yaml:"props"`
		Jwt        *jwtConf               `yaml:"jwt"`
		DataSource *dataSource            `yaml:"dataSource"`
		Redis      *redisConf             `yaml:"redis"`
		MongoDb    *mongoDb               `yaml:"mongo"`
	}

	jwtConf struct {
		ContextKey    string        `yaml:"contextKey"`
		SigningKey    string        `yaml:"signingKey"`
		AuthScheme    string        `yaml:"authScheme"`
		SigningMethod string        `yaml:"signingMethod"`
		Expires       time.Duration `yaml:"expires"`
	}
	redisConf struct {
		Addr         string `yaml:"addr"`
		Password     string `yaml:"password"`
		DB           int    `yaml:"db"`
		PoolSize     int    `yaml:"poolSize"`
		MaxRetries   int    `yaml:"maxRetries"`
		MinIdleConns int    `yaml:"minIdle"`
		Stats        bool   `yaml:"stats"`
	}
	dataSource struct {
		Host            string `yaml:"host"`
		Port            int16  `yaml:"port"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		DB              string `yaml:"db"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifetime int    `yaml:"connMaxLifetime"`
		Log             bool   `yaml:"log"`
	}
	mongoDb struct {
		Url       string `yaml:"url"`
		Database  string `yaml:"db"`
		User      string `yaml:"user"`
		Password  string `yaml:"password"`
		Timeout   int    `yaml:"timeout"`
		PoolLimit int    `yaml:"poolLimit"`
	}
)

var ServerConf = &serverConf{}

func init() {

}
func LoadDefaultConf() {
	currentDir, _ := os.Getwd()
	LoadConf(currentDir + "/conf.yaml")
}

func LoadConf(conf string) {
	yamlFile, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, ServerConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(-1)
	}
	// 设置go processor数量
	if ServerConf.MaxProc == 0 {
		ServerConf.MaxProc = runtime.NumCPU()
	}
	if len(ServerConf.Log) == 0 {
		ServerConf.Log = "server.log"
	}
}
