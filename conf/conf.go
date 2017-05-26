package conf

import (
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type serverConf struct {
	Port string `yaml:"port"`
	Jwt  jwtConf
}

type jwtConf struct {
	ContextKey    string `yaml:"contextKey"`
	SigningKey    string `yaml:"signingKey"`
	AuthScheme    string `yaml:"authScheme"`
	SigningMethod string `yaml:"signingMethod"`
}

var ServerConf = serverConf{}

func init() {
	loadConf()
}
func loadConf() {
	currentDir, _ := os.Getwd()
	yamlFile, err := ioutil.ReadFile(currentDir + "/conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		os.Exit(-1)
	}
	err = yaml.Unmarshal(yamlFile, &ServerConf)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		os.Exit(-1)
	}
}
