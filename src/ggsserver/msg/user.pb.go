// Code generated by protoc-gen-go.
// source: user.proto
// DO NOT EDIT!

package msg

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// id: 17
type User struct {
	Id               *string      `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Name             *string      `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Avatar           *string      `protobuf:"bytes,3,opt,name=avatar" json:"avatar,omitempty"`
	OfflineTime      *uint32      `protobuf:"varint,4,opt,name=offlineTime" json:"offlineTime,omitempty"`
	CountryName      *string      `protobuf:"bytes,5,opt,name=countryName" json:"countryName,omitempty"`
	KingLisence      *KingLisence `protobuf:"bytes,6,opt,name=kingLisence" json:"kingLisence,omitempty"`
	Charges          []float64    `protobuf:"fixed64,7,rep,name=charges" json:"charges,omitempty"`
	IsNew            *bool        `protobuf:"varint,8,opt,name=isNew" json:"isNew,omitempty"`
	Sex              *uint32      `protobuf:"varint,9,opt,name=sex" json:"sex,omitempty"`
	XXX_unrecognized []byte       `json:"-"`
}

func (m *User) Reset()                    { *m = User{} }
func (m *User) String() string            { return proto.CompactTextString(m) }
func (*User) ProtoMessage()               {}
func (*User) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{0} }

func (m *User) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *User) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *User) GetAvatar() string {
	if m != nil && m.Avatar != nil {
		return *m.Avatar
	}
	return ""
}

func (m *User) GetOfflineTime() uint32 {
	if m != nil && m.OfflineTime != nil {
		return *m.OfflineTime
	}
	return 0
}

func (m *User) GetCountryName() string {
	if m != nil && m.CountryName != nil {
		return *m.CountryName
	}
	return ""
}

func (m *User) GetKingLisence() *KingLisence {
	if m != nil {
		return m.KingLisence
	}
	return nil
}

func (m *User) GetCharges() []float64 {
	if m != nil {
		return m.Charges
	}
	return nil
}

func (m *User) GetIsNew() bool {
	if m != nil && m.IsNew != nil {
		return *m.IsNew
	}
	return false
}

func (m *User) GetSex() uint32 {
	if m != nil && m.Sex != nil {
		return *m.Sex
	}
	return 0
}

type KingLisence struct {
	Name              *string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Have              *uint32  `protobuf:"varint,2,opt,name=have" json:"have,omitempty"`
	Remain            *int32   `protobuf:"varint,3,opt,name=remain" json:"remain,omitempty"`
	GiveWings         *uint32  `protobuf:"varint,4,opt,name=giveWings" json:"giveWings,omitempty"`
	AttackIncr        *float64 `protobuf:"fixed64,5,opt,name=attackIncr" json:"attackIncr,omitempty"`
	SpeedIncr         *float64 `protobuf:"fixed64,6,opt,name=speedIncr" json:"speedIncr,omitempty"`
	OfflineIncomeIncr *float64 `protobuf:"fixed64,7,opt,name=offlineIncomeIncr" json:"offlineIncomeIncr,omitempty"`
	InvadeListIncr    *uint32  `protobuf:"varint,8,opt,name=invadeListIncr" json:"invadeListIncr,omitempty"`
	XXX_unrecognized  []byte   `json:"-"`
}

func (m *KingLisence) Reset()                    { *m = KingLisence{} }
func (m *KingLisence) String() string            { return proto.CompactTextString(m) }
func (*KingLisence) ProtoMessage()               {}
func (*KingLisence) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{1} }

func (m *KingLisence) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *KingLisence) GetHave() uint32 {
	if m != nil && m.Have != nil {
		return *m.Have
	}
	return 0
}

func (m *KingLisence) GetRemain() int32 {
	if m != nil && m.Remain != nil {
		return *m.Remain
	}
	return 0
}

func (m *KingLisence) GetGiveWings() uint32 {
	if m != nil && m.GiveWings != nil {
		return *m.GiveWings
	}
	return 0
}

