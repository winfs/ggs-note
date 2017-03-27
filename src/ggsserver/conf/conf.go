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
	Names     map[string]string // 服务器id => 名称
	EnvPath   string            // 配置文件路径
	RedisPort string            // Redis端口
	ApiPort   string            //Api端口
	StartTime string
}

var Env struct {
	// ...
}

func init() {
	loadServer()
	loadEnv()
}

// 加载服务器配置
func loadServer() {
	data, err := ioutil.ReadFile(path.Join(conf.EnvPath, "ybzt.json")) // eg. ~/dbsgz/servers/dev-winfs/10000/100000000/ybzt.json
	if err != nil {
		log.Fatal("read ybzt.json error: %v, path: %v", err, conf.EnvPath)
	}

	err = json.Unmarshal(data, Server) // 解码数据, 存放到Server中
	if err != nil {
		log.Fatal("unmarshal ybzt.json error: %v", err)
	}

	if len(Server.Names) == 0 {
		log.Fatal("server names no set")
	}
}

// 通用配置
func loadEnv() {
	data, err := ioutil.ReadFile(path.Join(Server.EnvPath, "conf.json")) // eg. ~/dbsgz/configs/configs/dev-winfs/conf.json
	if err != nil {
		log.Fatal("read conf file error: %v, path: %v", err, Server.EnvPath)
	}

	err = json.Unmarshal(data, &Env) // 解码数据，存放到Env中
	if err != nil {
		log.Fatal("unmarshal conf file error: %v", err)
	}
}
