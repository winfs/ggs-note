// 定时器相关
package timer

import (
	"ggs/conf"
	"ggs/log"
	"runtime"
	"time"
)

// one dispatcher per goroutine (goroutine not safe)

// 定时器
type Dispatcher struct {
	ChanTimer chan *Timer // 用于传输Timer信息的管道
}

// 创建定时器
func NewDispatcher(l int) *Dispatcher {
	disp := new(Dispatcher)               // 创建定时器
	disp.ChanTimer = make(chan *Timer, l) // 创建用于传递Timer信息的管道
	return disp
}

type Timer struct {
	t  *time.Timer // time包的Timer类型引用
	cb func()      // 回调函数
}

// 停止
func (t *Timer) Stop() {
	t.t.Stop() // 停止time.Timer的执行
	t.cb = nil // 清空回调函数
}

// 执行回调函数
func (t *Timer) Cb() {
	defer func() { // 延迟捕获异常
		t.cb = nil // 清空回调函数
		if r := recover(); r != nil {
			if conf.Env.StackBufLen > 0 {
				buf := make([]byte, conf.Env.StackBufLen)
				l := runtime.Stack(buf, false)
				log.Error("%v: %s", r, buf[:l])
			} else {
				log.Error("%v", r)
			}
		}
	}()

	if t.cb != nil {
		t.cb() // 执行回调函数
	}
}

// AfterFunc会等待d时长后调用f函数，这里的f函数将在AfterFunc的调用goroutine中执行
func (disp *Dispatcher) AfterFunc(d time.Duration, cb func()) *Timer {
	t := new(Timer)                  // 创建Timer
	t.cb = cb                        // 保存回调函数
	t.t = time.AfterFunc(d, func() { // 另起一个goroutine等待时间段d过去后调用func
		disp.ChanTimer <- t // 将t发送到定时器的Timer信息管道中
	})
	return t
}

type Cron struct {
	t *Timer // Timer类型引用
}

// 停止Cron
func (c *Cron) Stop() {
	if c.t != nil {
		c.t.Stop()
	}
}

// 支持cron表达式，用于实现诸如'每天9点执行'、'每周末6点执行'的逻辑
func (disp *Dispatcher) CronFunc(cronExpr *CronExpr, _cb func()) *Cron {
	c := new(Cron)

	now := time.Now()
	nextTime := cronExpr.Next(now)
	if nextTime.IsZero() {
		return c
	}

	// callback
	var cb func()
	cb = func() {
		defer _cb()

		now := time.Now()
		nextTime := cronExpr.Next(now)
		if nextTime.IsZero() {
			return
		}
		c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	}

	c.t = disp.AfterFunc(nextTime.Sub(now), cb)
	return c
}
