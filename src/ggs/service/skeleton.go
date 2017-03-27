package service

import (
	"ggs/chanrpc"
	"ggs/conf"
	"ggs/log"
	"ggs/timer"
	"time"
)

// 骨架
type Skeleton struct {
	chanRPCServer *chanrpc.Server   // PRC服务器引用
	dispatcher    *timer.Dispatcher // 定时器引用
}

// 创建骨架
func NewSkeleton() *Skeleton {
	skeleton := &Skeleton{
		chanRPCServer: chanrpc.NewServer(conf.Env.ChanRPCLen),           // 创建RPC服务器
		dispatcher:    timer.NewDispatcher(conf.Env.TimerDispatcherLen), // 创建定时器
	}
	return skeleton
}

// 获取创建的RPC服务器
func (s *Skeleton) ChanRPCServer() *chanrpc.Server {
	return s.chanRPCServer
}

// 实现Service接口的Run方法
func (s *Skeleton) Run(closeSig chan bool) {
	for { // 一直循环
		select {
		case <-closeSig: //读取到关闭信号
			s.chanRPCServer.Close() // 关闭RPC服务器
			return
		case ci := <-s.chanRPCServer.ChanCall: // 从RPC服务器读取调用信息
			err := s.chanRPCServer.Exec(ci) // 执行调用
			if err != nil {
				log.Error("%v", err)
			}
		case t := <-s.dispatcher.ChanTimer: // 从定时器中读取到定时信息
			t.Cb() // 执行定时器回调
		}
	}
}

// 向RPC服务器注册函数f
func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.chanRPCServer == nil { // 外部没有传入RPC服务器
		panic("invalid ChanRPCServer") // 抛错
	}

	s.chanRPCServer.Register(id, f) // 注册函数f
}

// 定时器
func (s *Skeleton) AfterFunc(d time.Duration, cb func()) *timer.Timer {
	return s.dispatcher.AfterFunc(d, cb)
}

// 支持cron表达式的定时器
func (s *Skeleton) CronFunc(cronExpr *timer.CronExpr, cb func()) *timer.Cron {
	return s.dispatcher.CronFunc(cronExpr, cb)
}
