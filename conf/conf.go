package conf

import (
	"github.com/buzzxu/boys/common/files"
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
		Logger     *Logger                `yaml:"logger"`
		Props      map[string]interface{} `yaml:"props"`
		Jwt        *jwtConf               `yaml:"jwt"`
		DataSource *dataSource            `yaml:"dataSource"`
		Redis      *redisConf             `yaml:"redis"`
		MongoDb    *mongoDb               `yaml:"mongo"`
		WorkDir    string
	}
	Logger struct {
		Dir        string              `yaml:"dir"`
		MaxSize    int                 `yaml:"maxSize"`
		MaxBackups int                 `yaml:"maxBackups"`
		MaxAge     int                 `yaml:"maxAge"`
		Compress   bool                `yaml:"compress"`
		Line       bool                `yaml:"line"`
		Console    bool                `yaml:"console"`
		Json       bool                `yaml:"json"`
		Loggers    map[string]*LogConf `yaml:"loggers"`
	}

	LogConf struct {
		Level   string `yaml:"level"`
		File    string `yaml:"file"`
		Json    bool   `yaml:"json"`
		Console bool   `yaml:"console"`
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
	confFile := currentDir + "/app.yml"
	if files.Exists(confFile) {
		LoadConf(confFile)
		return
	}
	confFile = currentDir + "/app.yaml"
	if files.Exists(confFile) {
		LoadConf(confFile)
	} else {
		log.Fatalf("未找到 %s,请确认约定的文件名[app.yml,app.yaml]", confFile)
	}
}

func LoadConf(conf string) {
	if !files.Exists(conf) {
		log.Fatalf("未找到 %s", conf)
	}
	yamlFile, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatalf("配置文件读取失败 #%v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, ServerConf)
	if err != nil {
		log.Fatalf("配置文件解析失败: %v", err)
		os.Exit(1)
	}
	// 设置go processor数量
	if ServerConf.MaxProc == 0 {
		ServerConf.MaxProc = runtime.NumCPU()
	}
	if ServerConf.Port == "" {
		ServerConf.Port = "3000"
	}
	currentDir, _ := os.Getwd()
	ServerConf.WorkDir = currentDir
	//默认日志配置
	if ServerConf.Logger == nil {
		ServerConf.Logger = &Logger{
			MaxSize:    50,
			MaxAge:     30,
			MaxBackups: 30,
			Compress:   true,
			Json:       true,
			Console:    true,
			Line:       true,
			Dir:        currentDir,
		}
	} else {
		if ServerConf.Logger.MaxSize == 0 {
			ServerConf.Logger.MaxSize = 50
		}
		if ServerConf.Logger.MaxAge == 0 {
			ServerConf.Logger.MaxAge = 30
		}
		if ServerConf.Logger.MaxBackups == 0 {
			ServerConf.Logger.MaxBackups = 30
		}
		if !ServerConf.Logger.Line {
			ServerConf.Logger.Line = true
		}
		if !ServerConf.Logger.Compress {
			ServerConf.Logger.Compress = true
		}
		if ServerConf.Logger.Dir == "" {
			ServerConf.Logger.Dir = ServerConf.WorkDir
		}
	}
}
