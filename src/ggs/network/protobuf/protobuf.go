package protobuf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"ggs/chanrpc"
	"ggs/log"

	"github.com/golang/protobuf/proto"
)

// -------------------------
// | id | protobuf message |
// -------------------------
type Processor struct {
	littleEndian bool                    // 字节排序
	msgInfo      []*MsgInfo              // 注册的消息集合
	msgID        map[reflect.Type]uint16 // 消息类型 => ID的映射
}

type MsgInfo struct {
	msgType    reflect.Type    // 消息类型
	msgRouter  *chanrpc.Server // rpc模块
	msgHandler MsgHandler      // 消息处理函数
}

type MsgHandler func([]interface{})

// 初始化
func NewProcessor() *Processor {
	p := new(Processor)
	p.littleEndian = false
	p.msgID = make(map[reflect.Type]uint16)
	return p
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
func (p *Processor) SetByteOrder(littleEndian bool) {
	p.littleEndian = littleEndian
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
// 注册proto消息
func (p *Processor) Register(msg proto.Message) {
	msgType := reflect.TypeOf(msg) // 判断消息类型
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal("protobuf message pointer required")
	}
	if _, ok := p.msgID[msgType]; ok { // 判断消息是否已被注册
		log.Fatal("message %s is already registered", msgType)
	}
	if len(p.msgInfo) >= math.MaxUint16 { // 判断是否已超出最大的注册消息个数
		log.Fatal("too many protobuf messages (max = %v)", math.MaxUint16)
	}

	// 追加一条消息
	i := new(MsgInfo)
	i.msgType = msgType
	p.msgInfo = append(p.msgInfo, i)
	p.msgID[msgType] = uint16(len(p.msgInfo) - 1) // 初始消息ID为0，之后每加一条消息ID加1 ==> 与msgInfo索引保持一致
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
// 设置消息交由哪个rpc模块来处理
func (p *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatal("message %s not registered", msgType)
	}

	p.msgInfo[id].msgRouter = msgRouter // 设置当前消息的rpc模块
}

// It's dangerous to call the method on routing or marshaling (unmarshaling)
// 设置消息交由哪个消息处理函数来处理
func (p *Processor) SetHandler(msg proto.Message, msgHandler MsgHandler) {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		log.Fatal("message %s not registered", msgType)
	}

	p.msgInfo[id].msgHandler = msgHandler // 设置当前消息的消息处理函数
}

// goroutine safe
// 消息的路由分发(将消息交由rpc模块或者消息处理函数来处理)
func (p *Processor) Route(msg interface{}, userData interface{}) error {
	msgType := reflect.TypeOf(msg)
	id, ok := p.msgID[msgType]
	if !ok {
		return fmt.Errorf("message %s not registered", msgType)
	}

	i := p.msgInfo[id] // 获取当前消息
	if i.msgHandler != nil {
		i.msgHandler([]interface{}{msg, userData})
	}
	if i.msgRouter != nil {
		i.msgRouter.Go(msgType, msg, userData)
	}
	return nil
}

// goroutine safe
// 消息解码
func (p *Processor) Unmarshal(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}

	// id
	var id uint16
	if p.littleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}

	// msg
	if id >= uint16(len(p.msgInfo)) {
		return nil, fmt.Errorf("message id %v not registered", id)
	}
	msg := reflect.New(p.msgInfo[id].msgType.Elem()).Interface()
	return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
}

// goroutine safe
// 消息编码
func (p *Processor) Marshal(msg interface{}) ([][]byte, error) {
	msgType := reflect.TypeOf(msg)

	// id
	_id, ok := p.msgID[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		return nil, err
	}

	id := make([]byte, 2)
	if p.littleEndian {
		binary.LittleEndian.PutUint16(id, _id)
	} else {
		binary.BigEndian.PutUint16(id, _id)
	}

	// data
	data, err := proto.Marshal(msg.(proto.Message))
	return [][]byte{id, data}, err
}

// goroutine safe
func (p *Processor) Range(f func(id uint16, t reflect.Type)) {
	for id, i := range p.msgInfo {
		f(uint16(id), i.msgType)
	}
}
