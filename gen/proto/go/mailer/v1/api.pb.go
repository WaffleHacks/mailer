// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: mailer/v1/api.proto

package mailerv1

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type BodyType int32

const (
	BodyType_BODY_TYPE_UNSPECIFIED BodyType = 0
	BodyType_BODY_TYPE_PLAIN       BodyType = 1
	BodyType_BODY_TYPE_HTML        BodyType = 2
)

// Enum value maps for BodyType.
var (
	BodyType_name = map[int32]string{
		0: "BODY_TYPE_UNSPECIFIED",
		1: "BODY_TYPE_PLAIN",
		2: "BODY_TYPE_HTML",
	}
	BodyType_value = map[string]int32{
		"BODY_TYPE_UNSPECIFIED": 0,
		"BODY_TYPE_PLAIN":       1,
		"BODY_TYPE_HTML":        2,
	}
)

func (x BodyType) Enum() *BodyType {
	p := new(BodyType)
	*p = x
	return p
}

func (x BodyType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BodyType) Descriptor() protoreflect.EnumDescriptor {
	return file_mailer_v1_api_proto_enumTypes[0].Descriptor()
}

func (BodyType) Type() protoreflect.EnumType {
	return &file_mailer_v1_api_proto_enumTypes[0]
}

func (x BodyType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BodyType.Descriptor instead.
func (BodyType) EnumDescriptor() ([]byte, []int) {
	return file_mailer_v1_api_proto_rawDescGZIP(), []int{0}
}

type SendRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To      string   `protobuf:"bytes,1,opt,name=to,proto3" json:"to,omitempty"`
	From    string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Subject string   `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	Body    string   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Type    BodyType `protobuf:"varint,5,opt,name=type,proto3,enum=mailer.v1.BodyType" json:"type,omitempty"`
	ReplyTo string   `protobuf:"bytes,6,opt,name=reply_to,json=replyTo,proto3" json:"reply_to,omitempty"`
}

func (x *SendRequest) Reset() {
	*x = SendRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mailer_v1_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendRequest) ProtoMessage() {}

func (x *SendRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mailer_v1_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendRequest.ProtoReflect.Descriptor instead.
func (*SendRequest) Descriptor() ([]byte, []int) {
	return file_mailer_v1_api_proto_rawDescGZIP(), []int{0}
}

func (x *SendRequest) GetTo() string {
	if x != nil {
		return x.To
	}
	return ""
}

func (x *SendRequest) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SendRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *SendRequest) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *SendRequest) GetType() BodyType {
	if x != nil {
		return x.Type
	}
	return BodyType_BODY_TYPE_UNSPECIFIED
}

func (x *SendRequest) GetReplyTo() string {
	if x != nil {
		return x.ReplyTo
	}
	return ""
}

type SendResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendResponse) Reset() {
	*x = SendResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mailer_v1_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendResponse) ProtoMessage() {}

func (x *SendResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mailer_v1_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendResponse.ProtoReflect.Descriptor instead.
func (*SendResponse) Descriptor() ([]byte, []int) {
	return file_mailer_v1_api_proto_rawDescGZIP(), []int{1}
}

type SendBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	To      []string `protobuf:"bytes,1,rep,name=to,proto3" json:"to,omitempty"`
	From    string   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	Subject string   `protobuf:"bytes,3,opt,name=subject,proto3" json:"subject,omitempty"`
	Body    string   `protobuf:"bytes,4,opt,name=body,proto3" json:"body,omitempty"`
	Type    BodyType `protobuf:"varint,5,opt,name=type,proto3,enum=mailer.v1.BodyType" json:"type,omitempty"`
	ReplyTo string   `protobuf:"bytes,6,opt,name=reply_to,json=replyTo,proto3" json:"reply_to,omitempty"`
}

func (x *SendBatchRequest) Reset() {
	*x = SendBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mailer_v1_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendBatchRequest) ProtoMessage() {}

func (x *SendBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_mailer_v1_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendBatchRequest.ProtoReflect.Descriptor instead.
func (*SendBatchRequest) Descriptor() ([]byte, []int) {
	return file_mailer_v1_api_proto_rawDescGZIP(), []int{2}
}

func (x *SendBatchRequest) GetTo() []string {
	if x != nil {
		return x.To
	}
	return nil
}

func (x *SendBatchRequest) GetFrom() string {
	if x != nil {
		return x.From
	}
	return ""
}

func (x *SendBatchRequest) GetSubject() string {
	if x != nil {
		return x.Subject
	}
	return ""
}

func (x *SendBatchRequest) GetBody() string {
	if x != nil {
		return x.Body
	}
	return ""
}

func (x *SendBatchRequest) GetType() BodyType {
	if x != nil {
		return x.Type
	}
	return BodyType_BODY_TYPE_UNSPECIFIED
}

func (x *SendBatchRequest) GetReplyTo() string {
	if x != nil {
		return x.ReplyTo
	}
	return ""
}

type SendBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SendBatchResponse) Reset() {
	*x = SendBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mailer_v1_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SendBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendBatchResponse) ProtoMessage() {}

func (x *SendBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_mailer_v1_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendBatchResponse.ProtoReflect.Descriptor instead.
func (*SendBatchResponse) Descriptor() ([]byte, []int) {
	return file_mailer_v1_api_proto_rawDescGZIP(), []int{3}
}

var File_mailer_v1_api_proto protoreflect.FileDescriptor

var file_mailer_v1_api_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x70, 0x69, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e,
	0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e,
	0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xcd, 0x03, 0x0a, 0x0b, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x3a, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2a, 0xe0, 0x41, 0x02,
	0x92, 0x41, 0x24, 0x32, 0x22, 0x54, 0x68, 0x65, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x4e, 0x0a, 0x04, 0x66,
	0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x3a, 0xe0, 0x41, 0x02, 0x92, 0x41,
	0x34, 0x32, 0x32, 0x54, 0x68, 0x65, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x20, 0x69, 0x6e, 0x20, 0x52, 0x46, 0x43, 0x20, 0x35, 0x33, 0x32, 0x32, 0x20, 0x66,
	0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x41, 0x0a, 0x07, 0x73,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x27, 0xe0, 0x41,
	0x02, 0x92, 0x41, 0x21, 0x32, 0x1f, 0x54, 0x68, 0x65, 0x20, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x20, 0x6c, 0x69, 0x6e, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x3b,
	0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x27, 0xe0, 0x41,
	0x02, 0x92, 0x41, 0x21, 0x32, 0x1f, 0x54, 0x68, 0x65, 0x20, 0x6d, 0x61, 0x69, 0x6e, 0x20, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x6a, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6c,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x42, 0x6f, 0x64, 0x79, 0x54, 0x79, 0x70, 0x65, 0x42, 0x41,
	0xe0, 0x41, 0x01, 0x92, 0x41, 0x3b, 0x32, 0x39, 0x54, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x20, 0x74, 0x79, 0x70, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2c, 0x20, 0x65, 0x69, 0x74, 0x68, 0x65, 0x72, 0x20,
	0x48, 0x54, 0x4d, 0x4c, 0x20, 0x6f, 0x72, 0x20, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78,
	0x74, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x46, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x79,
	0x5f, 0x74, 0x6f, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2b, 0xe0, 0x41, 0x01, 0x92, 0x41,
	0x25, 0x32, 0x23, 0x54, 0x68, 0x65, 0x20, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x20, 0x66,
	0x6f, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2d, 0x54, 0x6f, 0x20,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x22,
	0x0e, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0xd7, 0x03, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x3f, 0x0a, 0x02, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x42, 0x2f, 0xe0, 0x41, 0x02, 0x92, 0x41, 0x29, 0x32, 0x27, 0x41, 0x20, 0x6c, 0x69, 0x73, 0x74,
	0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65,
	0x73, 0x52, 0x02, 0x74, 0x6f, 0x12, 0x4e, 0x0a, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x3a, 0xe0, 0x41, 0x02, 0x92, 0x41, 0x34, 0x32, 0x32, 0x54, 0x68, 0x65,
	0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x20, 0x6f,
	0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x20, 0x69, 0x6e, 0x20,
	0x52, 0x46, 0x43, 0x20, 0x35, 0x33, 0x32, 0x32, 0x20, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x52,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x41, 0x0a, 0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x27, 0xe0, 0x41, 0x02, 0x92, 0x41, 0x21, 0x32, 0x1f,
	0x54, 0x68, 0x65, 0x20, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x20, 0x6c, 0x69, 0x6e, 0x65,
	0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x07, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x3b, 0x0a, 0x04, 0x62, 0x6f, 0x64, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x27, 0xe0, 0x41, 0x02, 0x92, 0x41, 0x21, 0x32, 0x1f,
	0x54, 0x68, 0x65, 0x20, 0x6d, 0x61, 0x69, 0x6e, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74,
	0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x52,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x6a, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x42, 0x6f, 0x64, 0x79, 0x54, 0x79, 0x70, 0x65, 0x42, 0x41, 0xe0, 0x41, 0x01, 0x92, 0x41, 0x3b,
	0x32, 0x39, 0x54, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x20, 0x74, 0x79,
	0x70, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x2c, 0x20, 0x65, 0x69, 0x74, 0x68, 0x65, 0x72, 0x20, 0x48, 0x54, 0x4d, 0x4c, 0x20, 0x6f,
	0x72, 0x20, 0x70, 0x6c, 0x61, 0x69, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x46, 0x0a, 0x08, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x5f, 0x74, 0x6f, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x2b, 0xe0, 0x41, 0x01, 0x92, 0x41, 0x25, 0x32, 0x23, 0x54, 0x68, 0x65,
	0x20, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x74, 0x68, 0x65,
	0x20, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2d, 0x54, 0x6f, 0x20, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x52, 0x07, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x54, 0x6f, 0x22, 0x13, 0x0a, 0x11, 0x53, 0x65, 0x6e,
	0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2a, 0x4e,
	0x0a, 0x08, 0x42, 0x6f, 0x64, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x15, 0x42, 0x4f,
	0x44, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46,
	0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x42, 0x4f, 0x44, 0x59, 0x5f, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x50, 0x4c, 0x41, 0x49, 0x4e, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x42, 0x4f,
	0x44, 0x59, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x48, 0x54, 0x4d, 0x4c, 0x10, 0x02, 0x32, 0xac,
	0x02, 0x0a, 0x0d, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x7e, 0x0a, 0x04, 0x53, 0x65, 0x6e, 0x64, 0x12, 0x16, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65,
	0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x17, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x45, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x0a, 0x3a, 0x01, 0x2a, 0x22, 0x05, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x92, 0x41, 0x32, 0x12, 0x04,
	0x53, 0x65, 0x6e, 0x64, 0x1a, 0x2a, 0x53, 0x65, 0x6e, 0x64, 0x20, 0x61, 0x20, 0x73, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x20, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x20, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x20, 0x74, 0x6f, 0x20, 0x61, 0x20, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74,
	0x12, 0x9a, 0x01, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x1b,
	0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x42,
	0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6d, 0x61,
	0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63,
	0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x52, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x10, 0x3a, 0x01, 0x2a, 0x22, 0x0b, 0x2f, 0x73, 0x65, 0x6e, 0x64, 0x2f, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x92, 0x41, 0x39, 0x12, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x61, 0x74, 0x63, 0x68, 0x1a,
	0x2c, 0x53, 0x65, 0x6e, 0x64, 0x20, 0x61, 0x20, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x20, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x20, 0x74, 0x6f, 0x20, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70,
	0x6c, 0x65, 0x20, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x42, 0xfc, 0x01,
	0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42,
	0x08, 0x41, 0x70, 0x69, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x61, 0x66, 0x66, 0x6c, 0x65, 0x48, 0x61,
	0x63, 0x6b, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x3b, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4d, 0x58, 0x58,
	0xaa, 0x02, 0x09, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09, 0x4d,
	0x61, 0x69, 0x6c, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x4d, 0x61, 0x69, 0x6c, 0x65,
	0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x0a, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x92, 0x41, 0x5c,
	0x12, 0x5a, 0x0a, 0x0e, 0x4d, 0x61, 0x69, 0x6c, 0x65, 0x72, 0x20, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2a, 0x43, 0x0a, 0x03, 0x4d, 0x49, 0x54, 0x12, 0x3c, 0x68, 0x74, 0x74, 0x70, 0x73,
	0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x57, 0x61,
	0x66, 0x66, 0x6c, 0x65, 0x48, 0x61, 0x63, 0x6b, 0x73, 0x2f, 0x6d, 0x61, 0x69, 0x6c, 0x65, 0x72,
	0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2f, 0x4c, 0x49, 0x43,
	0x45, 0x4e, 0x53, 0x45, 0x2e, 0x6d, 0x64, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_mailer_v1_api_proto_rawDescOnce sync.Once
	file_mailer_v1_api_proto_rawDescData = file_mailer_v1_api_proto_rawDesc
)

func file_mailer_v1_api_proto_rawDescGZIP() []byte {
	file_mailer_v1_api_proto_rawDescOnce.Do(func() {
		file_mailer_v1_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_mailer_v1_api_proto_rawDescData)
	})
	return file_mailer_v1_api_proto_rawDescData
}

var file_mailer_v1_api_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_mailer_v1_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mailer_v1_api_proto_goTypes = []interface{}{
	(BodyType)(0),             // 0: mailer.v1.BodyType
	(*SendRequest)(nil),       // 1: mailer.v1.SendRequest
	(*SendResponse)(nil),      // 2: mailer.v1.SendResponse
	(*SendBatchRequest)(nil),  // 3: mailer.v1.SendBatchRequest
	(*SendBatchResponse)(nil), // 4: mailer.v1.SendBatchResponse
}
var file_mailer_v1_api_proto_depIdxs = []int32{
	0, // 0: mailer.v1.SendRequest.type:type_name -> mailer.v1.BodyType
	0, // 1: mailer.v1.SendBatchRequest.type:type_name -> mailer.v1.BodyType
	1, // 2: mailer.v1.MailerService.Send:input_type -> mailer.v1.SendRequest
	3, // 3: mailer.v1.MailerService.SendBatch:input_type -> mailer.v1.SendBatchRequest
	2, // 4: mailer.v1.MailerService.Send:output_type -> mailer.v1.SendResponse
	4, // 5: mailer.v1.MailerService.SendBatch:output_type -> mailer.v1.SendBatchResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_mailer_v1_api_proto_init() }
func file_mailer_v1_api_proto_init() {
	if File_mailer_v1_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mailer_v1_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendRequest); i {
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
		file_mailer_v1_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendResponse); i {
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
		file_mailer_v1_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendBatchRequest); i {
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
		file_mailer_v1_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SendBatchResponse); i {
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
			RawDescriptor: file_mailer_v1_api_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mailer_v1_api_proto_goTypes,
		DependencyIndexes: file_mailer_v1_api_proto_depIdxs,
		EnumInfos:         file_mailer_v1_api_proto_enumTypes,
		MessageInfos:      file_mailer_v1_api_proto_msgTypes,
	}.Build()
	File_mailer_v1_api_proto = out.File
	file_mailer_v1_api_proto_rawDesc = nil
	file_mailer_v1_api_proto_goTypes = nil
	file_mailer_v1_api_proto_depIdxs = nil
}
