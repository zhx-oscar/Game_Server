// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: FightConst_Config.proto

package DataTables

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

type FightConst_Config struct {
	//* ID
	ID uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	//* 值
	Value string `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (m *FightConst_Config) Reset()         { *m = FightConst_Config{} }
func (m *FightConst_Config) String() string { return proto.CompactTextString(m) }
func (*FightConst_Config) ProtoMessage()    {}
func (*FightConst_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_81821a32efc38dbe, []int{0}
}
func (m *FightConst_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FightConst_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FightConst_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FightConst_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FightConst_Config.Merge(m, src)
}
func (m *FightConst_Config) XXX_Size() int {
	return m.Size()
}
func (m *FightConst_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_FightConst_Config.DiscardUnknown(m)
}

var xxx_messageInfo_FightConst_Config proto.InternalMessageInfo

func (m *FightConst_Config) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *FightConst_Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type FightConst_Config_Data struct {
	FightConst_ConfigItems map[uint32]*FightConst_Config `protobuf:"bytes,1,rep,name=FightConst_Config_items,json=FightConstConfigItems,proto3" json:"FightConst_Config_items,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *FightConst_Config_Data) Reset()         { *m = FightConst_Config_Data{} }
func (m *FightConst_Config_Data) String() string { return proto.CompactTextString(m) }
func (*FightConst_Config_Data) ProtoMessage()    {}
func (*FightConst_Config_Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_81821a32efc38dbe, []int{1}
}
func (m *FightConst_Config_Data) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FightConst_Config_Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_FightConst_Config_Data.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *FightConst_Config_Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FightConst_Config_Data.Merge(m, src)
}
func (m *FightConst_Config_Data) XXX_Size() int {
	return m.Size()
}
func (m *FightConst_Config_Data) XXX_DiscardUnknown() {
	xxx_messageInfo_FightConst_Config_Data.DiscardUnknown(m)
}

var xxx_messageInfo_FightConst_Config_Data proto.InternalMessageInfo

func (m *FightConst_Config_Data) GetFightConst_ConfigItems() map[uint32]*FightConst_Config {
	if m != nil {
		return m.FightConst_ConfigItems
	}
	return nil
}

func init() {
	proto.RegisterType((*FightConst_Config)(nil), "DataTables.FightConst_Config")
	proto.RegisterType((*FightConst_Config_Data)(nil), "DataTables.FightConst_Config_Data")
	proto.RegisterMapType((map[uint32]*FightConst_Config)(nil), "DataTables.FightConst_Config_Data.FightConstConfigItemsEntry")
}

func init() { proto.RegisterFile("FightConst_Config.proto", fileDescriptor_81821a32efc38dbe) }

var fileDescriptor_81821a32efc38dbe = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x77, 0xcb, 0x4c, 0xcf,
	0x28, 0x71, 0xce, 0xcf, 0x2b, 0x2e, 0x89, 0x77, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x72, 0x49, 0x2c, 0x49, 0x0c, 0x49, 0x4c, 0xca, 0x49, 0x2d, 0x56,
	0xb2, 0xe4, 0x12, 0xc4, 0x50, 0x26, 0xc4, 0xc7, 0xc5, 0xe4, 0xe9, 0x22, 0xc1, 0xa8, 0xc0, 0xa8,
	0xc1, 0x1b, 0xc4, 0xe4, 0xe9, 0x22, 0x24, 0xc2, 0xc5, 0x1a, 0x96, 0x98, 0x53, 0x9a, 0x2a, 0xc1,
	0xa4, 0xc0, 0xa8, 0xc1, 0x19, 0x04, 0xe1, 0x28, 0xfd, 0x60, 0xe4, 0x12, 0xc3, 0xd0, 0x1b, 0x0f,
	0x32, 0x5b, 0xa8, 0x14, 0x8b, 0xe5, 0xf1, 0x99, 0x25, 0xa9, 0xb9, 0xc5, 0x12, 0x8c, 0x0a, 0xcc,
	0x1a, 0xdc, 0x46, 0xb6, 0x7a, 0x08, 0x37, 0xe8, 0x61, 0x37, 0x04, 0x49, 0x18, 0x22, 0xea, 0x09,
	0xd2, 0xef, 0x9a, 0x57, 0x52, 0x54, 0x19, 0x24, 0x8a, 0x55, 0x4e, 0x2a, 0x9d, 0x4b, 0x0a, 0xb7,
	0x26, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0xa8, 0xb7, 0x40, 0x4c, 0x21, 0x63, 0x2e, 0xd6,
	0x32, 0xb8, 0xbf, 0xb8, 0x8d, 0x64, 0xf1, 0x3a, 0x2a, 0x08, 0xa2, 0xd6, 0x8a, 0xc9, 0x82, 0xd1,
	0x49, 0xff, 0xc4, 0x23, 0x39, 0xc6, 0x0b, 0x8f, 0xe4, 0x18, 0x1f, 0x3c, 0x92, 0x63, 0x9c, 0xf0,
	0x58, 0x8e, 0xe1, 0xc2, 0x63, 0x39, 0x86, 0x1b, 0x8f, 0xe5, 0x18, 0xb8, 0x44, 0x93, 0xf3, 0x73,
	0xf5, 0x5c, 0x12, 0x33, 0x8b, 0x2b, 0xf5, 0x8a, 0x53, 0x8b, 0xca, 0x52, 0x8b, 0xf4, 0x52, 0x12,
	0x4b, 0x12, 0x93, 0xd8, 0xc0, 0x21, 0x6f, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x28, 0x3b, 0xdb,
	0xd6, 0x94, 0x01, 0x00, 0x00,
}

