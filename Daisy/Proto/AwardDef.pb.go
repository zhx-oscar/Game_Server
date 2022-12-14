// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: AwardDef.proto

package Proto

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// 发送给客户端离线收益
type RspOfflineAwardData struct {
	OfflineTime       int64             `protobuf:"varint,1,opt,name=OfflineTime,proto3" json:"OfflineTime,omitempty"`
	OfflineAwardDatas *OfflineAwardData `protobuf:"bytes,2,opt,name=OfflineAwardDatas,proto3" json:"OfflineAwardDatas,omitempty"`
}

func (m *RspOfflineAwardData) Reset()         { *m = RspOfflineAwardData{} }
func (m *RspOfflineAwardData) String() string { return proto.CompactTextString(m) }
func (*RspOfflineAwardData) ProtoMessage()    {}
func (*RspOfflineAwardData) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce3b55a9d36e511c, []int{0}
}
func (m *RspOfflineAwardData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RspOfflineAwardData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RspOfflineAwardData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RspOfflineAwardData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RspOfflineAwardData.Merge(m, src)
}
func (m *RspOfflineAwardData) XXX_Size() int {
	return m.Size()
}
func (m *RspOfflineAwardData) XXX_DiscardUnknown() {
	xxx_messageInfo_RspOfflineAwardData.DiscardUnknown(m)
}

var xxx_messageInfo_RspOfflineAwardData proto.InternalMessageInfo

func (m *RspOfflineAwardData) GetOfflineTime() int64 {
	if m != nil {
		return m.OfflineTime
	}
	return 0
}

func (m *RspOfflineAwardData) GetOfflineAwardDatas() *OfflineAwardData {
	if m != nil {
		return m.OfflineAwardDatas
	}
	return nil
}

// 离线收益统计
type OfflineAwardData struct {
	ActorExp          uint32              `protobuf:"varint,1,opt,name=ActorExp,proto3" json:"ActorExp,omitempty"`
	SpecialAgentExp   uint32              `protobuf:"varint,2,opt,name=SpecialAgentExp,proto3" json:"SpecialAgentExp,omitempty"`
	Money             uint32              `protobuf:"varint,3,opt,name=Money,proto3" json:"Money,omitempty"`
	OfflineAwardItems []*OfflineAwardItem `protobuf:"bytes,4,rep,name=OfflineAwardItems,proto3" json:"OfflineAwardItems,omitempty"`
}

func (m *OfflineAwardData) Reset()         { *m = OfflineAwardData{} }
func (m *OfflineAwardData) String() string { return proto.CompactTextString(m) }
func (*OfflineAwardData) ProtoMessage()    {}
func (*OfflineAwardData) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce3b55a9d36e511c, []int{1}
}
func (m *OfflineAwardData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OfflineAwardData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OfflineAwardData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OfflineAwardData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OfflineAwardData.Merge(m, src)
}
func (m *OfflineAwardData) XXX_Size() int {
	return m.Size()
}
func (m *OfflineAwardData) XXX_DiscardUnknown() {
	xxx_messageInfo_OfflineAwardData.DiscardUnknown(m)
}

var xxx_messageInfo_OfflineAwardData proto.InternalMessageInfo

func (m *OfflineAwardData) GetActorExp() uint32 {
	if m != nil {
		return m.ActorExp
	}
	return 0
}

func (m *OfflineAwardData) GetSpecialAgentExp() uint32 {
	if m != nil {
		return m.SpecialAgentExp
	}
	return 0
}

func (m *OfflineAwardData) GetMoney() uint32 {
	if m != nil {
		return m.Money
	}
	return 0
}

func (m *OfflineAwardData) GetOfflineAwardItems() []*OfflineAwardItem {
	if m != nil {
		return m.OfflineAwardItems
	}
	return nil
}

