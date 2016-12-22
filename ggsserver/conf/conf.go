package conf

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"ggs/conf"
	"ggs/log"
)

var Server = new(ServerConf)

type ServerConf struct {
	Names     map[string]string
	EnvPath   string
	RedisPort string
	ApiPort   string
	StartTime string
}

var Env struct {
	// ...
}

func init() {
	loadServer()
	loadEnv()
}

// 每个服务器自己的配置
func loadServer() {
	data, err := ioutil.ReadFile(path.Join(conf.EnvPath, "ybzt.json")) // eg. ~/dbsgz/servers/predev/10000/100000000/servers/ybzt.json
	if err != nil {
		log.Fatal("read ybzt.json error: %v, path: %v", err, conf.EnvPath)
	}

	err := json.Unmarshal(data, Server) // 解码数据, 存放到Server中
	if err != nil {
		log.Fatal("unmarshal ybzt.json error: %v", err)
	}

	if len(Server.Names) == 0 {
		log.Fatal("server names no set")
	}
}

// 通用的配置
func loadEnv() {
	data, err := ioutil.ReadFile(path.Join(Server.EnvPath, "conf.json"))
	if err != nil {
		log.Fatal("read conf file error: %v, path: %v", err, Server.EnvPath)
	}

	err := json.Unmarshal(data, &Env) // 解码数据，存放到Env中
	if err != nil {
		log.Fatal("unmarshal conf file error: %v", err)
	}
}