func (m *KingLisence) GetAttackIncr() float64 {
	if m != nil && m.AttackIncr != nil {
		return *m.AttackIncr
	}
	return 0
}

func (m *KingLisence) GetSpeedIncr() float64 {
	if m != nil && m.SpeedIncr != nil {
		return *m.SpeedIncr
	}
	return 0
}

func (m *KingLisence) GetOfflineIncomeIncr() float64 {
	if m != nil && m.OfflineIncomeIncr != nil {
		return *m.OfflineIncomeIncr
	}
	return 0
}

func (m *KingLisence) GetInvadeListIncr() uint32 {
	if m != nil && m.InvadeListIncr != nil {
		return *m.InvadeListIncr
	}
	return 0
}

type UserInfo struct {
	Name             *string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Avatar           *string  `protobuf:"bytes,2,opt,name=avatar" json:"avatar,omitempty"`
	Speed            *float64 `protobuf:"fixed64,3,opt,name=speed" json:"speed,omitempty"`
	Attack           *float64 `protobuf:"fixed64,4,opt,name=attack" json:"attack,omitempty"`
	AttackCityNum    *uint32  `protobuf:"varint,5,opt,name=attackCityNum" json:"attackCityNum,omitempty"`
	OwnCityNum       *uint32  `protobuf:"varint,6,opt,name=ownCityNum" json:"ownCityNum,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *UserInfo) Reset()                    { *m = UserInfo{} }
func (m *UserInfo) String() string            { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()               {}
func (*UserInfo) Descriptor() ([]byte, []int) { return fileDescriptor5, []int{2} }

func (m *UserInfo) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *UserInfo) GetAvatar() string {
	if m != nil && m.Avatar != nil {
		return *m.Avatar
	}
	return ""
}

func (m *UserInfo) GetSpeed() float64 {
	if m != nil && m.Speed != nil {
		return *m.Speed
	}
	return 0
}

func (m *UserInfo) GetAttack() float64 {
	if m != nil && m.Attack != nil {
		return *m.Attack
	}
	return 0
}

func (m *UserInfo) GetAttackCityNum() uint32 {
	if m != nil && m.AttackCityNum != nil {
		return *m.AttackCityNum
	}
	return 0
}

func (m *UserInfo) GetOwnCityNum() uint32 {
	if m != nil && m.OwnCityNum != nil {
		return *m.OwnCityNum
	}
	return 0
}

func init() {
	proto.RegisterType((*User)(nil), "msg.User")
	proto.RegisterType((*KingLisence)(nil), "msg.KingLisence")
	proto.RegisterType((*UserInfo)(nil), "msg.UserInfo")
}

func init() { proto.RegisterFile("user.proto", fileDescriptor5) }

var fileDescriptor5 = []byte{
	// 384 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x4f, 0x8b, 0xd4, 0x40,
	0x10, 0xc5, 0xe9, 0xfc, 0x9b, 0x4c, 0x85, 0x2c, 0x6b, 0xb1, 0x48, 0x1f, 0x44, 0xc2, 0x20, 0x92,
	0x83, 0xcc, 0x61, 0xbf, 0x82, 0xa7, 0xc1, 0x65, 0x0e, 0x8d, 0xe2, 0xb9, 0x49, 0x6a, 0xb2, 0xcd,
	0x9a, 0xce, 0xd2, 0x9d, 0xc9, 0xba, 0x5f, 0xc9, 0xcf, 0x28, 0x28, 0xa9, 0x64, 0x4d, 0x46, 0xbd,
	0xd5, 0xfb, 0xf5, 0x4b, 0x53, 0xef, 0x75, 0x00, 0xce, 0x9e, 0xdc, 0xfe, 0xd1, 0x75, 0x7d, 0x87,
	0x61, 0xeb, 0x9b, 0xdd, 0x4f, 0x01, 0xd1, 0x17, 0x4f, 0x0e, 0xaf, 0x20, 0x30, 0xb5, 0x14, 0x85,
	0x28, 0xb7, 0x2a, 0x30, 0x35, 0x22, 0x44, 0x56, 0xb7, 0x24, 0x03, 0x26, 0x3c, 0xe3, 0x6b, 0x48,
	0xf4, 0xa0, 0x7b, 0xed, 0x64, 0xc8, 0x74, 0x56, 0x58, 0x40, 0xd6, 0x9d, 0x4e, 0xdf, 0x8c, 0xa5,
	0xcf, 0xa6, 0x25, 0x19, 0x15, 0xa2, 0xcc, 0xd5, 0x1a, 0x8d, 0x8e, 0xaa, 0x3b, 0xdb, 0xde, 0x3d,
	0x1f, 0xc7, 0x4b, 0x63, 0xfe, 0x7c, 0x8d, 0xf0, 0x16, 0xb2, 0x07, 0x63, 0x9b, 0x3b, 0xe3, 0xc9,
	0x56, 0x24, 0x93, 0x42, 0x94, 0xd9, 0xed, 0xf5, 0xbe, 0xf5, 0xcd, 0xfe, 0xd3, 0xc2, 0xd5, 0xda,
	0x84, 0x12, 0x36, 0xd5, 0xbd, 0x76, 0x0d, 0x79, 0xb9, 0x29, 0xc2, 0x52, 0xa8, 0x17, 0x89, 0x37,
	0x10, 0x1b, 0x7f, 0xa4, 0x27, 0x99, 0x16, 0xa2, 0x4c, 0xd5, 0x24, 0xf0, 0x1a, 0x42, 0x4f, 0xdf,
	0xe5, 0x96, 0xf7, 0x1b, 0xc7, 0xdd, 0x2f, 0x01, 0xd9, 0xea, 0xfa, 0x3f, 0xa9, 0xc5, 0x2a, 0x35,
	0x42, 0x74, 0xaf, 0x87, 0xa9, 0x89, 0x5c, 0xf1, 0x3c, 0x36, 0xe1, 0xa8, 0xd5, 0xc6, 0x72, 0x13,
	0xb1, 0x9a, 0x15, 0xbe, 0x81, 0x6d, 0x63, 0x06, 0xfa, 0x6a, 0x6c, 0xe3, 0xe7, 0x1e, 0x16, 0x80,
	0x6f, 0x01, 0x74, 0xdf, 0xeb, 0xea, 0xe1, 0x60, 0x2b, 0xc7, 0x25, 0x08, 0xb5, 0x22, 0xe3, 0xd7,
	0xfe, 0x91, 0xa8, 0xe6, 0xe3, 0x84, 0x8f, 0x17, 0x80, 0x1f, 0xe0, 0xd5, 0x5c, 0xe9, 0xc1, 0x56,
	0x5d, 0x4b, 0xec, 0xda, 0xb0, 0xeb, 0xdf, 0x03, 0x7c, 0x0f, 0x57, 0xc6, 0x0e, 0xba, 0xa6, 0x3b,
	0xe3, 0x7b, 0xb6, 0xa6, 0xbc, 0xce, 0x5f, 0x74, 0xf7, 0x43, 0x40, 0x3a, 0xfe, 0x00, 0x07, 0x7b,
	0xea, 0xfe, 0x1b, 0x7f, 0x79, 0xf4, 0xe0, 0xe2, 0xd1, 0x6f, 0x20, 0xe6, 0xdd, 0xb8, 0x01, 0xa1,
	0x26, 0xc1, 0x6e, 0x0e, 0xc4, 0xe9, 0x85, 0x9a, 0x15, 0xbe, 0x83, 0x7c, 0x9a, 0x3e, 0x9a, 0xfe,
	0xf9, 0x78, 0x6e, 0x39, 0x7d, 0xae, 0x2e, 0xe1, 0x58, 0x50, 0xf7, 0x64, 0x5f, 0x2c, 0x09, 0x5b,
	0x56, 0xe4, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x74, 0xea, 0xab, 0xbf, 0xbf, 0x02, 0x00, 0x00,
}