// 离线收益道具
type OfflineAwardItem struct {
	ID       string        `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Type     ItemEnum_Type `protobuf:"varint,2,opt,name=Type,proto3,enum=ItemEnum_Type" json:"Type,omitempty"`
	Num      uint32        `protobuf:"varint,3,opt,name=Num,proto3" json:"Num,omitempty"`
	ConfigID uint32        `protobuf:"varint,4,opt,name=ConfigID,proto3" json:"ConfigID,omitempty"`
}

func (m *OfflineAwardItem) Reset()         { *m = OfflineAwardItem{} }
func (m *OfflineAwardItem) String() string { return proto.CompactTextString(m) }
func (*OfflineAwardItem) ProtoMessage()    {}
func (*OfflineAwardItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_ce3b55a9d36e511c, []int{2}
}
func (m *OfflineAwardItem) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *OfflineAwardItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_OfflineAwardItem.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *OfflineAwardItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OfflineAwardItem.Merge(m, src)
}
func (m *OfflineAwardItem) XXX_Size() int {
	return m.Size()
}
func (m *OfflineAwardItem) XXX_DiscardUnknown() {
	xxx_messageInfo_OfflineAwardItem.DiscardUnknown(m)
}

var xxx_messageInfo_OfflineAwardItem proto.InternalMessageInfo

func (m *OfflineAwardItem) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *OfflineAwardItem) GetType() ItemEnum_Type {
	if m != nil {
		return m.Type
	}
	return ItemEnum_ZeroNoUse
}

func (m *OfflineAwardItem) GetNum() uint32 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *OfflineAwardItem) GetConfigID() uint32 {
	if m != nil {
		return m.ConfigID
	}
	return 0
}

func init() {
	proto.RegisterType((*RspOfflineAwardData)(nil), "RspOfflineAwardData")
	proto.RegisterType((*OfflineAwardData)(nil), "OfflineAwardData")
	proto.RegisterType((*OfflineAwardItem)(nil), "OfflineAwardItem")
}

func init() { proto.RegisterFile("AwardDef.proto", fileDescriptor_ce3b55a9d36e511c) }

var fileDescriptor_ce3b55a9d36e511c = []byte{
	// 300 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x4d, 0x4e, 0x83, 0x40,
	0x14, 0xc7, 0x19, 0xa0, 0x46, 0x5f, 0x23, 0xb6, 0xa3, 0x0b, 0xd2, 0xc5, 0x48, 0x58, 0xb1, 0x62,
	0x51, 0x0f, 0x60, 0xa8, 0x74, 0xc1, 0xc2, 0x8f, 0x8c, 0x5d, 0xb9, 0xc3, 0x3a, 0x34, 0x24, 0x05,
	0x26, 0x30, 0x8d, 0xed, 0x2d, 0xbc, 0x89, 0xd7, 0x70, 0xd9, 0xa5, 0x4b, 0x03, 0x17, 0x31, 0x33,
	0x54, 0xa3, 0xb0, 0xe3, 0xff, 0x41, 0xde, 0xef, 0xcd, 0x03, 0x2b, 0x78, 0x8d, 0xcb, 0x97, 0x90,
	0x25, 0x3e, 0x2f, 0x0b, 0x51, 0x4c, 0xac, 0x48, 0xb0, 0x6c, 0x16, 0x57, 0xac, 0xd5, 0xee, 0x16,
	0xce, 0x69, 0xc5, 0xef, 0x93, 0x64, 0x9d, 0xe6, 0xac, 0xed, 0xc6, 0x22, 0xc6, 0x0e, 0x0c, 0x0f,
	0xde, 0x22, 0xcd, 0x98, 0x8d, 0x1c, 0xe4, 0x19, 0xf4, 0xaf, 0x85, 0xaf, 0x61, 0xdc, 0xfd, 0xab,
	0xb2, 0x75, 0x07, 0x79, 0xc3, 0xe9, 0xd8, 0xef, 0x26, 0xb4, 0xdf, 0x75, 0xdf, 0x11, 0x8c, 0x7a,
	0x73, 0x27, 0x70, 0x1c, 0x2c, 0x45, 0x51, 0xce, 0xb7, 0x5c, 0x0d, 0x3d, 0xa5, 0xbf, 0x1a, 0x7b,
	0x70, 0xf6, 0xc8, 0xd9, 0x32, 0x8d, 0xd7, 0xc1, 0x8a, 0xe5, 0x42, 0x56, 0x74, 0x55, 0xe9, 0xda,
	0xf8, 0x02, 0x06, 0xb7, 0x45, 0xce, 0x76, 0xb6, 0xa1, 0xf2, 0x56, 0x74, 0x89, 0xe5, 0x43, 0x54,
	0xb6, 0xe9, 0x18, 0x3d, 0x62, 0x99, 0xd0, 0x7e, 0xd7, 0x15, 0xff, 0x81, 0xa5, 0x89, 0x2d, 0xd0,
	0xa3, 0x50, 0xa1, 0x9e, 0x50, 0x3d, 0x0a, 0xb1, 0x0b, 0xe6, 0x62, 0xc7, 0x99, 0x22, 0xb3, 0xa6,
	0x96, 0x2f, 0x4b, 0xf3, 0x7c, 0x93, 0xf9, 0xd2, 0xa5, 0x2a, 0xc3, 0x23, 0x30, 0xee, 0x36, 0xd9,
	0x01, 0x4e, 0x7e, 0xca, 0xb5, 0x6f, 0x8a, 0x3c, 0x49, 0x57, 0x51, 0x68, 0x9b, 0xed, 0xda, 0x3f,
	0x7a, 0x76, 0xf9, 0x51, 0x13, 0xb4, 0xaf, 0x09, 0xfa, 0xaa, 0x09, 0x7a, 0x6b, 0x88, 0xb6, 0x6f,
	0x88, 0xf6, 0xd9, 0x10, 0xed, 0x69, 0xf0, 0x20, 0x4f, 0xf8, 0x7c, 0xa4, 0x2e, 0x79, 0xf5, 0x1d,
	0x00, 0x00, 0xff, 0xff, 0xf8, 0x38, 0x2f, 0x1b, 0xeb, 0x01, 0x00, 0x00,
}

func (m *RspOfflineAwardData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RspOfflineAwardData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RspOfflineAwardData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.OfflineAwardDatas != nil {
		{
			size, err := m.OfflineAwardDatas.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAwardDef(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.OfflineTime != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.OfflineTime))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OfflineAwardData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OfflineAwardData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OfflineAwardData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.OfflineAwardItems) > 0 {
		for iNdEx := len(m.OfflineAwardItems) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.OfflineAwardItems[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintAwardDef(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Money != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.Money))
		i--
		dAtA[i] = 0x18
	}
	if m.SpecialAgentExp != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.SpecialAgentExp))
		i--
		dAtA[i] = 0x10
	}
	if m.ActorExp != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.ActorExp))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *OfflineAwardItem) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *OfflineAwardItem) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *OfflineAwardItem) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ConfigID != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.ConfigID))
		i--
		dAtA[i] = 0x20
	}
	if m.Num != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.Num))
		i--
		dAtA[i] = 0x18
	}
	if m.Type != 0 {
		i = encodeVarintAwardDef(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x10
	}
	if len(m.ID) > 0 {
		i -= len(m.ID)
		copy(dAtA[i:], m.ID)
		i = encodeVarintAwardDef(dAtA, i, uint64(len(m.ID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAwardDef(dAtA []byte, offset int, v uint64) int {
	offset -= sovAwardDef(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RspOfflineAwardData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.OfflineTime != 0 {
		n += 1 + sovAwardDef(uint64(m.OfflineTime))
	}
	if m.OfflineAwardDatas != nil {
		l = m.OfflineAwardDatas.Size()
		n += 1 + l + sovAwardDef(uint64(l))
	}
	return n
}

func (m *OfflineAwardData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ActorExp != 0 {
		n += 1 + sovAwardDef(uint64(m.ActorExp))
	}
	if m.SpecialAgentExp != 0 {
		n += 1 + sovAwardDef(uint64(m.SpecialAgentExp))
	}
	if m.Money != 0 {
		n += 1 + sovAwardDef(uint64(m.Money))
	}
	if len(m.OfflineAwardItems) > 0 {
		for _, e := range m.OfflineAwardItems {
			l = e.Size()
			n += 1 + l + sovAwardDef(uint64(l))
		}
	}
	return n
}

func (m *OfflineAwardItem) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.ID)
	if l > 0 {
		n += 1 + l + sovAwardDef(uint64(l))
	}
	if m.Type != 0 {
		n += 1 + sovAwardDef(uint64(m.Type))
	}
	if m.Num != 0 {
		n += 1 + sovAwardDef(uint64(m.Num))
	}
	if m.ConfigID != 0 {
		n += 1 + sovAwardDef(uint64(m.ConfigID))
	}
	return n
}

func sovAwardDef(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAwardDef(x uint64) (n int) {
	return sovAwardDef(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RspOfflineAwardData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAwardDef
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RspOfflineAwardData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RspOfflineAwardData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfflineTime", wireType)
			}
			m.OfflineTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OfflineTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfflineAwardDatas", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAwardDef
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAwardDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.OfflineAwardDatas == nil {
				m.OfflineAwardDatas = &OfflineAwardData{}
			}
			if err := m.OfflineAwardDatas.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAwardDef(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *OfflineAwardData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAwardDef
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: OfflineAwardData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OfflineAwardData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ActorExp", wireType)
			}
			m.ActorExp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ActorExp |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SpecialAgentExp", wireType)
			}
			m.SpecialAgentExp = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SpecialAgentExp |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Money", wireType)
			}
			m.Money = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Money |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OfflineAwardItems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAwardDef
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAwardDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OfflineAwardItems = append(m.OfflineAwardItems, &OfflineAwardItem{})
			if err := m.OfflineAwardItems[len(m.OfflineAwardItems)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAwardDef(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *OfflineAwardItem) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAwardDef
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: OfflineAwardItem: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: OfflineAwardItem: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAwardDef
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAwardDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= ItemEnum_Type(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Num", wireType)
			}
			m.Num = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Num |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConfigID", wireType)
			}
			m.ConfigID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ConfigID |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAwardDef(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAwardDef
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAwardDef(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAwardDef
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAwardDef
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthAwardDef
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAwardDef
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAwardDef
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAwardDef        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAwardDef          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAwardDef = fmt.Errorf("proto: unexpected end of group")
)
