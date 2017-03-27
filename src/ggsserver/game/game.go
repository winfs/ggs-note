package game

import (
	"ggs/gate"
	"ggs/log"
	"ggs/service"

	"ggsserver/manager"
	"ggsserver/msg"

	"github.com/golang/protobuf/proto"
)

var (
	skeleton      = service.NewSkeleton()
	Service       = new(Game)
	ChanRPCServer = skeleton.ChanRPCServer()
)

type Game struct {
	*service.Skeleton
}

func init() {
	manager.GameSkeleton = skeleton
}

func (g *Game) OnInit() {
	g.Skeleton = skeleton
}

func (g *Game) OnDestroy() {
	log.Info("game service destoryed.")
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	_ = a
}

func handleHello(args []interface{}) {
	m := args[0].(*msg.Hello)
	a := args[1].(gate.Agent)
	a.WriteMsg(&msg.Status{
		Code:        proto.Uint32(0),
		Description: proto.String(m.GetName()),
		ByAction:    proto.Uint32(0),
	})
}
