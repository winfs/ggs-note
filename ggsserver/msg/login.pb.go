// Code generated by protoc-gen-go.
// source: login.proto
// DO NOT EDIT!

package msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Login struct {
	User             *User  `protobuf:"bytes,1,req,name=user" json:"user,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Login) Reset()                    { *m = Login{} }
func (m *Login) String() string            { return proto.CompactTextString(m) }
func (*Login) ProtoMessage()               {}
func (*Login) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func (m *Login) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type LoginRequest struct {
	Token            *string `protobuf:"bytes,1,req,name=token" json:"token,omitempty"`
	DeviceId         *uint32 `protobuf:"varint,2,req,name=deviceId" json:"deviceId,omitempty"`
	Versiong         *uint32 `protobuf:"varint,3,req,name=versiong" json:"versiong,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LoginRequest) Reset()                    { *m = LoginRequest{} }
func (m *LoginRequest) String() string            { return proto.CompactTextString(m) }
func (*LoginRequest) ProtoMessage()               {}
func (*LoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

func (m *LoginRequest) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *LoginRequest) GetDeviceId() uint32 {
	if m != nil && m.DeviceId != nil {
		return *m.DeviceId
	}
	return 0
}

func (m *LoginRequest) GetVersiong() uint32 {
	if m != nil && m.Versiong != nil {
		return *m.Versiong
	}
	return 0
}

func init() {
	proto.RegisterType((*Login)(nil), "msg.Login")
	proto.RegisterType((*LoginRequest)(nil), "msg.loginRequest")
}

var fileDescriptor2 = []byte{
	// 133 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0xce, 0xc9, 0x4f, 0xcf,
	0xcc, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xce, 0x2d, 0x4e, 0x97, 0xe2, 0x2a, 0x2d,
	0x4e, 0x2d, 0x82, 0x08, 0x28, 0x29, 0x70, 0xb1, 0xfa, 0x80, 0xe4, 0x85, 0xc4, 0xb9, 0x58, 0x40,
	0xc2, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0xdc, 0x46, 0x9c, 0x7a, 0x40, 0x85, 0x7a, 0xa1, 0x40, 0x01,
	0x25, 0x47, 0x2e, 0x1e, 0xb0, 0x09, 0x41, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x42, 0xbc, 0x5c,
	0xac, 0x25, 0xf9, 0xd9, 0xa9, 0x79, 0x60, 0x95, 0x9c, 0x42, 0x02, 0x5c, 0x1c, 0x29, 0xa9, 0x65,
	0x99, 0xc9, 0xa9, 0x9e, 0x29, 0x12, 0x4c, 0x40, 0x11, 0x5e, 0x90, 0x48, 0x59, 0x6a, 0x51, 0x71,
	0x66, 0x7e, 0x5e, 0xba, 0x04, 0x33, 0x48, 0x04, 0x10, 0x00, 0x00, 0xff, 0xff, 0x9c, 0x4e, 0x12,
	0x35, 0x83, 0x00, 0x00, 0x00,
}
