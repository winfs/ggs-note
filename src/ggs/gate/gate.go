//网关模块，负责游戏客户端的接入
package gate

import (
	"ggs/chanrpc"
	"ggs/conf"
	"ggs/log"
	"ggs/network"
	"reflect"
)

var Service = new(Gate)

//网关对象
type Gate struct {
	Processor    network.Processor //消息处理器
	AgentChanRPC *chanrpc.Server   //rpc服务器引用
}

//初始化
func (gate *Gate) OnInit() {}

//运行
func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if conf.Env.WSAddr != "" { //如果配置中有设置WebSocket监听地址
		wsServer = new(network.WSServer) //创建一个WebSocket服务器
	}

	if wsServer != nil {
		//开始监听
		wsServer.Start(func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewAgent", a) //调用NewAgent并传递参数a
			}
			return a
		})
	}
	<-closeSig //阻塞，直到接收到一个退出标志
	if wsServer != nil {
		wsServer.Close() //关闭连接
	}
}

//销毁
func (gate *Gate) OnDestroy() {
	log.Debug("gate service destoryed.")
}

// 代理对象
type agent struct {
	conn     network.Conn //连接对象接口
	gate     *Gate        //网关对象引用
	userData interface{}  //任意数据
}

func (a *agent) Run() {
	for { // 一直循环
		data, err := a.conn.ReadMsg() // 读取消息
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil { // 网关对象的消息处理器不为空
			msg, err := a.gate.Processor.Unmarshal(data) // 解码消息
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a) //
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil { //rpc服务器不为空
		err := a.gate.AgentChanRPC.Open(0).Call0("CloseAgent", a) //
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

// 给客户端回应消息
func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil { // 网关对象的消息处理器不为空
		data, err := a.gate.Processor.Marshal(msg) // 编码消息
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...) // 通过具体连接对象写入消息
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
	}
}

// 关闭连接
func (a *agent) Close() {
	a.conn.Close() // 关闭连接
}

// 返回数据
func (a *agent) UserData() interface{} {
	return a.userData
}

// 设置数据
func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
