package msg

import (
	"ggs/network/protobuf"
)

// 使用protobuf消息处理器
var Processor = protobuf.NewProcessor()

// 注册消息
func init() {
	// response
	Processor.Register(&Hello{})  // 0
	Processor.Register(&Status{}) // 1
	Processor.Register(&Login{})  // 2
	Processor.Register(&Error{})  // 16
	Processor.Register(&User{})   // 17

	// request
	Processor.Register(&LoginRequest{}) // 5
}
