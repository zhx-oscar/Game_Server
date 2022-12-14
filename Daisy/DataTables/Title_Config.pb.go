// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: Title_Config.proto

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

type Title_Config struct {
	//* 称号ID
	ID uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	//* 称号名称
	Name string `protobuf:"bytes,2,opt,name=Name,proto3" json:"Name,omitempty"`
	//* 称号等级
	Level uint32 `protobuf:"varint,3,opt,name=Level,proto3" json:"Level,omitempty"`
	//* 称号来源
	Type uint32 `protobuf:"varint,4,opt,name=type,proto3" json:"type,omitempty"`
	//* 称号有效期
	TitleDuration uint32 `protobuf:"varint,5,opt,name=titleDuration,proto3" json:"titleDuration,omitempty"`
}

func (m *Title_Config) Reset()         { *m = Title_Config{} }
func (m *Title_Config) String() string { return proto.CompactTextString(m) }
func (*Title_Config) ProtoMessage()    {}
func (*Title_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_686f420a870b1345, []int{0}
}
func (m *Title_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Title_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Title_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Title_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Title_Config.Merge(m, src)
}
func (m *Title_Config) XXX_Size() int {
	return m.Size()
}
func (m *Title_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Title_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Title_Config proto.InternalMessageInfo

func (m *Title_Config) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *Title_Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Title_Config) GetLevel() uint32 {
	if m != nil {
		return m.Level
	}
	return 0
}

func (m *Title_Config) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Title_Config) GetTitleDuration() uint32 {
	if m != nil {
		return m.TitleDuration
	}
	return 0
}

