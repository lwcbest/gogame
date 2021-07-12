package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type EnvConfig struct {
	GameWebUrl    string `json:"gameWebUrl"`
	DataBaseUrl   string `json:"dataBaseUrl"`
	DataSecret    string `json:"dataSecret"`
	LotteryUrl    string `json:"lotteryUrl"`
	LotteryKey    string `json:"lotteryKey"`
	LotterySecret string `json:"lotterySecret"`
}

type RedisConfig struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

type MongoConfig struct {
	Url       string `json:"url"`
	Db        string `json:"db"`
	Authdb    string `json:"auth_db"`
	Authuser  string `json:"auth_user"`
	Authpass  string `json:"auth_pass"`
	Timeout   int64  `json:"timeout"`
	Poollimit int    `json:"pool_limit"`
	IsAuth    bool   `json:"is_auth"`
}

type MailConfig struct {
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	User    string   `json:"user"`
	Pass    string   `json:"pass"`
	ToMails []string `json:"toMails"`
}

var (
	EnvConf   *EnvConfig
	RedisConf *RedisConfig
	MongoConf *MongoConfig
)

func IsOnline() bool {
	env := os.Getenv("environmentName")
	return env == "online" || env == "pre"
}

func init() {
	env := os.Getenv("environmentName")
	if env == "" {
		env = "local"
	}
	EnvConf = &EnvConfig{}
	readConfig(env, "env.json", EnvConf)
	RedisConf = &RedisConfig{}
	readConfig(env, "redis.json", RedisConf)
	MongoConf = &MongoConfig{}
	readConfig(env, "mongo.json", MongoConf)
}

func readConfig(env string, filePath string, v interface{}) {
	fBytes, err := ReadFile(path.Join("conf", env, filePath))
	if err != nil {
		panic("read config error" + filePath)
	}
	if err = json.Unmarshal(fBytes, v); err != nil {
		fmt.Printf("%s %+v %+v\n", fBytes, err, v)
		panic("parse config error " + filePath)
	}
	return
}

func ReadFile(path string) ([]byte, error) {
	fBytes, err := ioutil.ReadFile(path)
	if err == nil {
		return fBytes, nil
	}
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		return nil, err
	}
	fBytes, err = ioutil.ReadFile(filepath.Join(dir, path))
	return fBytes, err
}

func ReDirPath(path string) string {
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	dir, err1 := filepath.Abs(filepath.Dir(os.Args[0]))
	if err1 != nil {
		return path
	}
	return filepath.Join(dir, path)
}
