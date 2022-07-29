// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ItemType_Config.proto

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

type ItemPackMaxSpace_Config struct {
	//* ID
	ID uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	//* 最大空间
	MaxSpace uint32 `protobuf:"varint,2,opt,name=MaxSpace,proto3" json:"MaxSpace,omitempty"`
}

func (m *ItemPackMaxSpace_Config) Reset()         { *m = ItemPackMaxSpace_Config{} }
func (m *ItemPackMaxSpace_Config) String() string { return proto.CompactTextString(m) }
func (*ItemPackMaxSpace_Config) ProtoMessage()    {}
func (*ItemPackMaxSpace_Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a3b29f8ec0f8099, []int{0}
}
func (m *ItemPackMaxSpace_Config) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ItemPackMaxSpace_Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ItemPackMaxSpace_Config.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ItemPackMaxSpace_Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemPackMaxSpace_Config.Merge(m, src)
}
func (m *ItemPackMaxSpace_Config) XXX_Size() int {
	return m.Size()
}
func (m *ItemPackMaxSpace_Config) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemPackMaxSpace_Config.DiscardUnknown(m)
}

var xxx_messageInfo_ItemPackMaxSpace_Config proto.InternalMessageInfo

func (m *ItemPackMaxSpace_Config) GetID() uint32 {
	if m != nil {
		return m.ID
	}
	return 0
}

func (m *ItemPackMaxSpace_Config) GetMaxSpace() uint32 {
	if m != nil {
		return m.MaxSpace
	}
	return 0
}

