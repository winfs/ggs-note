package login

import (
	"ggs/gate"
	"ggs/log"
	"ggsserver/manager"
	"ggsserver/msg"
)

func handleHelloRequest(args []interface{}) {
	m := args[0].(*msg.HelloRequest) //收到的消息
	a := args[1].(gate.Agent)        //消息的发送者

	log.Info("%s\n", m) //输出收到的消息内容

	//给发送者回应一个消息
	a.WriteMsg(&msg.Hello{
		Name: proto.String("111"),
	})
}
