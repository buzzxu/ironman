package conf

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type (
	serverConf struct {
		Port       string      `yaml:"port"`
		Jwt        *jwtConf    `yaml:"jwt"`
		DataSource *dataSource `yaml:"dataSource"`
		Redis      *redisConf  `yaml:"redis"`
	}

	jwtConf struct {
		ContextKey    string        `yaml:"contextKey"`
		SigningKey    string        `yaml:"signingKey"`
		AuthScheme    string        `yaml:"authScheme"`
		SigningMethod string        `yaml:"signingMethod"`
		Expires       time.Duration `yaml:"expires"`
	}
	redisConf struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
		PoolSize int    `yaml:"poolSize"`
	}
	dataSource struct {
		Host            string `yaml:"host"`
		Port            int16  `yaml:"port"`
		User            string `yaml:"user"`
		Password        string `yaml:"password"`
		DbName          string `yaml:"db"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		ConnMaxLifetime int    `yaml:"connMaxLifetime"`
	}
)

var ServerConf = serverConf{}

func init() {

}
func LoadConf() {
	currentDir, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(currentDir + "/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &ServerConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(-1)
	}
}
