package network

import (
	"net"
)

//连接对象接口
type Conn interface {
	ReadMsg() ([]byte, error)      //读取消息
	WriteMsg(args ...[]byte) error //写入消息
	LocalAddr() net.Addr           //本地地址
	RemoteAddr() net.Addr          //远程地址
	Close()                        //关闭连接
	Destroy()                      //退出销毁
}
