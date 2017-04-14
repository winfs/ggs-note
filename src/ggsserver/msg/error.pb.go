// Code generated by protoc-gen-go.
// source: error.proto
// DO NOT EDIT!

/*
Package msg is a generated protocol buffer package.

It is generated from these files:
	error.proto
	hello.proto
	login.proto
	spare.proto
	status.proto
	user.proto

It has these top-level messages:
	Error
	Hello
	Login
	LoginRequest
	Spare0
	Spare1
	Spare2
	Spare3
	Spare4
	Spare5
	Spare6
	Spare7
	Spare8
	Spare9
	Status
	User
	KingLisence
	UserInfo
*/
package msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// id: 16
type Error struct {
	Code             *uint32 `protobuf:"varint,1,req,name=code" json:"code,omitempty"`
	Desc             *string `protobuf:"bytes,2,opt,name=desc" json:"desc,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Error) GetCode() uint32 {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return 0
}

func (m *Error) GetDesc() string {
	if m != nil && m.Desc != nil {
		return *m.Desc
	}
	return ""
}

func init() {
	proto.RegisterType((*Error)(nil), "msg.Error")
}

func init() { proto.RegisterFile("error.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 80 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x2d, 0x2a, 0xca,
	0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xce, 0x2d, 0x4e, 0x57, 0xd2, 0xe7, 0x62,
	0x75, 0x05, 0x89, 0x09, 0x09, 0x71, 0xb1, 0x24, 0xe7, 0xa7, 0xa4, 0x4a, 0x30, 0x2a, 0x30, 0x69,
	0xf0, 0x06, 0x81, 0xd9, 0x20, 0xb1, 0x94, 0xd4, 0xe2, 0x64, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce,
	0x20, 0x30, 0x1b, 0x10, 0x00, 0x00, 0xff, 0xff, 0x02, 0xd3, 0x12, 0x2a, 0x43, 0x00, 0x00, 0x00,
}
