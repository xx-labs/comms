///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: messages.proto

package messages

import (
	context "context"
	any "github.com/golang/protobuf/ptypes/any"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Generic response message providing an error message from remote servers
type Ack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=Error,proto3" json:"Error,omitempty"`
}

func (x *Ack) Reset() {
	*x = Ack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ack) ProtoMessage() {}

func (x *Ack) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ack.ProtoReflect.Descriptor instead.
func (*Ack) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{0}
}

func (x *Ack) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// Empty message for requesting action from any type of server
type Ping struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Ping) Reset() {
	*x = Ping{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ping) ProtoMessage() {}

func (x *Ping) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ping.ProtoReflect.Descriptor instead.
func (*Ping) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{1}
}

// Wrapper for authenticated messages that also ensure integrity
type AuthenticatedMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID        []byte    `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Signature []byte    `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
	Token     []byte    `protobuf:"bytes,3,opt,name=Token,proto3" json:"Token,omitempty"`
	Client    *ClientID `protobuf:"bytes,4,opt,name=Client,proto3" json:"Client,omitempty"`
	Message   *any.Any  `protobuf:"bytes,5,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *AuthenticatedMessage) Reset() {
	*x = AuthenticatedMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthenticatedMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticatedMessage) ProtoMessage() {}

func (x *AuthenticatedMessage) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticatedMessage.ProtoReflect.Descriptor instead.
func (*AuthenticatedMessage) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{2}
}

func (x *AuthenticatedMessage) GetID() []byte {
	if x != nil {
		return x.ID
	}
	return nil
}

func (x *AuthenticatedMessage) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *AuthenticatedMessage) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

func (x *AuthenticatedMessage) GetClient() *ClientID {
	if x != nil {
		return x.Client
	}
	return nil
}

func (x *AuthenticatedMessage) GetMessage() *any.Any {
	if x != nil {
		return x.Message
	}
	return nil
}

// Message used for assembly of Client IDs in the system
type ClientID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Salt      []byte `protobuf:"bytes,1,opt,name=Salt,proto3" json:"Salt,omitempty"`
	PublicKey string `protobuf:"bytes,2,opt,name=PublicKey,proto3" json:"PublicKey,omitempty"`
}

func (x *ClientID) Reset() {
	*x = ClientID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientID) ProtoMessage() {}

func (x *ClientID) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientID.ProtoReflect.Descriptor instead.
func (*ClientID) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{3}
}

func (x *ClientID) GetSalt() []byte {
	if x != nil {
		return x.Salt
	}
	return nil
}

func (x *ClientID) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

// Provides a token to establish reverse identity to any type of client
type AssignToken struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token []byte `protobuf:"bytes,1,opt,name=Token,proto3" json:"Token,omitempty"`
}

func (x *AssignToken) Reset() {
	*x = AssignToken{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssignToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssignToken) ProtoMessage() {}

func (x *AssignToken) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AssignToken.ProtoReflect.Descriptor instead.
func (*AssignToken) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{4}
}

func (x *AssignToken) GetToken() []byte {
	if x != nil {
		return x.Token
	}
	return nil
}

// RSASignature is a digital signature for the RSA algorithm
type RSASignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce     []byte `protobuf:"bytes,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *RSASignature) Reset() {
	*x = RSASignature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RSASignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RSASignature) ProtoMessage() {}

func (x *RSASignature) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RSASignature.ProtoReflect.Descriptor instead.
func (*RSASignature) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{5}
}

