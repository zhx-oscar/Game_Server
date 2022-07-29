// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: RoleCacheDef.proto

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

type RoleCache struct {
	Base            *RoleBase             `protobuf:"bytes,1,opt,name=Base,proto3" json:"Base,omitempty"`
	BuildMap        map[string]*BuildData `protobuf:"bytes,2,rep,name=BuildMap,proto3" json:"BuildMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	FightingBuildID string                `protobuf:"bytes,3,opt,name=FightingBuildID,proto3" json:"FightingBuildID,omitempty"`
}

func (m *RoleCache) Reset()         { *m = RoleCache{} }
func (m *RoleCache) String() string { return proto.CompactTextString(m) }
func (*RoleCache) ProtoMessage()    {}
func (*RoleCache) Descriptor() ([]byte, []int) {
	return fileDescriptor_d1bdab0e4e6dcbbd, []int{0}
}
func (m *RoleCache) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *RoleCache) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_RoleCache.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *RoleCache) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RoleCache.Merge(m, src)
}
func (m *RoleCache) XXX_Size() int {
	return m.Size()
}
func (m *RoleCache) XXX_DiscardUnknown() {
	xxx_messageInfo_RoleCache.DiscardUnknown(m)
}

var xxx_messageInfo_RoleCache proto.InternalMessageInfo

func (m *RoleCache) GetBase() *RoleBase {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *RoleCache) GetBuildMap() map[string]*BuildData {
	if m != nil {
		return m.BuildMap
	}
	return nil
}

func (m *RoleCache) GetFightingBuildID() string {
	if m != nil {
		return m.FightingBuildID
	}
	return ""
}

func init() {
	proto.RegisterType((*RoleCache)(nil), "RoleCache")
	proto.RegisterMapType((map[string]*BuildData)(nil), "RoleCache.BuildMapEntry")
}

func init() { proto.RegisterFile("RoleCacheDef.proto", fileDescriptor_d1bdab0e4e6dcbbd) }

var fileDescriptor_d1bdab0e4e6dcbbd = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x0a, 0xca, 0xcf, 0x49,
	0x75, 0x4e, 0x4c, 0xce, 0x48, 0x75, 0x49, 0x4d, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x97, 0xe2,
	0x05, 0x89, 0x21, 0xb8, 0xa2, 0xc1, 0x05, 0xa9, 0xc9, 0x99, 0x89, 0x39, 0x8e, 0xe9, 0xa9, 0x79,
	0x25, 0x70, 0x61, 0xa5, 0xcb, 0x8c, 0x5c, 0x9c, 0x70, 0xcd, 0x42, 0xb2, 0x5c, 0x2c, 0x4e, 0x89,
	0xc5, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x9c, 0x7a, 0x20, 0x19, 0x90, 0x40, 0x10,
	0x58, 0x58, 0xc8, 0x84, 0x8b, 0xc3, 0xa9, 0x34, 0x33, 0x27, 0xc5, 0x37, 0xb1, 0x40, 0x82, 0x49,
	0x81, 0x59, 0x83, 0xdb, 0x48, 0x42, 0x0f, 0xae, 0x59, 0x0f, 0x26, 0xe5, 0x9a, 0x57, 0x52, 0x54,
	0x19, 0x04, 0x57, 0x29, 0xa4, 0xc1, 0xc5, 0xef, 0x96, 0x99, 0x9e, 0x51, 0x92, 0x99, 0x97, 0x0e,
	0x16, 0xf3, 0x74, 0x91, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x0c, 0x42, 0x17, 0x96, 0x72, 0xe7, 0xe2,
	0x45, 0x31, 0x44, 0x48, 0x80, 0x8b, 0x39, 0x3b, 0xb5, 0x12, 0xec, 0x1c, 0xce, 0x20, 0x10, 0x53,
	0x48, 0x81, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x34, 0x55, 0x82, 0x09, 0xec, 0x44, 0x2e, 0x88, 0xad,
	0x2e, 0x89, 0x25, 0x89, 0x41, 0x10, 0x09, 0x2b, 0x26, 0x0b, 0x46, 0x27, 0xf9, 0x13, 0x8f, 0xe4,
	0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71, 0xc2, 0x63, 0x39, 0x86, 0x0b, 0x8f,
	0xe5, 0x18, 0x6e, 0x3c, 0x96, 0x63, 0x88, 0x62, 0x0d, 0x00, 0x79, 0x3b, 0x89, 0x0d, 0xec, 0x7b,
	0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xd0, 0xa4, 0xba, 0x09, 0x39, 0x01, 0x00, 0x00,
}

func (m *RoleCache) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *RoleCache) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *RoleCache) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.FightingBuildID) > 0 {
		i -= len(m.FightingBuildID)
		copy(dAtA[i:], m.FightingBuildID)
		i = encodeVarintRoleCacheDef(dAtA, i, uint64(len(m.FightingBuildID)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.BuildMap) > 0 {
		for k := range m.BuildMap {
			v := m.BuildMap[k]
			baseI := i
			if v != nil {
				{
					size, err := v.MarshalToSizedBuffer(dAtA[:i])
					if err != nil {
						return 0, err
					}
					i -= size
					i = encodeVarintRoleCacheDef(dAtA, i, uint64(size))
				}
				i--
				dAtA[i] = 0x12
			}
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintRoleCacheDef(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintRoleCacheDef(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Base != nil {
		{
			size, err := m.Base.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintRoleCacheDef(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRoleCacheDef(dAtA []byte, offset int, v uint64) int {
	offset -= sovRoleCacheDef(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *RoleCache) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Base != nil {
		l = m.Base.Size()
		n += 1 + l + sovRoleCacheDef(uint64(l))
	}
	if len(m.BuildMap) > 0 {
		for k, v := range m.BuildMap {
			_ = k
			_ = v
			l = 0
			if v != nil {
				l = v.Size()
				l += 1 + sovRoleCacheDef(uint64(l))
			}
			mapEntrySize := 1 + len(k) + sovRoleCacheDef(uint64(len(k))) + l
			n += mapEntrySize + 1 + sovRoleCacheDef(uint64(mapEntrySize))
		}
	}
	l = len(m.FightingBuildID)
	if l > 0 {
		n += 1 + l + sovRoleCacheDef(uint64(l))
	}
	return n
}

func sovRoleCacheDef(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRoleCacheDef(x uint64) (n int) {
	return sovRoleCacheDef(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *RoleCache) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRoleCacheDef
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
			return fmt.Errorf("proto: RoleCache: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RoleCache: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Base", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoleCacheDef
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
				return ErrInvalidLengthRoleCacheDef
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRoleCacheDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Base == nil {
				m.Base = &RoleBase{}
			}
			if err := m.Base.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BuildMap", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoleCacheDef
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
				return ErrInvalidLengthRoleCacheDef
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRoleCacheDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.BuildMap == nil {
				m.BuildMap = make(map[string]*BuildData)
			}
			var mapkey string
			var mapvalue *BuildData
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowRoleCacheDef
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
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRoleCacheDef
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthRoleCacheDef
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthRoleCacheDef
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var mapmsglen int
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowRoleCacheDef
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
						return ErrInvalidLengthRoleCacheDef
					}
					postmsgIndex := iNdEx + mapmsglen
					if postmsgIndex < 0 {
						return ErrInvalidLengthRoleCacheDef
					}
					if postmsgIndex > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = &BuildData{}
					if err := mapvalue.Unmarshal(dAtA[iNdEx:postmsgIndex]); err != nil {
						return err
					}
					iNdEx = postmsgIndex
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipRoleCacheDef(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if skippy < 0 {
						return ErrInvalidLengthRoleCacheDef
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.BuildMap[mapkey] = mapvalue
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FightingBuildID", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRoleCacheDef
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
				return ErrInvalidLengthRoleCacheDef
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRoleCacheDef
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.FightingBuildID = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRoleCacheDef(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthRoleCacheDef
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthRoleCacheDef
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
func skipRoleCacheDef(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRoleCacheDef
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
					return 0, ErrIntOverflowRoleCacheDef
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
					return 0, ErrIntOverflowRoleCacheDef
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
				return 0, ErrInvalidLengthRoleCacheDef
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRoleCacheDef
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRoleCacheDef
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRoleCacheDef        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRoleCacheDef          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRoleCacheDef = fmt.Errorf("proto: unexpected end of group")
)
