package login

import (
	"fmt"

	"ggs/gate"
	"ggs/log"
	"ggs/service"

	"ggsserver/db"
	"ggsserver/manager"
	"ggsserver/msg"

	"github.com/golang/protobuf/proto"
	"github.com/mediocregopher/radix.v2/redis"
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
	a := args[0].(gate.Agent)
	dm, ok := a.UserData().(*manager.DataManager)
	if ok {
		log.Info("---------- close user: %v ------------", dm.UM.UserName())
		dm.Clear()
	} else {
		log.Info("--------------- no data manager ------------")
	}
}

func handleLoginRequest(args []interface{}) {
	var err error
	var DM *manager.DataManager
	defer func() {
		if err != nil {
			log.Error("create data manager error: %v", err)
			if DM != nil {
				DM.ResponseErr(4, "创建DataManager失败", 24)
			}
		}
	}()

	m := args[0].(*msg.LoginRequest) // 收到的登录请求消息
	a := args[1].(gate.Agent)        // 消息的发送者
	deviceId := m.GetDeviceId()
	token := m.GetToken()
	log.Info("++++++++++++++++ receive login request token: %v", token)

	userId, err := db.Redis.Cmd("HGET", "tokens", token).Str() // 根据token获取到userId
	if err != nil {
		err = fmt.Errorf("no corresponding user id for token: %v", err)
		return
	}

	shielded, err := checkUserShielded(userId) // 检查用户是否禁止登录(黑名单用户)
	if err != nil {
		return
	}
	if shielded {
		responseError(a, 7, userId) // 如果用户禁止登录，则给客户端回复错误消息,并返回
		return
	}

	DM, err = manager.NewDataManager(userId, a, deviceId) // 初始化DataManager
	if err != nil {
		return
	}
	a.SetUserData(DM)

	DM.Response() // 给客户端回复登录响应消息
	//DM.CheckResidentNotice()
	log.Info("++++++++++++++++ user login, token: %v, user id: %v, name: %v", token, userId, DM.UM.User.Name)
}

// 检查用户是否禁止登录
func checkUserShielded(userId string) (shielded bool, err error) {
	res := db.Redis.Cmd("HGET", "user:"+userId+":info", "shielded")
	if res.IsType(redis.Nil) {
		return
	}
	num, err := res.Int()
	if err != nil {
		return
	}
	shielded = num == 1
	return
}

// 给客户端回复错误消息
func responseError(a gate.Agent, code uint32, desc string) {
	log.Info("send error, code: %v, desc: %v", code, desc)
	a.WriteMsg(&msg.Error{
		Code: proto.Uint32(code),
		Desc: proto.String(desc),
	})
}
