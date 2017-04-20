# ggs
Go Game Server


# 执行流程

游戏服务器在启动时需要进行模块的注册，例如：
```
ggs.Run(login.Service, game.Service, gate.Service)
hl
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

由于 ggs 中每个模块都跑在一个单独的 goroutine 上，为了模块间方便的相互调用就有了基于 channel 的RPC机制。
一个 ChanRPC 需要在游戏服务器初始化的时候进行注册(注意：注册过程不是并发安全的)。
例如 ggsserver 中 login 模块注册了 NewAgent 和 CloseAgent 两个ChanRPC：
```
package login

import (
    "ggs/service"
    ...
)

var skeleton = service.NewSkeleton()

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {

}

func rpcCloseAgent(args []interface{}) {

}
```

使用 skeleton 来注册 ChanRPC。RegisterChanRPC 的第一个参数是 ChanRPC 的名字，第二个参数是 ChanRPC 的实现。
这里的 NewAgent 和 CloseAgent 会被 ggsserver 的 gate 模块在连接建立和连接中断时调用。

查看 service 包中的 RegisterChanRPC 函数：
```
func (s *Skeleton) RegisterChanRPC(id interface{}, f interface{}) {
	if s.chanRPCServer == nil {
		panic("invalid ChanRPCServer")
	}

	s.chanRPCServer.Register(id, f)
}
```

进一步查看 chanrpc 包的 Register 函数：
```
func (s *Server) Register(id interface{}, f interface{}) {
	switch f.(type) {
	case func([]interface{}):
	case func([]interface{}) interface{}:
	case func([]interface{}) []interface{}:
	default:
		panic(fmt.Sprintf("function id %v: definition of function is invalid", id))
	}

	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.functions[id] = f
}
```

可以看到它会存储 id => func 的映射，即存储了 ChanRPC 的名称到其实现的映射关系。

当有游戏客户端接入时，经由 gate 模块的 Run 方法时：
```
func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if conf.Env.WSAddr != "" {
		wsServer = new(network.WSServer)
	}

	if wsServer != nil {
		wsServer.Start(func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a)
			}
			return a
		})
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
}
```

执行到 `gate.AgentChanRPC.Go("NewAgent", a)` (注意这里只实现了 websocket 方式)，查看 serveces 包的 Go 方法：
```
func (s *Server) Go(id interface{}, args ...interface{}) {
	f := s.functions[id]
	if f == nil {
		return
	}

	defer func() {
		recover()
	}()

	s.ChanCall <- &CallInfo{
		f:    f,
		args: args,
	}
}
```

这里只是将包装的 CallInfo 对象(ChanRPC的实现和需要的参数) 通过 channel 管道传输。

由于之前在 ggserver.go 中：
```
gate.Service.AgentChanRPC = login.ChanRPCServer
```

指定了 gate 模块的代理 ChanRPC 为 login 模块的 ChanRPC。
因此在 login 模块的 Run 方法中：
```
func (s *Skeleton) Run(closeSig chan bool) {
	for {
		select {
		case <-closeSig:
			s.chanRPCServer.Close()
			return
		case ci := <-s.chanRPCServer.ChanCall:
			err := s.chanRPCServer.Exec(ci)
			if err != nil {
				log.Error("%v", err)
			}
		case t := <-s.dispatcher.ChanTimer:
			t.Cb()
		}
	}
}
```

会执行到 `s.chanRPCServer.Exec(ci)`，注意这里 gate 和 login 模块的 Run 方法是并发运行的。

查看 chanRPC 包的 Exec 方法：
```
func (s *Server) Exec(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if conf.Env.StackBufLen > 0 {
				buf := make([]byte, conf.Env.StackBufLen)
				l := runtime.Stack(buf, false)
				err = fmt.Errorf("%v: %s", r, buf[:l])
			} else {
				err = fmt.Errorf("%v", r)
			}

			s.ret(ci, &RetInfo{err: fmt.Errorf("%v", r)})
		}
	}()

	// execute
	switch ci.f.(type) {
	case func([]interface{}):
		ci.f.(func([]interface{}))(ci.args)
		return s.ret(ci, &RetInfo{})
	case func([]interface{}) interface{}:
		ret := ci.f.(func([]interface{}) interface{})(ci.args)
		return s.ret(ci, &RetInfo{ret: ret})
	case func([]interface{}) []interface{}:
		ret := ci.f.(func([]interface{}) []interface{})(ci.args)
		return s.ret(ci, &RetInfo{ret: ret})
	}

	panic("bug")
}
```

NewAgent 这个 ChanRPC 会 执行到 `return s.ret(ci, &RetInfo{})`，查看 ret 方法：
```
func (s *Server) ret(ci *CallInfo, ri *RetInfo) (err error) {
	if ci.chanRet == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}
```

此时 `ci.chanRet` 为 nil 直接返回。

同理在游戏服务器退出时会调用 gate 模块的 OnClose 方法，CloseAgent 也会被调用。
