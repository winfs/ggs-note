package game

import (
	"ggsserver/msg"
	"reflect"
)

func init() {
	skeleton.RegisterChanRPC(reflect.TypeOf(&msg.Hello{}), handleHello)
}
