// Code generated by protoc-gen-go. DO NOT EDIT.
// source: messages.proto

package messages

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
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

// Generic response message providing an error message from remote servers
type Ack struct {
	Error                string   `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{0}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

func (m *Ack) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

// Empty message for requesting action from any type of server
type Ping struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{1}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

// Wrapper for authenticated messages that also ensure integrity
type AuthenticatedMessage struct {
	ID                   []byte    `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Signature            []byte    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Token                []byte    `protobuf:"bytes,3,opt,name=Token,proto3" json:"Token,omitempty"`
	Client               *ClientID `protobuf:"bytes,4,opt,name=Client,proto3" json:"Client,omitempty"`
	Message              *any.Any  `protobuf:"bytes,5,opt,name=Message,proto3" json:"Message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *AuthenticatedMessage) Reset()         { *m = AuthenticatedMessage{} }
func (m *AuthenticatedMessage) String() string { return proto.CompactTextString(m) }
func (*AuthenticatedMessage) ProtoMessage()    {}
func (*AuthenticatedMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{2}
}

func (m *AuthenticatedMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AuthenticatedMessage.Unmarshal(m, b)
}
func (m *AuthenticatedMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AuthenticatedMessage.Marshal(b, m, deterministic)
}
func (m *AuthenticatedMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthenticatedMessage.Merge(m, src)
}
func (m *AuthenticatedMessage) XXX_Size() int {
	return xxx_messageInfo_AuthenticatedMessage.Size(m)
}
func (m *AuthenticatedMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthenticatedMessage.DiscardUnknown(m)
}

var xxx_messageInfo_AuthenticatedMessage proto.InternalMessageInfo

func (m *AuthenticatedMessage) GetID() []byte {
	if m != nil {
		return m.ID
	}
	return nil
}

func (m *AuthenticatedMessage) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *AuthenticatedMessage) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

func (m *AuthenticatedMessage) GetClient() *ClientID {
	if m != nil {
		return m.Client
	}
	return nil
}

func (m *AuthenticatedMessage) GetMessage() *any.Any {
	if m != nil {
		return m.Message
	}
	return nil
}

// Message used for assembly of Client IDs in the system
type ClientID struct {
	Salt                 []byte   `protobuf:"bytes,1,opt,name=Salt,proto3" json:"Salt,omitempty"`
	PublicKey            string   `protobuf:"bytes,2,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClientID) Reset()         { *m = ClientID{} }
func (m *ClientID) String() string { return proto.CompactTextString(m) }
func (*ClientID) ProtoMessage()    {}
func (*ClientID) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{3}
}

func (m *ClientID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClientID.Unmarshal(m, b)
}
func (m *ClientID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClientID.Marshal(b, m, deterministic)
}
func (m *ClientID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClientID.Merge(m, src)
}
func (m *ClientID) XXX_Size() int {
	return xxx_messageInfo_ClientID.Size(m)
}
func (m *ClientID) XXX_DiscardUnknown() {
	xxx_messageInfo_ClientID.DiscardUnknown(m)
}

var xxx_messageInfo_ClientID proto.InternalMessageInfo

func (m *ClientID) GetSalt() []byte {
	if m != nil {
		return m.Salt
	}
	return nil
}

func (m *ClientID) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

// Provides a token to establish reverse identity to any type of client
type AssignToken struct {
	Token                []byte   `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AssignToken) Reset()         { *m = AssignToken{} }
func (m *AssignToken) String() string { return proto.CompactTextString(m) }
func (*AssignToken) ProtoMessage()    {}
func (*AssignToken) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{4}
}

func (m *AssignToken) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AssignToken.Unmarshal(m, b)
}
func (m *AssignToken) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AssignToken.Marshal(b, m, deterministic)
}
func (m *AssignToken) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AssignToken.Merge(m, src)
}
func (m *AssignToken) XXX_Size() int {
	return xxx_messageInfo_AssignToken.Size(m)
}
func (m *AssignToken) XXX_DiscardUnknown() {
	xxx_messageInfo_AssignToken.DiscardUnknown(m)
}

var xxx_messageInfo_AssignToken proto.InternalMessageInfo

func (m *AssignToken) GetToken() []byte {
	if m != nil {
		return m.Token
	}
	return nil
}

// RSASignature is a digital signature for the RSA algorithm
type RSASignature struct {
	Nonce                []byte   `protobuf:"bytes,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Signature            []byte   `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RSASignature) Reset()         { *m = RSASignature{} }
func (m *RSASignature) String() string { return proto.CompactTextString(m) }
func (*RSASignature) ProtoMessage()    {}
func (*RSASignature) Descriptor() ([]byte, []int) {
	return fileDescriptor_4dc296cbfe5ffcd5, []int{5}
}

func (m *RSASignature) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RSASignature.Unmarshal(m, b)
}
func (m *RSASignature) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RSASignature.Marshal(b, m, deterministic)
}
func (m *RSASignature) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RSASignature.Merge(m, src)
}
func (m *RSASignature) XXX_Size() int {
	return xxx_messageInfo_RSASignature.Size(m)
}
func (m *RSASignature) XXX_DiscardUnknown() {
	xxx_messageInfo_RSASignature.DiscardUnknown(m)
}

