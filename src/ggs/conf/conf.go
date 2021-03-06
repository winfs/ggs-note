// 配置相关
package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var Env struct {
	StackBufLen int
	LogLevel    string // 日志的错误级别(Debug、Info、Warn、Error、Fatal),不区分大小写
	LogPath     string // 日志文件路径, 不设置则打印到标准输出

	MaxConnNum      int    // 最大连接数
	PendingWriteNum int    // network中用于传递消息的管道长度
	MaxMsgLen       uint32 // 最大消息长度

	WSAddr      string        // WebSockeyt监听地址
	HTTPTimeout time.Duration // 请求超时时长

	ChanRPCLen         int // RPC服务器中用于传递调用信息的管道长度
	TimerDispatcherLen int // 定时器中用于传递Timer信息的管道长度

	// console
	ConsolePort   int
	ConsolePrompt string
	ProfilePath   string
}

var EnvPath string // 游戏服务器的配置文件路径
//var CrossEnvPath string // 跨服服务器的配置文件路径

func init() {
	initFlags()

	data, err := ioutil.ReadFile(path.Join(EnvPath, "ggs.env")) // 读取配置文件
	if err != nil {
		fmt.Println("env file not found, path: " + EnvPath)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &Env) // 解码数据, 存储到Env结构中
	if err != nil {
		fmt.Printf("env file format error: %v\n", err)
		os.Exit(1)
	}

	if Env.WSAddr != "" {
		if Env.MaxMsgLen <= 0 {
			Env.MaxMsgLen = 4096
			fmt.Println("invalid MaxMsgLen, reset to %v", Env.MaxMsgLen)
		}
		if Env.HTTPTimeout <= 0 {
			Env.HTTPTimeout = 10
			fmt.Println("invalid HTTPTimeout, reset to %v", Env.HTTPTimeout)
		}
		Env.HTTPTimeout *= time.Second
	}
}

// 注册命令行参数
func initFlags() {
	flag.StringVar(&EnvPath, "env", "", "path of env file")
	//flag.StringVar(&CrossEnvPath, "cross-env", "", "path of cross env file")

	flag.Parse()

	if EnvPath == "" {
		log.Fatal("flag env not set")
	}
	/*
		if CrossEnvPath == "" {
			log.Fatal("falg cross-env not set")
		}
	*/
}
