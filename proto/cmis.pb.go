// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/cmis.proto

package cmisproto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Repository struct {
	Id                   int32             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description          string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	RootFolder           *CmisObject       `protobuf:"bytes,4,opt,name=root_folder,json=rootFolder,proto3" json:"root_folder,omitempty"`
	TypeDefinitions      []*TypeDefinition `protobuf:"bytes,5,rep,name=type_definitions,json=typeDefinitions,proto3" json:"type_definitions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Repository) Reset()         { *m = Repository{} }
func (m *Repository) String() string { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()    {}
func (*Repository) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{0}
}

func (m *Repository) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Repository.Unmarshal(m, b)
}
func (m *Repository) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Repository.Marshal(b, m, deterministic)
}
func (m *Repository) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Repository.Merge(m, src)
}
func (m *Repository) XXX_Size() int {
	return xxx_messageInfo_Repository.Size(m)
}
func (m *Repository) XXX_DiscardUnknown() {
	xxx_messageInfo_Repository.DiscardUnknown(m)
}

var xxx_messageInfo_Repository proto.InternalMessageInfo

func (m *Repository) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Repository) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Repository) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Repository) GetRootFolder() *CmisObject {
	if m != nil {
		return m.RootFolder
	}
	return nil
}

func (m *Repository) GetTypeDefinitions() []*TypeDefinition {
	if m != nil {
		return m.TypeDefinitions
	}
	return nil
}

type TypeDefinition struct {
	Name                 string                `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description          string                `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	PropertyDefinitions  []*PropertyDefinition `protobuf:"bytes,4,rep,name=property_definitions,json=propertyDefinitions,proto3" json:"property_definitions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *TypeDefinition) Reset()         { *m = TypeDefinition{} }
func (m *TypeDefinition) String() string { return proto.CompactTextString(m) }
func (*TypeDefinition) ProtoMessage()    {}
func (*TypeDefinition) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{1}
}

func (m *TypeDefinition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TypeDefinition.Unmarshal(m, b)
}
func (m *TypeDefinition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TypeDefinition.Marshal(b, m, deterministic)
}
func (m *TypeDefinition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TypeDefinition.Merge(m, src)
}
func (m *TypeDefinition) XXX_Size() int {
	return xxx_messageInfo_TypeDefinition.Size(m)
}
func (m *TypeDefinition) XXX_DiscardUnknown() {
	xxx_messageInfo_TypeDefinition.DiscardUnknown(m)
}

var xxx_messageInfo_TypeDefinition proto.InternalMessageInfo

func (m *TypeDefinition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *TypeDefinition) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *TypeDefinition) GetPropertyDefinitions() []*PropertyDefinition {
	if m != nil {
		return m.PropertyDefinitions
	}
	return nil
}

type PropertyDefinition struct {
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description          string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Datatype             string   `protobuf:"bytes,4,opt,name=datatype,proto3" json:"datatype,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PropertyDefinition) Reset()         { *m = PropertyDefinition{} }
func (m *PropertyDefinition) String() string { return proto.CompactTextString(m) }
func (*PropertyDefinition) ProtoMessage()    {}
func (*PropertyDefinition) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{2}
}

func (m *PropertyDefinition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PropertyDefinition.Unmarshal(m, b)
}
func (m *PropertyDefinition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PropertyDefinition.Marshal(b, m, deterministic)
}
func (m *PropertyDefinition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PropertyDefinition.Merge(m, src)
}
func (m *PropertyDefinition) XXX_Size() int {
	return xxx_messageInfo_PropertyDefinition.Size(m)
}
func (m *PropertyDefinition) XXX_DiscardUnknown() {
	xxx_messageInfo_PropertyDefinition.DiscardUnknown(m)
}

var xxx_messageInfo_PropertyDefinition proto.InternalMessageInfo

func (m *PropertyDefinition) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PropertyDefinition) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *PropertyDefinition) GetDatatype() string {
	if m != nil {
		return m.Datatype
	}
	return ""
}

type CmisObject struct {
	Id                   *CmisObjectId   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	TypeDefinition       *TypeDefinition `protobuf:"bytes,2,opt,name=type_definition,json=typeDefinition,proto3" json:"type_definition,omitempty"`
	Properties           []*CmisProperty `protobuf:"bytes,3,rep,name=properties,proto3" json:"properties,omitempty"`
	Children             []*CmisObject   `protobuf:"bytes,4,rep,name=children,proto3" json:"children,omitempty"`
	Parents              []*CmisObject   `protobuf:"bytes,5,rep,name=parents,proto3" json:"parents,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CmisObject) Reset()         { *m = CmisObject{} }
func (m *CmisObject) String() string { return proto.CompactTextString(m) }
func (*CmisObject) ProtoMessage()    {}
func (*CmisObject) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{3}
}