var xxx_messageInfo_RSASignature proto.InternalMessageInfo

func (m *RSASignature) GetNonce() []byte {
	if m != nil {
		return m.Nonce
	}
	return nil
}

func (m *RSASignature) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*Ack)(nil), "messages.Ack")
	proto.RegisterType((*Ping)(nil), "messages.Ping")
	proto.RegisterType((*AuthenticatedMessage)(nil), "messages.AuthenticatedMessage")
	proto.RegisterType((*ClientID)(nil), "messages.ClientID")
	proto.RegisterType((*AssignToken)(nil), "messages.AssignToken")
	proto.RegisterType((*RSASignature)(nil), "messages.RSASignature")
}

func init() { proto.RegisterFile("messages.proto", fileDescriptor_4dc296cbfe5ffcd5) }

var fileDescriptor_4dc296cbfe5ffcd5 = []byte{
	// 340 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x51, 0xdd, 0x4a, 0xc3, 0x30,
	0x18, 0x5d, 0xf7, 0xbf, 0x6f, 0x73, 0x60, 0x98, 0x50, 0xa7, 0xc8, 0x88, 0x37, 0xc3, 0x8b, 0x0e,
	0xe6, 0x85, 0x37, 0xde, 0x54, 0x2b, 0x32, 0x44, 0x19, 0x99, 0x2f, 0xd0, 0xd5, 0xcf, 0x1a, 0x56,
	0x13, 0x6d, 0xd2, 0x8b, 0xbd, 0x81, 0x4f, 0xe4, 0xf3, 0xc9, 0x92, 0x76, 0xad, 0x20, 0xde, 0xe5,
	0x9c, 0x9c, 0x7c, 0xdf, 0x39, 0x27, 0x30, 0x7c, 0x47, 0xa5, 0xc2, 0x18, 0x95, 0xf7, 0x91, 0x4a,
	0x2d, 0x49, 0xb7, 0xc0, 0xe3, 0xe3, 0x58, 0xca, 0x38, 0xc1, 0x99, 0xe1, 0xd7, 0xd9, 0xeb, 0x2c,
	0x14, 0x5b, 0x2b, 0xa2, 0x27, 0xd0, 0xf0, 0xa3, 0x0d, 0x19, 0x41, 0xeb, 0x2e, 0x4d, 0x65, 0xea,
	0x3a, 0x13, 0x67, 0xda, 0x63, 0x16, 0xd0, 0x36, 0x34, 0x97, 0x5c, 0xc4, 0xf4, 0xdb, 0x81, 0x91,
	0x9f, 0xe9, 0x37, 0x14, 0x9a, 0x47, 0xa1, 0xc6, 0x97, 0x47, 0x3b, 0x99, 0x0c, 0xa1, 0xbe, 0x08,
	0xcc, 0x9b, 0x01, 0xab, 0x2f, 0x02, 0x72, 0x0a, 0xbd, 0x15, 0x8f, 0x45, 0xa8, 0xb3, 0x14, 0xdd,
	0xba, 0xa1, 0x4b, 0x62, 0xb7, 0xe4, 0x59, 0x6e, 0x50, 0xb8, 0x0d, 0x73, 0x63, 0x01, 0xb9, 0x80,
	0xf6, 0x6d, 0xc2, 0x51, 0x68, 0xb7, 0x39, 0x71, 0xa6, 0xfd, 0x39, 0xf1, 0xf6, 0x39, 0x2c, 0xbf,
	0x08, 0x58, 0xae, 0x20, 0x1e, 0x74, 0xf2, 0xd5, 0x6e, 0xcb, 0x88, 0x47, 0x9e, 0x8d, 0xe6, 0x15,
	0xd1, 0x3c, 0x5f, 0x6c, 0x59, 0x21, 0xa2, 0xd7, 0xd0, 0x2d, 0x66, 0x10, 0x02, 0xcd, 0x55, 0x98,
	0xe8, 0xdc, 0xad, 0x39, 0xef, 0xfc, 0x2e, 0xb3, 0x75, 0xc2, 0xa3, 0x07, 0xdc, 0x1a, 0xbf, 0x3d,
	0x56, 0x12, 0xf4, 0x1c, 0xfa, 0xbe, 0x52, 0x3c, 0x16, 0xd6, 0xe8, 0xde, 0xbe, 0x53, 0xb1, 0x4f,
	0x6f, 0x60, 0xc0, 0x56, 0xfe, 0xaf, 0x90, 0x4f, 0x52, 0x44, 0x58, 0xa8, 0x0c, 0xf8, 0xbf, 0x98,
	0xf9, 0x97, 0x03, 0x9d, 0x7b, 0x14, 0x98, 0xf2, 0x88, 0x04, 0x70, 0x58, 0xad, 0xda, 0xae, 0x3e,
	0x2b, 0x3b, 0xf9, 0xeb, 0x1f, 0xc6, 0x07, 0x95, 0xfb, 0x68, 0x43, 0x6b, 0xe4, 0x0a, 0x06, 0x0c,
	0x3f, 0x33, 0x54, 0xda, 0x0e, 0x18, 0x96, 0x82, 0xdd, 0x8f, 0x8e, 0x8f, 0x2a, 0x0f, 0xca, 0x88,
	0xb4, 0xb6, 0x6e, 0x9b, 0x22, 0x2f, 0x7f, 0x02, 0x00, 0x00, 0xff, 0xff, 0xfc, 0xe5, 0xbd, 0x15,
	0x4d, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GenericClient is the client API for Generic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GenericClient interface {
	// Authenticate a token with the server
	AuthenticateToken(ctx context.Context, in *AuthenticatedMessage, opts ...grpc.CallOption) (*Ack, error)
	// Request a token from the server
	RequestToken(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*AssignToken, error)
}

type genericClient struct {
	cc *grpc.ClientConn
}

func NewGenericClient(cc *grpc.ClientConn) GenericClient {
	return &genericClient{cc}
}

func (c *genericClient) AuthenticateToken(ctx context.Context, in *AuthenticatedMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/messages.Generic/AuthenticateToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genericClient) RequestToken(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*AssignToken, error) {
	out := new(AssignToken)
	err := c.cc.Invoke(ctx, "/messages.Generic/RequestToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GenericServer is the server API for Generic service.
type GenericServer interface {
	// Authenticate a token with the server
	AuthenticateToken(context.Context, *AuthenticatedMessage) (*Ack, error)
	// Request a token from the server
	RequestToken(context.Context, *Ping) (*AssignToken, error)
}

// UnimplementedGenericServer can be embedded to have forward compatible implementations.
type UnimplementedGenericServer struct {
}

func (*UnimplementedGenericServer) AuthenticateToken(ctx context.Context, req *AuthenticatedMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateToken not implemented")
}
func (*UnimplementedGenericServer) RequestToken(ctx context.Context, req *Ping) (*AssignToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestToken not implemented")
}

func RegisterGenericServer(s *grpc.Server, srv GenericServer) {
	s.RegisterService(&_Generic_serviceDesc, srv)
}

func _Generic_AuthenticateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthenticatedMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericServer).AuthenticateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Generic/AuthenticateToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericServer).AuthenticateToken(ctx, req.(*AuthenticatedMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Generic_RequestToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenericServer).RequestToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/messages.Generic/RequestToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenericServer).RequestToken(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

var _Generic_serviceDesc = grpc.ServiceDesc{
	ServiceName: "messages.Generic",
	HandlerType: (*GenericServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthenticateToken",
			Handler:    _Generic_AuthenticateToken_Handler,
		},
		{
			MethodName: "RequestToken",
			Handler:    _Generic_RequestToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "messages.proto",
}