type Title_Config_Data struct {
	Title_ConfigItems map[uint32]*Title_Config `protobuf:"bytes,1,rep,name=Title_Config_items,json=TitleConfigItems,proto3" json:"Title_Config_items,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *Title_Config_Data) Reset()         { *m = Title_Config_Data{} }
func (m *Title_Config_Data) String() string { return proto.CompactTextString(m) }
func (*Title_Config_Data) ProtoMessage()    {}
func (*Title_Config_Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_686f420a870b1345, []int{1}
}
func (m *Title_Config_Data) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Title_Config_Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Title_Config_Data.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Title_Config_Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Title_Config_Data.Merge(m, src)
}
func (m *Title_Config_Data) XXX_Size() int {
	return m.Size()
}
func (m *Title_Config_Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Title_Config_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Title_Config_Data proto.InternalMessageInfo

func (m *Title_Config_Data) GetTitle_ConfigItems() map[uint32]*Title_Config {
	if m != nil {
		return m.Title_ConfigItems
	}
	return nil
}

func init() {
	proto.RegisterType((*Title_Config)(nil), "DataTables.Title_Config")
	proto.RegisterType((*Title_Config_Data)(nil), "DataTables.Title_Config_Data")
	proto.RegisterMapType((map[uint32]*Title_Config)(nil), "DataTables.Title_Config_Data.TitleConfigItemsEntry")
}

func init() { proto.RegisterFile("Title_Config.proto", fileDescriptor_686f420a870b1345) }

var fileDescriptor_686f420a870b1345 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x0a, 0xc9, 0x2c, 0xc9,
	0x49, 0x8d, 0x77, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2,
	0x72, 0x49, 0x2c, 0x49, 0x0c, 0x49, 0x4c, 0xca, 0x49, 0x2d, 0x56, 0x6a, 0x62, 0xe4, 0xe2, 0x41,
	0x56, 0x22, 0xc4, 0xc7, 0xc5, 0xe4, 0xe9, 0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1b, 0xc4, 0xe4,
	0xe9, 0x22, 0x24, 0xc4, 0xc5, 0xe2, 0x97, 0x98, 0x9b, 0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x19,
	0x04, 0x66, 0x0b, 0x89, 0x70, 0xb1, 0xfa, 0xa4, 0x96, 0xa5, 0xe6, 0x48, 0x30, 0x83, 0x95, 0x41,
	0x38, 0x20, 0x95, 0x25, 0x95, 0x05, 0xa9, 0x12, 0x2c, 0x60, 0x41, 0x30, 0x5b, 0x48, 0x85, 0x8b,
	0xb7, 0x04, 0x64, 0xba, 0x4b, 0x69, 0x51, 0x62, 0x49, 0x66, 0x7e, 0x9e, 0x04, 0x2b, 0x58, 0x12,
	0x55, 0x50, 0xe9, 0x2a, 0x23, 0x97, 0x20, 0xb2, 0x23, 0xe2, 0x41, 0x0e, 0x14, 0x4a, 0x44, 0x75,
	0x7c, 0x7c, 0x66, 0x49, 0x6a, 0x6e, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0xb1, 0x1e,
	0xc2, 0x0f, 0x7a, 0x18, 0x5a, 0x21, 0x22, 0x10, 0x01, 0x4f, 0x90, 0x2e, 0xd7, 0xbc, 0x92, 0xa2,
	0xca, 0x20, 0x01, 0x74, 0x61, 0xa9, 0x58, 0x2e, 0x51, 0xac, 0x4a, 0x85, 0x04, 0xb8, 0x98, 0xb3,
	0x53, 0x2b, 0xa1, 0xc1, 0x00, 0x62, 0x0a, 0xe9, 0x71, 0xb1, 0x96, 0x25, 0xe6, 0x94, 0x42, 0x02,
	0x82, 0xdb, 0x48, 0x02, 0x97, 0x03, 0x82, 0x20, 0xca, 0xac, 0x98, 0x2c, 0x18, 0x9d, 0xf4, 0x4f,
	0x3c, 0x92, 0x63, 0xbc, 0xf0, 0x48, 0x8e, 0xf1, 0xc1, 0x23, 0x39, 0xc6, 0x09, 0x8f, 0xe5, 0x18,
	0x2e, 0x3c, 0x96, 0x63, 0xb8, 0xf1, 0x58, 0x8e, 0x81, 0x4b, 0x34, 0x39, 0x3f, 0x57, 0xcf, 0x25,
	0x31, 0xb3, 0xb8, 0x52, 0xaf, 0x38, 0xb5, 0xa8, 0x2c, 0xb5, 0x48, 0x2f, 0x25, 0xb1, 0x24, 0x31,
	0x89, 0x0d, 0x1c, 0x41, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xee, 0xae, 0x20, 0x12, 0xb6,
	0x01, 0x00, 0x00,
}

func (m *Title_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Title_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Title_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.TitleDuration != 0 {
		i = encodeVarintTitle_Config(dAtA, i, uint64(m.TitleDuration))
		i--
		dAtA[i] = 0x28
	}
	if m.Type != 0 {
		i = encodeVarintTitle_Config(dAtA, i, uint64(m.Type))
		i--
		dAtA[i] = 0x20
	}
	if m.Level != 0 {
		i = encodeVarintTitle_Config(dAtA, i, uint64(m.Level))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintTitle_Config(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if m.ID != 0 {
		i = encodeVarintTitle_Config(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Title_Config_Data) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Title_Config_Data) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Title_Config_Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Title_ConfigItems) > 0 {
		for k := range m.Title_ConfigItems {
			v := m.Title_ConfigItems[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintTitle_Config(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i = encodeVarintTitle_Config(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintTitle_Config(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTitle_Config(dAtA []byte, offset int, v uint64) int {
	offset -= sovTitle_Config(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Title_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovTitle_Config(uint64(m.ID))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovTitle_Config(uint64(l))
	}
	if m.Level != 0 {
		n += 1 + sovTitle_Config(uint64(m.Level))
	}
	if m.Type != 0 {
		n += 1 + sovTitle_Config(uint64(m.Type))
	}
	if m.TitleDuration != 0 {
		n += 1 + sovTitle_Config(uint64(m.TitleDuration))
	}
	return n
}

func (m *Title_Config_Data) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Title_ConfigItems) > 0 {
		for k, v := range m.Title_ConfigItems {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovTitle_Config(uint64(l))
			}
			mapEntrySize := 1 + sovTitle_Config(uint64(k)) + l
			n += mapEntrySize + 1 + sovTitle_Config(uint64(mapEntrySize))
		}
	}
	return n
}

func sovTitle_Config(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTitle_Config(x uint64) (n int) {
	return sovTitle_Config(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Title_Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTitle_Config
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
			return fmt.Errorf("proto: Title_Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Title_Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
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
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
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
				return ErrInvalidLengthTitle_Config
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTitle_Config
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Level", wireType)
			}
			m.Level = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Level |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			m.Type = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Type |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TitleDuration", wireType)
			}
			m.TitleDuration = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TitleDuration |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTitle_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTitle_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTitle_Config
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
func (m *Title_Config_Data) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTitle_Config
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
			return fmt.Errorf("proto: Title_Config_Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Title_Config_Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title_ConfigItems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTitle_Config
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
				return ErrInvalidLengthTitle_Config
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTitle_Config
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Title_ConfigItems == nil {
				m.Title_ConfigItems = make(map[uint32]*Title_Config)
			}
			var mapkey uint32
			var mapvalue *Title_Config
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowTitle_Config
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
							return ErrIntOverflowTitle_Config
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
							return ErrIntOverflowTitle_Config
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
						return ErrInvalidLengthTitle_Config
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthTitle_Config
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &Title_Config{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipTitle_Config(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthTitle_Config
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Title_ConfigItems[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTitle_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTitle_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTitle_Config
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
func skipTitle_Config(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTitle_Config
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
					return 0, ErrIntOverflowTitle_Config
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
					return 0, ErrIntOverflowTitle_Config
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
				return 0, ErrInvalidLengthTitle_Config
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTitle_Config
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTitle_Config
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTitle_Config        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTitle_Config          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTitle_Config = fmt.Errorf("proto: unexpected end of group")
)
