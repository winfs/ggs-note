package main

import (
	"ggs"
	"ggs/gate"
	"ggsserver/game"
	"ggsserver/login"
	"ggsserver/msg"
	"ggsserver/router"
)

func main() {
	gate.Service.Processor = msg.Processor          // 设置网关的消息处理器
	gate.Service.AgentChanRPC = login.ChanRPCServer // 设置网关的RPC服务器

	router.Init()
	ggs.Run(login.Service, game.Service, gate.Service)
}
