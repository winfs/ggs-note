//配置相关
package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var Env struct {
	StackBufLen int
	LogLevel    string //日志的错误级别(Debug、Info、Warn、Error、Fatal),不区分大小写
	LogPath     string //日志文件的目录

	MaxConnNum      int    //最大连接数
	PendingWriteNum int    //network中用于传递消息的管道长度
	MaxMsgLen       uint32 //大小消息长度

	WSAddr      string        //WebSockeyt监听地址
	HTTPTimeout time.Duration //请求超时时长

	ChanRPCLen         int //RPC服务器中用于传递调用信息的管道长度
	TimerDispatcherLen int //定时分发器中用于传递Timer信息的管道长度

	ConsolePort   int
	ConsolePrompt string
	ProfilePath   string
}

var EnvPath string

func init() {
	flag.StringVar(&EnvPath, "env", "", "path of env file")
	flag.Parse()

	if EnvPath == "" {
		fmt.Println("param -env no set")
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(path.Join(EnvPath, "ggs.env"))
	if err != nil {
		fmt.Println("env file not found, path: " + EnvPath)
		os.Exit(1)
	}

	err = json.Unmarshal(data, &Env)
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
