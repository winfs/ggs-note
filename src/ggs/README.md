# ggs
Go Game Server


# 执行流程

游戏服务器在启动时需要进行模块的注册，例如：
```
ggs.Run(login.Service, game.Service, gate.Service)
```

查看 ggs 包的 Run 函数，首先这里按顺序注册了 login、game、gate 三个模块：
```
for i := 0; i < len(services); i++ {
    service.Register(services[i])
}
```

其中每个模块都要实现 Service 接口：
```
type Service interface {
	OnInit()
	Run(closeSig chan bool)
	OnDestroy()
}
```

查看 service 包的 Register 函数可以发现：
```
var services []*service

func Register(si Service) {
    ...
    ...

	services = append(services, s)
}
```

注册的模块会按顺序追加到模块数组 services 中。

继续查看 ggs 包的 Run 函数：
```
service.Init()
```

查看 service 包的 Init 函数可以发现：
```
func Init() {
	for i := 0; i < len(services); i++ {
		services[i].si.OnInit()
	}

	for i := 0; i < len(services); i++ {
		go run(services[i])
	}
}
```

它会遍历我们的模块数组，然后按模块的注册顺序执行各模块的 OnInit 方法。
等到所有模块的 OnInit 方法执行完后，再次遍历我们的模块数组，为每一个模块分别启动一个goroutine来执行模块的 Run 方法。

其中 run 函数:
```
func run(s *service) {
	s.wg.Add(1)
	s.si.Run(s.closeSig)
	s.wg.Done()
}
```

继续查看 ggs 包的 Run 函数：
```
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, os.Kill)
sig := <-c
```

可以看到此时主goroutine会阻塞，直接接收到一个退出标志(如 Ctrl + C 关闭游戏服务器)才继续往下执行。

继续查看 ggs 包的 Run 函数：
```
service.Destroy()
```

查看 service 包的 Destroy 函数可以发现：
```
func Destroy() {
	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]
		s.closeSig <- true
		s.wg.Wait()
		destroy(s)
	}
}
```

当关闭游戏服务器时，会按与模块注册相反的顺序执行各模块的 OnDestroy 方法。

其中 destroy 函数：
```
func destroy(s *service) {
    ...
    ...

	s.si.OnDestroy()
}
```


# ChanRPC

由于ggs中每个模块都跑在一个独立的goroutine上，为了模块间方便的相互调用就有了基于channel的RPC机制。
一个ChanRPC需要在游戏服务器初始化的时候进行注册(注意：注册过程不是并发安全的)


# 消息发送流程 
客户端发送到游戏服务器的消息需要通过gate模块路由，简而言之，gate模块决定了某个消息具体交给内部的哪个模块来处理。
