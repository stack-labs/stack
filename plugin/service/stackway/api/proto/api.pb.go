// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/stack-labs/stack/plugin/service/stackway/api/proto/api.proto

package stack_rpc_api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type Pair struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values               []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Pair) Reset()         { *m = Pair{} }
func (m *Pair) String() string { return proto.CompactTextString(m) }
func (*Pair) ProtoMessage()    {}
func (*Pair) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e9b8f29cbb979d, []int{0}
}

func (m *Pair) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Pair.Unmarshal(m, b)
}
func (m *Pair) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Pair.Marshal(b, m, deterministic)
}
func (m *Pair) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Pair.Merge(m, src)
}
func (m *Pair) XXX_Size() int {
	return xxx_messageInfo_Pair.Size(m)
}
func (m *Pair) XXX_DiscardUnknown() {
	xxx_messageInfo_Pair.DiscardUnknown(m)
}

var xxx_messageInfo_Pair proto.InternalMessageInfo

func (m *Pair) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *Pair) GetValues() []string {
	if m != nil {
		return m.Values
	}
	return nil
}

type Request struct {
	Method               string           `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Path                 string           `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Header               map[string]*Pair `protobuf:"bytes,3,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Get                  map[string]*Pair `protobuf:"bytes,4,rep,name=get,proto3" json:"get,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Post                 map[string]*Pair `protobuf:"bytes,5,rep,name=post,proto3" json:"post,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 string           `protobuf:"bytes,6,opt,name=body,proto3" json:"body,omitempty"`
	Url                  string           `protobuf:"bytes,7,opt,name=url,proto3" json:"url,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e9b8f29cbb979d, []int{1}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *Request) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Request) GetHeader() map[string]*Pair {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Request) GetGet() map[string]*Pair {
	if m != nil {
		return m.Get
	}
	return nil
}

func (m *Request) GetPost() map[string]*Pair {
	if m != nil {
		return m.Post
	}
	return nil
}

func (m *Request) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *Request) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

type Response struct {
	StatusCode           int32            `protobuf:"varint,1,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Header               map[string]*Pair `protobuf:"bytes,2,rep,name=header,proto3" json:"header,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Body                 string           `protobuf:"bytes,3,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_e5e9b8f29cbb979d, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *Response) GetHeader() map[string]*Pair {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Response) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func init() {
	proto.RegisterType((*Pair)(nil), "stack.rpc.api.Pair")
	proto.RegisterType((*Request)(nil), "stack.rpc.api.Request")
	proto.RegisterMapType((map[string]*Pair)(nil), "stack.rpc.api.Request.GetEntry")
	proto.RegisterMapType((map[string]*Pair)(nil), "stack.rpc.api.Request.HeaderEntry")
	proto.RegisterMapType((map[string]*Pair)(nil), "stack.rpc.api.Request.PostEntry")
	proto.RegisterType((*Response)(nil), "stack.rpc.api.Response")
	proto.RegisterMapType((map[string]*Pair)(nil), "stack.rpc.api.Response.HeaderEntry")
}

func init() {
	proto.RegisterFile("github.com/stack-labs/stack/plugin/service/stackway/api/proto/api.proto", fileDescriptor_e5e9b8f29cbb979d)
}

var fileDescriptor_e5e9b8f29cbb979d = []byte{
	// 375 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0x4f, 0x6e, 0xe2, 0x30,
	0x18, 0xc5, 0x95, 0x3f, 0x04, 0xf8, 0xd0, 0x48, 0x23, 0x8f, 0x34, 0xb2, 0x58, 0xcc, 0x44, 0x99,
	0x0d, 0xb3, 0x20, 0x69, 0x69, 0x17, 0x15, 0x5d, 0x56, 0x55, 0x2b, 0x15, 0x55, 0x28, 0x37, 0x70,
	0x12, 0x8b, 0x44, 0x04, 0xec, 0xda, 0x0e, 0x55, 0xce, 0xd8, 0xa3, 0xf4, 0x12, 0x95, 0x9d, 0x40,
	0x69, 0x0b, 0x2b, 0xba, 0x7b, 0xb6, 0xbf, 0xf7, 0xf2, 0xf2, 0xb3, 0x0c, 0xb3, 0x45, 0xa1, 0xf2,
	0x2a, 0x09, 0x53, 0xb6, 0x8a, 0xa4, 0x22, 0xe9, 0x72, 0x5c, 0x92, 0x44, 0xb6, 0x52, 0xf0, 0x74,
	0xcc, 0xcb, 0x6a, 0x51, 0xac, 0x65, 0x24, 0xa9, 0xd8, 0x14, 0x29, 0x6d, 0x4e, 0x9e, 0x49, 0x1d,
	0x11, 0x5e, 0x44, 0x5c, 0x30, 0xc5, 0xb4, 0x0a, 0x8d, 0x42, 0x3f, 0xcc, 0x69, 0x28, 0x78, 0x1a,
	0x12, 0x5e, 0x04, 0x67, 0xe0, 0xce, 0x49, 0x21, 0xd0, 0x4f, 0x70, 0x96, 0xb4, 0xc6, 0x96, 0x6f,
	0x8d, 0xfa, 0xb1, 0x96, 0xe8, 0x37, 0x78, 0x1b, 0x52, 0x56, 0x54, 0x62, 0xdb, 0x77, 0x46, 0xfd,
	0xb8, 0x5d, 0x05, 0xaf, 0x0e, 0x74, 0x63, 0xfa, 0x54, 0x51, 0xa9, 0xf4, 0xcc, 0x8a, 0xaa, 0x9c,
	0x65, 0xad, 0xb1, 0x5d, 0x21, 0x04, 0x2e, 0x27, 0x2a, 0xc7, 0xb6, 0xd9, 0x35, 0x1a, 0x4d, 0xc1,
	0xcb, 0x29, 0xc9, 0xa8, 0xc0, 0x8e, 0xef, 0x8c, 0x06, 0x93, 0x20, 0xfc, 0xd0, 0x24, 0x6c, 0x33,
	0xc3, 0x7b, 0x33, 0x74, 0xbb, 0x56, 0xa2, 0x8e, 0x5b, 0x07, 0x3a, 0x07, 0x67, 0x41, 0x15, 0x76,
	0x8d, 0xf1, 0xef, 0x11, 0xe3, 0x1d, 0x55, 0x8d, 0x4b, 0xcf, 0xa2, 0x4b, 0x70, 0x39, 0x93, 0x0a,
	0x77, 0x8c, 0xc7, 0x3f, 0xe2, 0x99, 0x33, 0xd9, 0x9a, 0xcc, 0xb4, 0x2e, 0x9e, 0xb0, 0xac, 0xc6,
	0x5e, 0x53, 0x5c, 0x6b, 0x8d, 0xa6, 0x12, 0x25, 0xee, 0x36, 0x68, 0x2a, 0x51, 0x0e, 0x1f, 0x61,
	0xb0, 0xd7, 0xf2, 0x00, 0xbb, 0xff, 0xd0, 0x31, 0xb4, 0x0c, 0x80, 0xc1, 0xe4, 0xd7, 0xa7, 0xaf,
	0x6b, 0xe2, 0x71, 0x33, 0x31, 0xb5, 0xaf, 0xac, 0xe1, 0x03, 0xf4, 0xb6, 0xe5, 0x4f, 0x0f, 0x9b,
	0x41, 0x7f, 0xf7, 0x57, 0x27, 0xa7, 0x05, 0x2f, 0x16, 0xf4, 0x62, 0x2a, 0x39, 0x5b, 0x4b, 0x8a,
	0xfe, 0x00, 0x48, 0x45, 0x54, 0x25, 0x6f, 0x58, 0x46, 0x4d, 0x68, 0x27, 0xde, 0xdb, 0x41, 0xd7,
	0xbb, 0x2b, 0xb6, 0x0d, 0xf5, 0x7f, 0x5f, 0xa8, 0x37, 0x41, 0x07, 0xef, 0x78, 0x8b, 0xde, 0x79,
	0x47, 0xff, 0xdd, 0xa0, 0x13, 0xcf, 0xbc, 0x81, 0x8b, 0xb7, 0x00, 0x00, 0x00, 0xff, 0xff, 0x49,
	0xc3, 0x7d, 0x73, 0x53, 0x03, 0x00, 0x00,
}