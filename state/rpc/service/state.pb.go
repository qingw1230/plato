// Code generated by protoc-gen-go. DO NOT EDIT.
// source: state.proto

package service

import (
	context "context"
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type StateRequest struct {
	Endpoint             string   `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Fd                   int32    `protobuf:"varint,2,opt,name=fd,proto3" json:"fd,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateRequest) Reset()         { *m = StateRequest{} }
func (m *StateRequest) String() string { return proto.CompactTextString(m) }
func (*StateRequest) ProtoMessage()    {}
func (*StateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a888679467bb7853, []int{0}
}

func (m *StateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateRequest.Unmarshal(m, b)
}
func (m *StateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateRequest.Marshal(b, m, deterministic)
}
func (m *StateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateRequest.Merge(m, src)
}
func (m *StateRequest) XXX_Size() int {
	return xxx_messageInfo_StateRequest.Size(m)
}
func (m *StateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StateRequest proto.InternalMessageInfo

func (m *StateRequest) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

func (m *StateRequest) GetFd() int32 {
	if m != nil {
		return m.Fd
	}
	return 0
}

func (m *StateRequest) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type StateResponse struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StateResponse) Reset()         { *m = StateResponse{} }
func (m *StateResponse) String() string { return proto.CompactTextString(m) }
func (*StateResponse) ProtoMessage()    {}
func (*StateResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_a888679467bb7853, []int{1}
}

func (m *StateResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StateResponse.Unmarshal(m, b)
}
func (m *StateResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StateResponse.Marshal(b, m, deterministic)
}
func (m *StateResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StateResponse.Merge(m, src)
}
func (m *StateResponse) XXX_Size() int {
	return xxx_messageInfo_StateResponse.Size(m)
}
func (m *StateResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StateResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StateResponse proto.InternalMessageInfo

func (m *StateResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *StateResponse) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*StateRequest)(nil), "service.StateRequest")
	proto.RegisterType((*StateResponse)(nil), "service.StateResponse")
}

func init() { proto.RegisterFile("state.proto", fileDescriptor_a888679467bb7853) }

var fileDescriptor_a888679467bb7853 = []byte{
	// 212 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2e, 0x2e, 0x49, 0x2c,
	0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e,
	0x55, 0xf2, 0xe3, 0xe2, 0x09, 0x06, 0x89, 0x07, 0xa5, 0x16, 0x96, 0xa6, 0x16, 0x97, 0x08, 0x49,
	0x71, 0x71, 0xa4, 0xe6, 0xa5, 0x14, 0xe4, 0x67, 0xe6, 0x95, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0xc1, 0xf9, 0x42, 0x7c, 0x5c, 0x4c, 0x69, 0x29, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xac, 0x41,
	0x4c, 0x69, 0x29, 0x42, 0x42, 0x5c, 0x2c, 0x29, 0x89, 0x25, 0x89, 0x12, 0xcc, 0x0a, 0x8c, 0x1a,
	0x3c, 0x41, 0x60, 0xb6, 0x92, 0x29, 0x17, 0x2f, 0xd4, 0xbc, 0xe2, 0x82, 0xfc, 0xbc, 0xe2, 0x54,
	0x90, 0xa2, 0xe4, 0xfc, 0x94, 0x54, 0xb0, 0x61, 0xac, 0x41, 0x60, 0xb6, 0x90, 0x00, 0x17, 0x73,
	0x6e, 0x71, 0x3a, 0xd8, 0x24, 0xce, 0x20, 0x10, 0xd3, 0xa8, 0x8e, 0x8b, 0x15, 0xec, 0x3c, 0x21,
	0x6b, 0x2e, 0x2e, 0xe7, 0xc4, 0xbc, 0xe4, 0xd4, 0x1c, 0xe7, 0xfc, 0xbc, 0x3c, 0x21, 0x51, 0x3d,
	0xa8, 0x3b, 0xf5, 0x90, 0x1d, 0x29, 0x25, 0x86, 0x2e, 0x0c, 0xb5, 0xcb, 0x82, 0x8b, 0x3d, 0x38,
	0x35, 0x2f, 0xc5, 0xb7, 0x38, 0x9d, 0x44, 0x9d, 0x4e, 0x3c, 0x51, 0x5c, 0x7a, 0xfa, 0xd6, 0x50,
	0xb9, 0x24, 0x36, 0x70, 0x20, 0x19, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x99, 0x42, 0x68, 0x74,
	0x33, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// StateClient is the client API for State service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type StateClient interface {
	CancelConn(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error)
	SendMsg(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error)
}

type stateClient struct {
	cc *grpc.ClientConn
}

func NewStateClient(cc *grpc.ClientConn) StateClient {
	return &stateClient{cc}
}

func (c *stateClient) CancelConn(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error) {
	out := new(StateResponse)
	err := c.cc.Invoke(ctx, "/service.state/CancelConn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateClient) SendMsg(ctx context.Context, in *StateRequest, opts ...grpc.CallOption) (*StateResponse, error) {
	out := new(StateResponse)
	err := c.cc.Invoke(ctx, "/service.state/SendMsg", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateServer is the server API for State service.
type StateServer interface {
	CancelConn(context.Context, *StateRequest) (*StateResponse, error)
	SendMsg(context.Context, *StateRequest) (*StateResponse, error)
}

// UnimplementedStateServer can be embedded to have forward compatible implementations.
type UnimplementedStateServer struct {
}

func (*UnimplementedStateServer) CancelConn(ctx context.Context, req *StateRequest) (*StateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelConn not implemented")
}
func (*UnimplementedStateServer) SendMsg(ctx context.Context, req *StateRequest) (*StateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}

func RegisterStateServer(s *grpc.Server, srv StateServer) {
	s.RegisterService(&_State_serviceDesc, srv)
}

func _State_CancelConn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServer).CancelConn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.state/CancelConn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServer).CancelConn(ctx, req.(*StateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _State_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.state/SendMsg",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateServer).SendMsg(ctx, req.(*StateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _State_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.state",
	HandlerType: (*StateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CancelConn",
			Handler:    _State_CancelConn_Handler,
		},
		{
			MethodName: "SendMsg",
			Handler:    _State_SendMsg_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "state.proto",
}
