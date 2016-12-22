package msg

import (
	"ggs/network/protobuf"
)

// 使用protobuf消息处理器
var Processor = protobuf.NewProcessor()

// 注册消息
func init() {
	// response
	Processor.Register(&Hello{}) //0

	// request
	Processor.Register(&HelloRequest{}) //1
}