func (m *CmisObject) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmisObject.Unmarshal(m, b)
}
func (m *CmisObject) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmisObject.Marshal(b, m, deterministic)
}
func (m *CmisObject) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmisObject.Merge(m, src)
}
func (m *CmisObject) XXX_Size() int {
	return xxx_messageInfo_CmisObject.Size(m)
}
func (m *CmisObject) XXX_DiscardUnknown() {
	xxx_messageInfo_CmisObject.DiscardUnknown(m)
}

var xxx_messageInfo_CmisObject proto.InternalMessageInfo

func (m *CmisObject) GetId() *CmisObjectId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *CmisObject) GetTypeDefinition() *TypeDefinition {
	if m != nil {
		return m.TypeDefinition
	}
	return nil
}

func (m *CmisObject) GetProperties() []*CmisProperty {
	if m != nil {
		return m.Properties
	}
	return nil
}

func (m *CmisObject) GetChildren() []*CmisObject {
	if m != nil {
		return m.Children
	}
	return nil
}

func (m *CmisObject) GetParents() []*CmisObject {
	if m != nil {
		return m.Parents
	}
	return nil
}

type CmisProperty struct {
	PropertyDefinition   *PropertyDefinition `protobuf:"bytes,1,opt,name=property_definition,json=propertyDefinition,proto3" json:"property_definition,omitempty"`
	Value                string              `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CmisProperty) Reset()         { *m = CmisProperty{} }
func (m *CmisProperty) String() string { return proto.CompactTextString(m) }
func (*CmisProperty) ProtoMessage()    {}
func (*CmisProperty) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{4}
}

func (m *CmisProperty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmisProperty.Unmarshal(m, b)
}
func (m *CmisProperty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmisProperty.Marshal(b, m, deterministic)
}
func (m *CmisProperty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmisProperty.Merge(m, src)
}
func (m *CmisProperty) XXX_Size() int {
	return xxx_messageInfo_CmisProperty.Size(m)
}
func (m *CmisProperty) XXX_DiscardUnknown() {
	xxx_messageInfo_CmisProperty.DiscardUnknown(m)
}

var xxx_messageInfo_CmisProperty proto.InternalMessageInfo

func (m *CmisProperty) GetPropertyDefinition() *PropertyDefinition {
	if m != nil {
		return m.PropertyDefinition
	}
	return nil
}

func (m *CmisProperty) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type CmisObjectId struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CmisObjectId) Reset()         { *m = CmisObjectId{} }
func (m *CmisObjectId) String() string { return proto.CompactTextString(m) }
func (*CmisObjectId) ProtoMessage()    {}
func (*CmisObjectId) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{5}
}

func (m *CmisObjectId) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CmisObjectId.Unmarshal(m, b)
}
func (m *CmisObjectId) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CmisObjectId.Marshal(b, m, deterministic)
}
func (m *CmisObjectId) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CmisObjectId.Merge(m, src)
}
func (m *CmisObjectId) XXX_Size() int {
	return xxx_messageInfo_CmisObjectId.Size(m)
}
func (m *CmisObjectId) XXX_DiscardUnknown() {
	xxx_messageInfo_CmisObjectId.DiscardUnknown(m)
}

var xxx_messageInfo_CmisObjectId proto.InternalMessageInfo

func (m *CmisObjectId) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

type CreateObjectReq struct {
	Name                 string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type                 string        `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	ParentId             *CmisObjectId `protobuf:"bytes,3,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	RepositoryId         int32         `protobuf:"varint,4,opt,name=repository_id,json=repositoryId,proto3" json:"repository_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *CreateObjectReq) Reset()         { *m = CreateObjectReq{} }
func (m *CreateObjectReq) String() string { return proto.CompactTextString(m) }
func (*CreateObjectReq) ProtoMessage()    {}
func (*CreateObjectReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_5ed4f7310ba83cea, []int{6}
}

func (m *CreateObjectReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateObjectReq.Unmarshal(m, b)
}
func (m *CreateObjectReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateObjectReq.Marshal(b, m, deterministic)
}
func (m *CreateObjectReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateObjectReq.Merge(m, src)
}
func (m *CreateObjectReq) XXX_Size() int {
	return xxx_messageInfo_CreateObjectReq.Size(m)
}
func (m *CreateObjectReq) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateObjectReq.DiscardUnknown(m)
}

