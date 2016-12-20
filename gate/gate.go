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

type agent struct {
	conn     network.Conn //连接对象接口
	gate     *Gate        //网关引用
	userData interface{}  //任意数据
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		if a.gate.Processor != nil {
			msg, err := a.gate.Processor.Unmarshal(data)
			if err != nil {
				log.Debug("unmarshal message error: %v", err)
				break
			}
			err = a.gate.Processor.Route(msg, a)
			if err != nil {
				log.Debug("route message error: %v", err)
				break
			}
		}
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Open(0).Call0("CloseAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(msg interface{}) {
	if a.gate.Processor != nil {
		data, err := a.gate.Processor.Marshal(msg)
		if err != nil {
			log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
		err = a.conn.WriteMsg(data...)
		if err != nil {
			log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
			return
		}
	}
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
