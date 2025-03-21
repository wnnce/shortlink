// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v6.30.0
// source: protos/shortlink.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type QueryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BaseKey       string                 `protobuf:"bytes,1,opt,name=baseKey,proto3" json:"baseKey,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryRequest) Reset() {
	*x = QueryRequest{}
	mi := &file_protos_shortlink_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryRequest) ProtoMessage() {}

func (x *QueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryRequest.ProtoReflect.Descriptor instead.
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{0}
}

func (x *QueryRequest) GetBaseKey() string {
	if x != nil {
		return x.BaseKey
	}
	return ""
}

type QueryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginUrl     string                 `protobuf:"bytes,1,opt,name=originUrl,proto3" json:"originUrl,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueryResponse) Reset() {
	*x = QueryResponse{}
	mi := &file_protos_shortlink_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryResponse) ProtoMessage() {}

func (x *QueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryResponse.ProtoReflect.Descriptor instead.
func (*QueryResponse) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{1}
}

func (x *QueryResponse) GetOriginUrl() string {
	if x != nil {
		return x.OriginUrl
	}
	return ""
}

type CreateRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OriginUrl     string                 `protobuf:"bytes,1,opt,name=originUrl,proto3" json:"originUrl,omitempty"`    // 原始url
	ClientIp      string                 `protobuf:"bytes,2,opt,name=clientIp,proto3" json:"clientIp,omitempty"`      // 客户端ip
	UserAgent     string                 `protobuf:"bytes,3,opt,name=userAgent,proto3" json:"userAgent,omitempty"`    // 客户端userAgent
	IsLasting     bool                   `protobuf:"varint,4,opt,name=isLasting,proto3" json:"isLasting,omitempty"`   // 是否长期有效
	ValidHour     int32                  `protobuf:"varint,5,opt,name=validHour,proto3" json:"validHour,omitempty"`   // 链接的有效时间 单位：小时
	ExpireMode    int32                  `protobuf:"varint,6,opt,name=expireMode,proto3" json:"expireMode,omitempty"` // 链接的过期删除模式，1:精确模式 2:模糊模式
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	mi := &file_protos_shortlink_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{2}
}

func (x *CreateRequest) GetOriginUrl() string {
	if x != nil {
		return x.OriginUrl
	}
	return ""
}

func (x *CreateRequest) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *CreateRequest) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

func (x *CreateRequest) GetIsLasting() bool {
	if x != nil {
		return x.IsLasting
	}
	return false
}

func (x *CreateRequest) GetValidHour() int32 {
	if x != nil {
		return x.ValidHour
	}
	return 0
}

func (x *CreateRequest) GetExpireMode() int32 {
	if x != nil {
		return x.ExpireMode
	}
	return 0
}

type CreateResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BaseKey       string                 `protobuf:"bytes,1,opt,name=baseKey,proto3" json:"baseKey,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	mi := &file_protos_shortlink_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{3}
}

func (x *CreateResponse) GetBaseKey() string {
	if x != nil {
		return x.BaseKey
	}
	return ""
}

type DeleteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	BaseKey       string                 `protobuf:"bytes,1,opt,name=baseKey,proto3" json:"baseKey,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	mi := &file_protos_shortlink_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRequest) GetBaseKey() string {
	if x != nil {
		return x.BaseKey
	}
	return ""
}

type LinkRecord struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty" db:"id"'`                 // 数据库自增id
	UniqueId      string                 `protobuf:"bytes,2,opt,name=UniqueId,proto3" json:"uniqueId,omitempty" db:"unique_id"`      // 唯一Id，由服务端生成
	BaseValue     string                 `protobuf:"bytes,3,opt,name=BaseValue,proto3" json:"baseValue,omitempty" db:"base_value"`    // 唯一Id的base62值
	OriginUrl     string                 `protobuf:"bytes,4,opt,name=OriginUrl,proto3" json:"originUrl,omitempty" db:"origin_url"`    // 对应的原链接Url
	ValidHour     int32                  `protobuf:"varint,5,opt,name=ValidHour,proto3" json:"validHour,omitempty" db:"valid_hour"`   // 链接的有效时间，长期链接为空
	IsLasting     bool                   `protobuf:"varint,6,opt,name=IsLasting,proto3" json:"isLasting,omitempty" db:"is_lasting"`   // 是否为长期有效的链接
	CreateTime    string                 `protobuf:"bytes,7,opt,name=CreateTime,proto3" json:"createTime,omitempty" db:"create_time"`  // 创建时间
	ExpireTime    int64                  `protobuf:"varint,8,opt,name=ExpireTime,proto3" json:"expireTime,omitempty" db:"expire_time"` // 过期时间
	ExpireMode    int32                  `protobuf:"varint,9,opt,name=ExpireMode,proto3" json:"expireMode" db:"expire_mode"` // 连接的过期删除策略 0:模糊 1:精确
	ClientIp      string                 `protobuf:"bytes,10,opt,name=ClientIp,proto3" json:"clientIp,omitempty" db:"client_ip"`     // 客户端Ip
	UserAgent     string                 `protobuf:"bytes,11,opt,name=UserAgent,proto3" json:"userAgent" db:"user_agent"`   // 客户端的UserAgent
	Status        int32                  `protobuf:"varint,12,opt,name=Status,proto3" json:"status" db:"status"`        // 此链接的状态 0:正常 1:禁用
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LinkRecord) Reset() {
	*x = LinkRecord{}
	mi := &file_protos_shortlink_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LinkRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LinkRecord) ProtoMessage() {}

func (x *LinkRecord) ProtoReflect() protoreflect.Message {
	mi := &file_protos_shortlink_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LinkRecord.ProtoReflect.Descriptor instead.
func (*LinkRecord) Descriptor() ([]byte, []int) {
	return file_protos_shortlink_proto_rawDescGZIP(), []int{5}
}

func (x *LinkRecord) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LinkRecord) GetUniqueId() string {
	if x != nil {
		return x.UniqueId
	}
	return ""
}

func (x *LinkRecord) GetBaseValue() string {
	if x != nil {
		return x.BaseValue
	}
	return ""
}

func (x *LinkRecord) GetOriginUrl() string {
	if x != nil {
		return x.OriginUrl
	}
	return ""
}

func (x *LinkRecord) GetValidHour() int32 {
	if x != nil {
		return x.ValidHour
	}
	return 0
}

func (x *LinkRecord) GetIsLasting() bool {
	if x != nil {
		return x.IsLasting
	}
	return false
}

func (x *LinkRecord) GetCreateTime() string {
	if x != nil {
		return x.CreateTime
	}
	return ""
}

func (x *LinkRecord) GetExpireTime() int64 {
	if x != nil {
		return x.ExpireTime
	}
	return 0
}

func (x *LinkRecord) GetExpireMode() int32 {
	if x != nil {
		return x.ExpireMode
	}
	return 0
}

func (x *LinkRecord) GetClientIp() string {
	if x != nil {
		return x.ClientIp
	}
	return ""
}

func (x *LinkRecord) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

func (x *LinkRecord) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_protos_shortlink_proto protoreflect.FileDescriptor

