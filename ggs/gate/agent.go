package gate

// 代理接口
type Agent interface {
	WriteMsg(msg interface{})     // 给客户端回应消息
	Close()                       // 关闭连接
	UserData() interface{}        // 返回数据
	SetUserData(data interface{}) // 设置数据
}