var xxx_messageInfo_CreateObjectReq proto.InternalMessageInfo

func (m *CreateObjectReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateObjectReq) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *CreateObjectReq) GetParentId() *CmisObjectId {
	if m != nil {
		return m.ParentId
	}
	return nil
}

func (m *CreateObjectReq) GetRepositoryId() int32 {
	if m != nil {
		return m.RepositoryId
	}
	return 0
}

func init() {
	proto.RegisterType((*Repository)(nil), "cmis.Repository")
	proto.RegisterType((*TypeDefinition)(nil), "cmis.TypeDefinition")
	proto.RegisterType((*PropertyDefinition)(nil), "cmis.PropertyDefinition")
	proto.RegisterType((*CmisObject)(nil), "cmis.CmisObject")
	proto.RegisterType((*CmisProperty)(nil), "cmis.CmisProperty")
	proto.RegisterType((*CmisObjectId)(nil), "cmis.CmisObjectId")
	proto.RegisterType((*CreateObjectReq)(nil), "cmis.CreateObjectReq")
}

func init() { proto.RegisterFile("proto/cmis.proto", fileDescriptor_5ed4f7310ba83cea) }

var fileDescriptor_5ed4f7310ba83cea = []byte{
	// 551 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xcd, 0x8e, 0x12, 0x41,
	0x10, 0x4e, 0xb3, 0x83, 0x42, 0x0d, 0x0b, 0xa4, 0x16, 0xcd, 0x04, 0x13, 0x43, 0xc6, 0x0b, 0x31,
	0x06, 0x5c, 0x34, 0x31, 0xc6, 0x18, 0x13, 0x57, 0x4d, 0x88, 0x07, 0xcd, 0xe8, 0xc9, 0x0b, 0x19,
	0xa6, 0x0b, 0x6c, 0x03, 0xf4, 0xd8, 0xd3, 0x6c, 0xc2, 0x53, 0x78, 0xf0, 0x21, 0x7c, 0x13, 0xdf,
	0xca, 0xc4, 0x74, 0xcf, 0xcf, 0x0e, 0x0c, 0x7b, 0xd8, 0x3d, 0x51, 0x5d, 0xf3, 0x15, 0xf5, 0xd5,
	0xf7, 0x55, 0x37, 0x74, 0x63, 0x25, 0xb5, 0x1c, 0x47, 0x6b, 0x91, 0x8c, 0x6c, 0x88, 0x8e, 0x89,
	0xfb, 0x0f, 0x96, 0x52, 0x2e, 0x57, 0x34, 0xb6, 0xb9, 0xf9, 0x76, 0x31, 0xa6, 0x75, 0xac, 0x77,
	0x29, 0xc4, 0xff, 0xcb, 0x00, 0x02, 0x8a, 0x65, 0x22, 0xb4, 0x54, 0x3b, 0x6c, 0x43, 0x4d, 0x70,
	0x8f, 0x0d, 0xd8, 0xb0, 0x1e, 0xd4, 0x04, 0x47, 0x04, 0x67, 0x13, 0xae, 0xc9, 0xab, 0x0d, 0xd8,
	0xb0, 0x19, 0xd8, 0x18, 0x07, 0xe0, 0x72, 0x4a, 0x22, 0x25, 0x62, 0x2d, 0xe4, 0xc6, 0x3b, 0xb1,
	0x9f, 0xca, 0x29, 0x3c, 0x07, 0x57, 0x49, 0xa9, 0x67, 0x0b, 0xb9, 0xe2, 0xa4, 0x3c, 0x67, 0xc0,
	0x86, 0xee, 0xa4, 0x3b, 0xb2, 0xcc, 0x2e, 0xd6, 0x22, 0xf9, 0x34, 0xff, 0x41, 0x91, 0x0e, 0xc0,
	0x80, 0x3e, 0x58, 0x0c, 0xbe, 0x81, 0xae, 0xde, 0xc5, 0x34, 0xe3, 0xb4, 0x10, 0x1b, 0x61, 0xfe,
	0x25, 0xf1, 0xea, 0x83, 0x93, 0xa1, 0x3b, 0xe9, 0xa5, 0x75, 0x5f, 0x77, 0x31, 0xbd, 0x2b, 0x3e,
	0x06, 0x1d, 0xbd, 0x77, 0x4e, 0xfc, 0xdf, 0x0c, 0xda, 0xfb, 0x98, 0x5b, 0x92, 0xff, 0x08, 0xbd,
	0x58, 0xc9, 0x98, 0x94, 0xde, 0xed, 0xb1, 0x71, 0x2c, 0x1b, 0x2f, 0x65, 0xf3, 0x39, 0x43, 0x94,
	0x18, 0x9d, 0xc5, 0x95, 0x5c, 0xe2, 0x2f, 0x00, 0xab, 0xd0, 0x5b, 0x12, 0xeb, 0x43, 0x83, 0x87,
	0x3a, 0x34, 0x83, 0x5b, 0x49, 0x9b, 0x41, 0x71, 0xf6, 0xff, 0x31, 0x80, 0x2b, 0x65, 0xd1, 0x2f,
	0x6c, 0x74, 0x27, 0x78, 0xa8, 0xfb, 0x94, 0x5b, 0x6b, 0x5f, 0x43, 0xe7, 0x40, 0x71, 0xcb, 0xe7,
	0x3a, 0xc1, 0xdb, 0xfb, 0x82, 0xe3, 0x04, 0x20, 0x1b, 0x58, 0x50, 0xe2, 0x9d, 0x58, 0x71, 0x4a,
	0xad, 0xf2, 0xa9, 0x83, 0x12, 0x0a, 0x9f, 0x40, 0x23, 0xfa, 0x2e, 0x56, 0x5c, 0xd1, 0x26, 0x93,
	0xb3, 0xba, 0x14, 0x05, 0x02, 0x1f, 0xc3, 0xdd, 0x38, 0x54, 0xb4, 0xd1, 0xf9, 0x26, 0x54, 0xc1,
	0x39, 0xc0, 0x97, 0xd0, 0x2a, 0x77, 0xc5, 0x29, 0x9c, 0x1d, 0x31, 0x31, 0x53, 0xe4, 0x7a, 0x0f,
	0xb1, 0xea, 0x21, 0xf6, 0xa0, 0x7e, 0x19, 0xae, 0xb6, 0xb9, 0x5b, 0xe9, 0xc1, 0x7f, 0x98, 0x36,
	0xcc, 0x15, 0x3d, 0xbc, 0x38, 0xfe, 0x2f, 0x06, 0x9d, 0x0b, 0x45, 0xa1, 0xa6, 0x8c, 0x2a, 0xfd,
	0x2c, 0x6c, 0x67, 0x25, 0xdb, 0x11, 0x1c, 0x6b, 0x68, 0xb6, 0x0a, 0x26, 0xc6, 0x31, 0x34, 0xd3,
	0xb9, 0x66, 0x82, 0xdb, 0x45, 0x38, 0x6e, 0x62, 0x23, 0x05, 0x4d, 0x39, 0x3e, 0x82, 0x53, 0x55,
	0xdc, 0x61, 0x53, 0xe4, 0x58, 0x1e, 0xad, 0xab, 0xe4, 0x94, 0x4f, 0xfe, 0xd4, 0xc0, 0x35, 0xf5,
	0x5f, 0x48, 0x5d, 0x8a, 0x88, 0xf0, 0x25, 0x9c, 0x2e, 0x49, 0x97, 0xee, 0xfe, 0xfd, 0x51, 0xfa,
	0x50, 0x8c, 0xf2, 0x87, 0x62, 0xf4, 0xde, 0x3c, 0x14, 0xfd, 0x4c, 0xf6, 0x12, 0xf2, 0x1c, 0x9a,
	0x4b, 0xd2, 0xd9, 0xae, 0x1d, 0xa1, 0xd6, 0xaf, 0x38, 0x85, 0xaf, 0xa0, 0x93, 0x6c, 0xe7, 0x66,
	0x99, 0xe7, 0x74, 0x93, 0xc2, 0x21, 0x7b, 0xca, 0xf0, 0x05, 0xb4, 0xa2, 0x92, 0x96, 0x78, 0x2f,
	0x43, 0xed, 0xeb, 0x7b, 0xa4, 0xeb, 0x73, 0x68, 0x71, 0x5a, 0x91, 0xbe, 0x51, 0xcb, 0xb7, 0xee,
	0xb7, 0xa6, 0x49, 0xa5, 0x32, 0xdc, 0xb1, 0x3f, 0xcf, 0xfe, 0x07, 0x00, 0x00, 0xff, 0xff, 0x70,
	0x83, 0x16, 0xd6, 0x5e, 0x05, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// CmisServiceClient is the client API for CmisService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CmisServiceClient interface {
	GetRepository(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Repository, error)
	GetObject(ctx context.Context, in *CmisObjectId, opts ...grpc.CallOption) (*CmisObject, error)
	SubscribeObject(ctx context.Context, opts ...grpc.CallOption) (CmisService_SubscribeObjectClient, error)
	CreateObject(ctx context.Context, in *CreateObjectReq, opts ...grpc.CallOption) (*CmisObject, error)
	DeleteObject(ctx context.Context, in *CmisObjectId, opts ...grpc.CallOption) (*CmisObject, error)
}

type cmisServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCmisServiceClient(cc grpc.ClientConnInterface) CmisServiceClient {
	return &cmisServiceClient{cc}
}

func (c *cmisServiceClient) GetRepository(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*Repository, error) {
	out := new(Repository)
	err := c.cc.Invoke(ctx, "/cmis.CmisService/getRepository", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmisServiceClient) GetObject(ctx context.Context, in *CmisObjectId, opts ...grpc.CallOption) (*CmisObject, error) {
	out := new(CmisObject)
	err := c.cc.Invoke(ctx, "/cmis.CmisService/getObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmisServiceClient) SubscribeObject(ctx context.Context, opts ...grpc.CallOption) (CmisService_SubscribeObjectClient, error) {
	stream, err := c.cc.NewStream(ctx, &_CmisService_serviceDesc.Streams[0], "/cmis.CmisService/subscribeObject", opts...)
	if err != nil {
		return nil, err
	}
	x := &cmisServiceSubscribeObjectClient{stream}
	return x, nil
}

type CmisService_SubscribeObjectClient interface {
	Send(*CmisObjectId) error
	Recv() (*CmisObject, error)
	grpc.ClientStream
}

type cmisServiceSubscribeObjectClient struct {
	grpc.ClientStream
}

func (x *cmisServiceSubscribeObjectClient) Send(m *CmisObjectId) error {
	return x.ClientStream.SendMsg(m)
}

func (x *cmisServiceSubscribeObjectClient) Recv() (*CmisObject, error) {
	m := new(CmisObject)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *cmisServiceClient) CreateObject(ctx context.Context, in *CreateObjectReq, opts ...grpc.CallOption) (*CmisObject, error) {
	out := new(CmisObject)
	err := c.cc.Invoke(ctx, "/cmis.CmisService/createObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cmisServiceClient) DeleteObject(ctx context.Context, in *CmisObjectId, opts ...grpc.CallOption) (*CmisObject, error) {
	out := new(CmisObject)
	err := c.cc.Invoke(ctx, "/cmis.CmisService/deleteObject", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CmisServiceServer is the server API for CmisService service.
type CmisServiceServer interface {
	GetRepository(context.Context, *empty.Empty) (*Repository, error)
	GetObject(context.Context, *CmisObjectId) (*CmisObject, error)
	SubscribeObject(CmisService_SubscribeObjectServer) error
	CreateObject(context.Context, *CreateObjectReq) (*CmisObject, error)
	DeleteObject(context.Context, *CmisObjectId) (*CmisObject, error)
}

// UnimplementedCmisServiceServer can be embedded to have forward compatible implementations.
type UnimplementedCmisServiceServer struct {
}

func (*UnimplementedCmisServiceServer) GetRepository(ctx context.Context, req *empty.Empty) (*Repository, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRepository not implemented")
}
func (*UnimplementedCmisServiceServer) GetObject(ctx context.Context, req *CmisObjectId) (*CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObject not implemented")
}
func (*UnimplementedCmisServiceServer) SubscribeObject(srv CmisService_SubscribeObjectServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeObject not implemented")
}
func (*UnimplementedCmisServiceServer) CreateObject(ctx context.Context, req *CreateObjectReq) (*CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateObject not implemented")
}
func (*UnimplementedCmisServiceServer) DeleteObject(ctx context.Context, req *CmisObjectId) (*CmisObject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}

func RegisterCmisServiceServer(s *grpc.Server, srv CmisServiceServer) {
	s.RegisterService(&_CmisService_serviceDesc, srv)
}

func _CmisService_GetRepository_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CmisServiceServer).GetRepository(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cmis.CmisService/GetRepository",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CmisServiceServer).GetRepository(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CmisService_GetObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CmisObjectId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CmisServiceServer).GetObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cmis.CmisService/GetObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CmisServiceServer).GetObject(ctx, req.(*CmisObjectId))
	}
	return interceptor(ctx, in, info, handler)
}

func _CmisService_SubscribeObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CmisServiceServer).SubscribeObject(&cmisServiceSubscribeObjectServer{stream})
}

type CmisService_SubscribeObjectServer interface {
	Send(*CmisObject) error
	Recv() (*CmisObjectId, error)
	grpc.ServerStream
}

type cmisServiceSubscribeObjectServer struct {
	grpc.ServerStream
}

func (x *cmisServiceSubscribeObjectServer) Send(m *CmisObject) error {
	return x.ServerStream.SendMsg(m)
}

func (x *cmisServiceSubscribeObjectServer) Recv() (*CmisObjectId, error) {
	m := new(CmisObjectId)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _CmisService_CreateObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateObjectReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CmisServiceServer).CreateObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cmis.CmisService/CreateObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CmisServiceServer).CreateObject(ctx, req.(*CreateObjectReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _CmisService_DeleteObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CmisObjectId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CmisServiceServer).DeleteObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cmis.CmisService/DeleteObject",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CmisServiceServer).DeleteObject(ctx, req.(*CmisObjectId))
	}
	return interceptor(ctx, in, info, handler)
}

var _CmisService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cmis.CmisService",
	HandlerType: (*CmisServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getRepository",
			Handler:    _CmisService_GetRepository_Handler,
		},
		{
			MethodName: "getObject",
			Handler:    _CmisService_GetObject_Handler,
		},
		{
			MethodName: "createObject",
			Handler:    _CmisService_CreateObject_Handler,
		},
		{
			MethodName: "deleteObject",
			Handler:    _CmisService_DeleteObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "subscribeObject",
			Handler:       _CmisService_SubscribeObject_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/cmis.proto",
}