var file_protos_shortlink_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69,
	0x6e, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c,
	0x69, 0x6e, 0x6b, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x28, 0x0a, 0x0c, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x62, 0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x22, 0x2d, 0x0a, 0x0d, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x22, 0xc3, 0x01, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x69, 0x73, 0x4c, 0x61, 0x73, 0x74, 0x69, 0x6e, 0x67,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x69, 0x73, 0x4c, 0x61, 0x73, 0x74, 0x69, 0x6e,
	0x67, 0x12, 0x1c, 0x0a, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x12,
	0x1e, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x22,
	0x2a, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x62, 0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x22, 0x29, 0x0a, 0x0d, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x62, 0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x62,
	0x61, 0x73, 0x65, 0x4b, 0x65, 0x79, 0x22, 0xe2, 0x02, 0x0a, 0x0a, 0x4c, 0x69, 0x6e, 0x6b, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x55, 0x6e, 0x69, 0x71, 0x75, 0x65, 0x49,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x42, 0x61, 0x73, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x42, 0x61, 0x73, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x1c, 0x0a,
	0x09, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x48, 0x6f, 0x75, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x49,
	0x73, 0x4c, 0x61, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09,
	0x49, 0x73, 0x4c, 0x61, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x45, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x45, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x70, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x49, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x55, 0x73, 0x65, 0x72, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0xa5, 0x02, 0x0a, 0x0b,
	0x4c, 0x69, 0x6e, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0e, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x17, 0x2e,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69,
	0x6e, 0x6b, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x46, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4c,
	0x69, 0x6e, 0x6b, 0x12, 0x18, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0f, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x12, 0x18, 0x2e, 0x73, 0x68,
	0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x44, 0x0a,
	0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x17, 0x2e, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e,
	0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73,
	0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e, 0x6b, 0x2e, 0x4c, 0x69, 0x6e, 0x6b, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x42, 0x0f, 0x5a, 0x0d, 0x2f, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x6c, 0x69, 0x6e,
	0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_protos_shortlink_proto_rawDescOnce sync.Once
	file_protos_shortlink_proto_rawDescData []byte
)

func file_protos_shortlink_proto_rawDescGZIP() []byte {
	file_protos_shortlink_proto_rawDescOnce.Do(func() {
		file_protos_shortlink_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_protos_shortlink_proto_rawDesc), len(file_protos_shortlink_proto_rawDesc)))
	})
	return file_protos_shortlink_proto_rawDescData
}

var file_protos_shortlink_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_protos_shortlink_proto_goTypes = []any{
	(*QueryRequest)(nil),   // 0: shortlink.QueryRequest
	(*QueryResponse)(nil),  // 1: shortlink.QueryResponse
	(*CreateRequest)(nil),  // 2: shortlink.CreateRequest
	(*CreateResponse)(nil), // 3: shortlink.CreateResponse
	(*DeleteRequest)(nil),  // 4: shortlink.DeleteRequest
	(*LinkRecord)(nil),     // 5: shortlink.LinkRecord
	(*emptypb.Empty)(nil),  // 6: google.protobuf.Empty
}
var file_protos_shortlink_proto_depIdxs = []int32{
	0, // 0: shortlink.LinkService.QueryOriginUrl:input_type -> shortlink.QueryRequest
	2, // 1: shortlink.LinkService.CreateShortLink:input_type -> shortlink.CreateRequest
	4, // 2: shortlink.LinkService.DeleteShortLink:input_type -> shortlink.DeleteRequest
	0, // 3: shortlink.LinkService.QueryShortLinkInfo:input_type -> shortlink.QueryRequest
	1, // 4: shortlink.LinkService.QueryOriginUrl:output_type -> shortlink.QueryResponse
	3, // 5: shortlink.LinkService.CreateShortLink:output_type -> shortlink.CreateResponse
	6, // 6: shortlink.LinkService.DeleteShortLink:output_type -> google.protobuf.Empty
	5, // 7: shortlink.LinkService.QueryShortLinkInfo:output_type -> shortlink.LinkRecord
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_shortlink_proto_init() }
func file_protos_shortlink_proto_init() {
	if File_protos_shortlink_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_protos_shortlink_proto_rawDesc), len(file_protos_shortlink_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_shortlink_proto_goTypes,
		DependencyIndexes: file_protos_shortlink_proto_depIdxs,
		MessageInfos:      file_protos_shortlink_proto_msgTypes,
	}.Build()
	File_protos_shortlink_proto = out.File
	file_protos_shortlink_proto_goTypes = nil
	file_protos_shortlink_proto_depIdxs = nil
}
