package ggs

import (
	"ggs/console"
	"ggs/log"
	"ggs/service"
	"os"
	"os/signal"
)

func Run(services ...service.Service) {
	log.Info("GGS is starting up...")

	// 注册所有模块, 保存到模块数组中
	for i := 0; i < len(services); i++ {
		service.Register(services[i])
	}

	// 遍历模块数组,调用各模块的OnInit方法,等到所有模块的OnInit方法执行完后则为每一个模块启动一个goroutine来执行模块的Run方法
	service.Init()

	console.Init()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c // 接收到退出标志(如ctrl+c关闭游戏服务器)

	log.Info("GGS is closing down... (signal: %v)", sig)
	console.Destroy()

	// 按与模块注册相反顺序在同一个goroutine中执行模块的OnDestroy方法
	service.Destroy()

	log.Info("BYE")
}
