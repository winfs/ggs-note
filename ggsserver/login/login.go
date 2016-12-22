package login

import (
	"ggs/gate"
	"ggs/log"
	"ggs/service"
	"ggsserver/msg"

	"github.com/golang/protobuf/proto"
)

var skeleton = service.NewSkeleton()

var (
	Service       = new(Login)
	ChanRPCServer = skeleton.ChanRPCServer()
)

type Login struct {
	*service.Skeleton
}

func (l *Login) OnInit() {
	l.Skeleton = skeleton
}

func (l *Login) OnDestroy() {
	log.Info("login service destoryed.")
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	log.Info("--------------- agent close -------------")
}

func handleHelloRequest(args []interface{}) {
	m := args[0].(*msg.HelloRequest) //收到的消息
	a := args[1].(gate.Agent)        //消息的发送者

	log.Info("%s\n", m) //输出收到的消息内容

	//给发送者回应一个消息
	a.WriteMsg(&msg.Hello{
		Name: proto.String("111"),
	})
}
