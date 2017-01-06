package network

import (
	"errors"
	"net"
	"sync"

	"ggs/conf"
	"ggs/log"

	"github.com/gorilla/websocket"
)

type WebsocketConnSet map[*websocket.Conn]struct{}

type WSConn struct {
	sync.Mutex
	conn      *websocket.Conn //WebSocket连接对象
	writeChan chan []byte     //传递消息的管道(消息类型为byte[])
	maxMsgLen uint32          //可以传递的最大消息长度
	closeFlag bool            //连接关闭标志
}

//创建WSConn对象
func newWSConn(conn *websocket.Conn) *WSConn {
	wsConn := new(WSConn)                                          //创建一个WSConn对象
	wsConn.conn = conn                                             //保存WebSocket的连接对象
	wsConn.writeChan = make(chan []byte, conf.Env.PendingWriteNum) //创建一个用于传递消息的管道
	wsConn.maxMsgLen = conf.Env.MaxMsgLen                          //设置可以传递的最大消息长度

	//开启一个新的goroutine,读取消息管道中的数据，并写入消息
	go func() {
		for b := range wsConn.writeChan { // 遍历消息管道的内容
			if b == nil {
				break
			}

			err := conn.WriteMessage(websocket.BinaryMessage, b) // 写入二进制数据消息
			if err != nil {
				break
			}
		}

		conn.Close() //关闭WebSocket连接
		wsConn.Lock()
		wsConn.closeFlag = true //设置连接关闭标志
		wsConn.Unlock()
	}()

	return wsConn
}

//关闭连接
func (wsConn *WSConn) destroy() {
	wsConn.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	wsConn.conn.Close()

	if !wsConn.closeFlag {
		close(wsConn.writeChan)
		wsConn.closeFlag = true
	}
}

//关闭连接
func (wsConn *WSConn) Destroy() {
	wsConn.Lock()
	defer wsConn.Unlock()

	wsConn.destroy()
}

//设置连接关闭标志
func (wsConn *WSConn) Close() {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return
	}

	wsConn.write(nil)
	wsConn.closeFlag = true
}

//将消息发送到消息管道中
func (wsConn *WSConn) write(b []byte) {
	if len(wsConn.writeChan) == cap(wsConn.writeChan) {
		log.Debug("close conn: channel full")
		wsConn.destroy()
		return
	}

	wsConn.writeChan <- b //将消息发送到消息管道中
}

//返回本地地址
func (wsConn *WSConn) LocalAddr() net.Addr {
	return wsConn.conn.LocalAddr()
}

//返回远程地址
func (wsConn *WSConn) RemoteAddr() net.Addr {
	return wsConn.conn.RemoteAddr()
}

// goroutine not safe

//读取消息
func (wsConn *WSConn) ReadMsg() ([]byte, error) {
	_, b, err := wsConn.conn.ReadMessage()
	return b, err
}

// args must not be modified by the others goroutines

//写入消息
func (wsConn *WSConn) WriteMsg(args ...[]byte) error {
	wsConn.Lock()
	defer wsConn.Unlock()
	if wsConn.closeFlag {
		return nil
	}

	// get len
	var msgLen uint32
	for i := 0; i < len(args); i++ {
		msgLen += uint32(len(args[i]))
	}

	// check len
	if msgLen > wsConn.maxMsgLen {
		return errors.New("message too long")
	} else if msgLen < 1 {
		return errors.New("message too short")
	}

	// don't copy
	if len(args) == 1 {
		wsConn.write(args[0])
		return nil
	}

	// merge the args
	msg := make([]byte, msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		copy(msg[l:], args[i])
		l += len(args[i])
	}

	wsConn.write(msg)

	return nil
}
