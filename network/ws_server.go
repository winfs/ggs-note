package network

import (
	"ggs/conf"
	"ggs/log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//WebSocket对象
type WSServer struct {
	ln      net.Listener //网络监听器接口
	handler *WSHandler   //WebSocket处理器对象
}

//WebSocket处理器对象
type WSHandler struct {
	newAgent   func(*WSConn) Agent
	upgrader   websocket.Upgrader
	conns      WebsocketConnSet
	mutexConns sync.Mutex
	wg         sync.WaitGroup
}

func (handler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	conn, err := handler.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Debug("upgrade error: %v", err)
		return
	}
	conn.SetReadLimit(int64(conf.Env.MaxMsgLen))

	handler.wg.Add(1)
	defer handler.wg.Done()

	handler.mutexConns.Lock()
	if handler.conns == nil {
		handler.mutexConns.Unlock()
		conn.Close()
		return
	}
	if len(handler.conns) >= conf.Env.MaxConnNum {
		handler.mutexConns.Unlock()
		conn.Close()
		log.Debug("too many connections")
		return
	}
	handler.conns[conn] = struct{}{}
	handler.mutexConns.Unlock()

	wsConn := newWSConn(conn)
	agent := handler.newAgent(wsConn)
	agent.Run()

	wsConn.Close()
	handler.mutexConns.Lock()
	delete(handler.conns, conn)
	handler.mutexConns.Unlock()
	agent.OnClose()
}

//开始监听
func (server *WSServer) Start(newAgent func(*WSConn) Agent) {
	if newAgent == nil {
		log.Fatal("newAgent must not be nil")
	}

	ln, err := net.Listen("tcp", conf.Env.WSAddr) //建立一个tcp连接
	if err != nil {
		log.Fatal("%v", err)
	}

	server.ln = ln               //保存监听器接口
	server.handler = &WSHandler{ //设置WebSocket处理器对象
		newAgent: newAgent,
		conns:    make(WebsocketConnSet),
		upgrader: websocket.Upgrader{
			HandshakeTimeout: conf.Env.HTTPTimeout,
			CheckOrigin:      func(_ *http.Request) bool { return true },
		},
	}

	//创建一个http服务端
	httpServer := &http.Server{
		Addr:           conf.Env.WSAddr,      //监听的地址
		Handler:        server.handler,       //调用的处理器
		ReadTimeout:    conf.Env.HTTPTimeout, //请求的读取操作在超时前的最大持续时间
		WriteTimeout:   conf.Env.HTTPTimeout, //回复的写入操作在超时前的最大持续时间
		MaxHeaderBytes: 1024,
	}

	go httpServer.Serve(ln) //开启一个goroutine来接收监听地址收到的每一个连接
	log.Info("Websocket host: %v", conf.Env.WSAddr)
}

//关闭连接
func (server *WSServer) Close() {
	server.ln.Close()

	server.handler.mutexConns.Lock()
	for conn := range server.handler.conns {
		conn.Close()
	}
	server.handler.conns = nil
	server.handler.mutexConns.Unlock()

	server.handler.wg.Wait()
}