func (m *FightConst_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FightConst_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FightConst_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Value) > 0 {
		i -= len(m.Value)
		copy(dAtA[i:], m.Value)
		i = encodeVarintFightConst_Config(dAtA, i, uint64(len(m.Value)))
		i--
		dAtA[i] = 0x12
	}
	if m.ID != 0 {
		i = encodeVarintFightConst_Config(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *FightConst_Config_Data) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FightConst_Config_Data) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FightConst_Config_Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FightConst_ConfigItems) > 0 {
		for k := range m.FightConst_ConfigItems {
			v := m.FightConst_ConfigItems[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintFightConst_Config(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i = encodeVarintFightConst_Config(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintFightConst_Config(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintFightConst_Config(dAtA []byte, offset int, v uint64) int {
	offset -= sovFightConst_Config(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *FightConst_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovFightConst_Config(uint64(m.ID))
	}
	l = len(m.Value)
	if l > 0 {
		n += 1 + l + sovFightConst_Config(uint64(l))
	}
	return n
}

func (m *FightConst_Config_Data) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.FightConst_ConfigItems) > 0 {
		for k, v := range m.FightConst_ConfigItems {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovFightConst_Config(uint64(l))
			}
			mapEntrySize := 1 + sovFightConst_Config(uint64(k)) + l
			n += mapEntrySize + 1 + sovFightConst_Config(uint64(mapEntrySize))
		}
	}
	return n
}

func sovFightConst_Config(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFightConst_Config(x uint64) (n int) {
	return sovFightConst_Config(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *FightConst_Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFightConst_Config
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
			return fmt.Errorf("proto: FightConst_Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FightConst_Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFightConst_Config
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ID |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFightConst_Config
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
				return ErrInvalidLengthFightConst_Config
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthFightConst_Config
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Value = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFightConst_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFightConst_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFightConst_Config
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
func (m *FightConst_Config_Data) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFightConst_Config
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
			return fmt.Errorf("proto: FightConst_Config_Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FightConst_Config_Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FightConst_ConfigItems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFightConst_Config
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
				return ErrInvalidLengthFightConst_Config
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFightConst_Config
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.FightConst_ConfigItems == nil {
				m.FightConst_ConfigItems = make(map[uint32]*FightConst_Config)
			}
			var mapkey uint32
			var mapvalue *FightConst_Config
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowFightConst_Config
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
				if fieldNum == 1 {
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFightConst_Config
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapkey |= uint32(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowFightConst_Config
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						mapmsglen |= int(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					if mapmsglen < 0 {
						return ErrInvalidLengthFightConst_Config
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthFightConst_Config
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &FightConst_Config{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipFightConst_Config(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthFightConst_Config
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.FightConst_ConfigItems[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFightConst_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFightConst_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFightConst_Config
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
func skipFightConst_Config(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFightConst_Config
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
					return 0, ErrIntOverflowFightConst_Config
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
					return 0, ErrIntOverflowFightConst_Config
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
				return 0, ErrInvalidLengthFightConst_Config
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupFightConst_Config
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthFightConst_Config
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthFightConst_Config        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFightConst_Config          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupFightConst_Config = fmt.Errorf("proto: unexpected end of group")
)
