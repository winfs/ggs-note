package login

import (
	"ggsserver/msg"
	"reflect"
)

//向当前模块注册某个消息对应的消息处理函数
func init() {
	skeleton.RegisterChanRPC(reflect.TypeOf(&msg.LoginRequest{}), handleLoginRequest)
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}