type ItemType_Config_Data struct {
	ItemPackMaxSpace_ConfigItems map[uint32]*ItemPackMaxSpace_Config `protobuf:"bytes,1,rep,name=ItemPackMaxSpace_Config_items,json=ItemPackMaxSpaceConfigItems,proto3" json:"ItemPackMaxSpace_Config_items,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (m *ItemType_Config_Data) Reset()         { *m = ItemType_Config_Data{} }
func (m *ItemType_Config_Data) String() string { return proto.CompactTextString(m) }
func (*ItemType_Config_Data) ProtoMessage()    {}
func (*ItemType_Config_Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_4a3b29f8ec0f8099, []int{1}
}
func (m *ItemType_Config_Data) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ItemType_Config_Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ItemType_Config_Data.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ItemType_Config_Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ItemType_Config_Data.Merge(m, src)
}
func (m *ItemType_Config_Data) XXX_Size() int {
	return m.Size()
}
func (m *ItemType_Config_Data) XXX_DiscardUnknown() {
	xxx_messageInfo_ItemType_Config_Data.DiscardUnknown(m)
}

var xxx_messageInfo_ItemType_Config_Data proto.InternalMessageInfo

func (m *ItemType_Config_Data) GetItemPackMaxSpace_ConfigItems() map[uint32]*ItemPackMaxSpace_Config {
	if m != nil {
		return m.ItemPackMaxSpace_ConfigItems
	}
	return nil
}

func init() {
	proto.RegisterType((*ItemPackMaxSpace_Config)(nil), "DataTables.ItemPackMaxSpace_Config")
	proto.RegisterType((*ItemType_Config_Data)(nil), "DataTables.ItemType_Config_Data")
	proto.RegisterMapType((map[uint32]*ItemPackMaxSpace_Config)(nil), "DataTables.ItemType_Config_Data.ItemPackMaxSpaceConfigItemsEntry")
}

func init() { proto.RegisterFile("ItemType_Config.proto", fileDescriptor_4a3b29f8ec0f8099) }

var fileDescriptor_4a3b29f8ec0f8099 = []byte{
	// 260 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0xf5, 0x2c, 0x49, 0xcd,
	0x0d, 0xa9, 0x2c, 0x48, 0x8d, 0x77, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0xd7, 0x2b, 0x28, 0xca, 0x2f,
	0xc9, 0x17, 0xe2, 0x72, 0x49, 0x2c, 0x49, 0x0c, 0x49, 0x4c, 0xca, 0x49, 0x2d, 0x56, 0x72, 0xe5,
	0x12, 0x07, 0x29, 0x0a, 0x48, 0x4c, 0xce, 0xf6, 0x4d, 0xac, 0x08, 0x2e, 0x48, 0x4c, 0x86, 0x29,
	0x16, 0xe2, 0xe3, 0x62, 0xf2, 0x74, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0d, 0x62, 0xf2, 0x74,
	0x11, 0x92, 0xe2, 0xe2, 0x80, 0x29, 0x91, 0x60, 0x02, 0x8b, 0xc2, 0xf9, 0x4a, 0x93, 0x99, 0xb8,
	0x44, 0xd0, 0x2c, 0x8b, 0x07, 0xd9, 0x22, 0xd4, 0xca, 0xc8, 0x25, 0x8b, 0xc3, 0x82, 0xf8, 0xcc,
	0x92, 0xd4, 0xdc, 0x62, 0x09, 0x46, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x47, 0x3d, 0x84, 0xa3, 0xf4,
	0xb0, 0x99, 0xa4, 0x87, 0x6e, 0x0a, 0x44, 0x0e, 0x24, 0x5a, 0xec, 0x9a, 0x57, 0x52, 0x54, 0x19,
	0x24, 0x8d, 0x47, 0x85, 0x54, 0x31, 0x97, 0x02, 0x21, 0x03, 0x84, 0x04, 0xb8, 0x98, 0xb3, 0x53,
	0x2b, 0xa1, 0x3e, 0x06, 0x31, 0x85, 0x2c, 0xb9, 0x58, 0xcb, 0x12, 0x73, 0x4a, 0x21, 0xfe, 0xe5,
	0x36, 0x52, 0x46, 0x77, 0x24, 0x16, 0x5f, 0x05, 0x41, 0x74, 0x58, 0x31, 0x59, 0x30, 0x3a, 0xe9,
	0x9f, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c, 0x13, 0x1e, 0xcb,
	0x31, 0x5c, 0x78, 0x2c, 0xc7, 0x70, 0xe3, 0xb1, 0x1c, 0x03, 0x97, 0x68, 0x72, 0x7e, 0xae, 0x9e,
	0x4b, 0x62, 0x66, 0x71, 0xa5, 0x5e, 0x71, 0x6a, 0x51, 0x59, 0x6a, 0x91, 0x5e, 0x4a, 0x62, 0x49,
	0x62, 0x12, 0x1b, 0x38, 0x82, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xf6, 0x2a, 0xd3,
	0xb9, 0x01, 0x00, 0x00,
}

func (m *ItemPackMaxSpace_Config) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ItemPackMaxSpace_Config) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ItemPackMaxSpace_Config) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MaxSpace != 0 {
		i = encodeVarintItemType_Config(dAtA, i, uint64(m.MaxSpace))
		i--
		dAtA[i] = 0x10
	}
	if m.ID != 0 {
		i = encodeVarintItemType_Config(dAtA, i, uint64(m.ID))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *ItemType_Config_Data) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ItemType_Config_Data) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ItemType_Config_Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ItemPackMaxSpace_ConfigItems) > 0 {
		for k := range m.ItemPackMaxSpace_ConfigItems {
			v := m.ItemPackMaxSpace_ConfigItems[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintItemType_Config(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i = encodeVarintItemType_Config(dAtA, i, uint64(k))
			i--
			dAtA[i] = 0x8
			i = encodeVarintItemType_Config(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintItemType_Config(dAtA []byte, offset int, v uint64) int {
	offset -= sovItemType_Config(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ItemPackMaxSpace_Config) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ID != 0 {
		n += 1 + sovItemType_Config(uint64(m.ID))
	}
	if m.MaxSpace != 0 {
		n += 1 + sovItemType_Config(uint64(m.MaxSpace))
	}
	return n
}

func (m *ItemType_Config_Data) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.ItemPackMaxSpace_ConfigItems) > 0 {
		for k, v := range m.ItemPackMaxSpace_ConfigItems {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovItemType_Config(uint64(l))
			}
			mapEntrySize := 1 + sovItemType_Config(uint64(k)) + l
			n += mapEntrySize + 1 + sovItemType_Config(uint64(mapEntrySize))
		}
	}
	return n
}

func sovItemType_Config(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozItemType_Config(x uint64) (n int) {
	return sovItemType_Config(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ItemPackMaxSpace_Config) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowItemType_Config
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
			return fmt.Errorf("proto: ItemPackMaxSpace_Config: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ItemPackMaxSpace_Config: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ID", wireType)
			}
			m.ID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowItemType_Config
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
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxSpace", wireType)
			}
			m.MaxSpace = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowItemType_Config
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxSpace |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipItemType_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthItemType_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthItemType_Config
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
func (m *ItemType_Config_Data) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowItemType_Config
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
			return fmt.Errorf("proto: ItemType_Config_Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ItemType_Config_Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ItemPackMaxSpace_ConfigItems", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowItemType_Config
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
				return ErrInvalidLengthItemType_Config
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthItemType_Config
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ItemPackMaxSpace_ConfigItems == nil {
				m.ItemPackMaxSpace_ConfigItems = make(map[uint32]*ItemPackMaxSpace_Config)
			}
			var mapkey uint32
			var mapvalue *ItemPackMaxSpace_Config
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowItemType_Config
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
							return ErrIntOverflowItemType_Config
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
							return ErrIntOverflowItemType_Config
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
						return ErrInvalidLengthItemType_Config
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthItemType_Config
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &ItemPackMaxSpace_Config{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipItemType_Config(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthItemType_Config
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.ItemPackMaxSpace_ConfigItems[mapkey] = mapvalue
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipItemType_Config(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthItemType_Config
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthItemType_Config
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
func skipItemType_Config(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowItemType_Config
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
					return 0, ErrIntOverflowItemType_Config
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
					return 0, ErrIntOverflowItemType_Config
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
				return 0, ErrInvalidLengthItemType_Config
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupItemType_Config
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthItemType_Config
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthItemType_Config        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowItemType_Config          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupItemType_Config = fmt.Errorf("proto: unexpected end of group")
)
