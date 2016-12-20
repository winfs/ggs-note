package network

//消息处理器接口
type Processor interface {
	// must goroutine safe
	Route(msg interface{}, userData interface{}) error
	// must goroutine safe
	Unmarshal(data []byte) (interface{}, error) //消息解码
	// must goroutine safe
	Marshal(msg interface{}) ([][]byte, error) //消息编码
}
