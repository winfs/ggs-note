package router

import (
	"ggsserver/game"
	"ggsserver/login"
	"ggsserver/msg"
)

// 设置某个消息具体交给哪个模块来处理
func Init() {
	msg.Processor.SetRouter(&msg.Hello{}, game.ChanRPCServer)
	msg.Processor.SetRouter(&msg.LoginRequest{}, login.ChanRPCServer)
}