func (x *RSASignature) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *RSASignature) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// ECCSignature is a digital signature for the ECC algorithm
type ECCSignature struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nonce     []byte `protobuf:"bytes,1,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	Signature []byte `protobuf:"bytes,2,opt,name=Signature,proto3" json:"Signature,omitempty"`
}

func (x *ECCSignature) Reset() {
	*x = ECCSignature{}
	if protoimpl.UnsafeEnabled {
		mi := &file_messages_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ECCSignature) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ECCSignature) ProtoMessage() {}

func (x *ECCSignature) ProtoReflect() protoreflect.Message {
	mi := &file_messages_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ECCSignature.ProtoReflect.Descriptor instead.
func (*ECCSignature) Descriptor() ([]byte, []int) {
	return file_messages_proto_rawDescGZIP(), []int{6}
}

func (x *ECCSignature) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *ECCSignature) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

var File_messages_proto protoreflect.FileDescriptor

var file_messages_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1b, 0x0a, 0x03, 0x41, 0x63, 0x6b, 0x12, 0x14, 0x0a, 0x05,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x06, 0x0a, 0x04, 0x50, 0x69, 0x6e, 0x67, 0x22, 0xb6, 0x01, 0x0a, 0x14, 0x41,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72,
	0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2a, 0x0a, 0x06, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x06, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x3c, 0x0a, 0x08, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12,
	0x12, 0x0a, 0x04, 0x53, 0x61, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x53,
	0x61, 0x6c, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x22, 0x23, 0x0a, 0x0b, 0x41, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x42, 0x0a, 0x0c, 0x52, 0x53, 0x41, 0x53, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x42, 0x0a, 0x0c, 0x45, 0x43,
	0x43, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f,
	0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x09, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x32, 0x88,
	0x01, 0x0a, 0x07, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x69, 0x63, 0x12, 0x44, 0x0a, 0x11, 0x41, 0x75,
	0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12,
	0x1e, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x65,
	0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x64, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a,
	0x0d, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x41, 0x63, 0x6b, 0x22, 0x00,
	0x12, 0x37, 0x0a, 0x0c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x0e, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x50, 0x69, 0x6e, 0x67,
	0x1a, 0x15, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x69,
	0x67, 0x6e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x00, 0x42, 0x26, 0x5a, 0x24, 0x67, 0x69, 0x74,
	0x6c, 0x61, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x78, 0x5f, 0x6e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_messages_proto_rawDescOnce sync.Once
	file_messages_proto_rawDescData = file_messages_proto_rawDesc
)

func file_messages_proto_rawDescGZIP() []byte {
	file_messages_proto_rawDescOnce.Do(func() {
		file_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_messages_proto_rawDescData)
	})
	return file_messages_proto_rawDescData
}

var file_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_messages_proto_goTypes = []interface{}{
	(*Ack)(nil),                  // 0: messages.Ack
	(*Ping)(nil),                 // 1: messages.Ping
	(*AuthenticatedMessage)(nil), // 2: messages.AuthenticatedMessage
	(*ClientID)(nil),             // 3: messages.ClientID
	(*AssignToken)(nil),          // 4: messages.AssignToken
	(*RSASignature)(nil),         // 5: messages.RSASignature
	(*ECCSignature)(nil),         // 6: messages.ECCSignature
	(*any.Any)(nil),              // 7: google.protobuf.Any
}
var file_messages_proto_depIdxs = []int32{
	3, // 0: messages.AuthenticatedMessage.Client:type_name -> messages.ClientID
	7, // 1: messages.AuthenticatedMessage.Message:type_name -> google.protobuf.Any
	2, // 2: messages.Generic.AuthenticateToken:input_type -> messages.AuthenticatedMessage
	1, // 3: messages.Generic.RequestToken:input_type -> messages.Ping
	0, // 4: messages.Generic.AuthenticateToken:output_type -> messages.Ack
	4, // 5: messages.Generic.RequestToken:output_type -> messages.AssignToken
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_messages_proto_init() }
func file_messages_proto_init() {
	if File_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ack); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ping); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthenticatedMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssignToken); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RSASignature); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_messages_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ECCSignature); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_messages_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_messages_proto_goTypes,
		DependencyIndexes: file_messages_proto_depIdxs,
		MessageInfos:      file_messages_proto_msgTypes,
	}.Build()
	File_messages_proto = out.File
	file_messages_proto_rawDesc = nil
	file_messages_proto_goTypes = nil
	file_messages_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

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
	cc grpc.ClientConnInterface
}

func NewGenericClient(cc grpc.ClientConnInterface) GenericClient {
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

func (*UnimplementedGenericServer) AuthenticateToken(context.Context, *AuthenticatedMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthenticateToken not implemented")
}
func (*UnimplementedGenericServer) RequestToken(context.Context, *Ping) (*AssignToken, error) {
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
