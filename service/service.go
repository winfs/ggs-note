package service

import (
	"ggs/conf"
	"ggs/log"
	"runtime"
	"sync"
)

//模块接口(每个模块都需要实现这个接口)
type Service interface {
	OnInit()                //初始化函数
	Run(closeSig chan bool) //运行函数
	OnDestroy()             //销毁函数
}

//模块
type service struct {
	si       Service        //实现了模块接口的某个对象
	closeSig chan bool      //传输关闭信号的管道
	wg       sync.WaitGroup //等待组
}

var services []*service //模块数组，用于保存注册的模块

//注册某个模块
func Register(si Service) {
	s := new(service)               //创建一个模块
	s.si = si                       //保存实现了模块接口的对象si
	s.closeSig = make(chan bool, 1) //创建用于传输关闭信号的管道

	services = append(services, s) //保存模块到模块数组中
}

//初始化所有模块
func Init() {
	for i := 0; i < len(services); i++ {
		services[i].si.OnInit() //先调用各模块的初始化函数
	}

	for i := 0; i < len(services); i++ {
		go run(services[i]) //然后在一个新的goroutine中运行模块
	}
}

//倒序销毁所有模块
func Destroy() {
	for i := len(services) - 1; i >= 0; i-- {
		s := services[i]   //取得索引对应的模块
		s.closeSig <- true //向管道发送关闭信号
		s.wg.Wait()        //等待该模块所在的goroutine执行完成
		destroy(s)         //销毁该模块
	}
}

//运行某个模块
func run(s *service) {
	s.wg.Add(1)          //要等待的goroutine数量加1
	s.si.Run(s.closeSig) //调用模块的运行函数
	s.wg.Done()          //要等待的goroutine数量减1
}

//销毁某个模块
func destroy(s *service) {
	defer func() { //延迟捕获异常
		if r := recover(); r != nil {
			if conf.Env.StackBufLen > 0 {
				buf := make([]byte, conf.Env.StackBufLen) //创建一个字节切片用于存储格式化后的stack trace
				l := runtime.Stack(buf, false)            //格式化调用Stack函数的goroutine的stack trace
				log.Error("%v: %s", r, buf[:l])           //打印错误消息和stack trace
			} else {
				log.Error("%v", r) //只打印错误消息
			}
		}
	}()

	s.si.OnDestroy() //先调用模块的销毁函数，再执行上面的延迟函数
}
