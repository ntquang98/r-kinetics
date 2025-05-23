// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: app/v1/analytics_job.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	AnalyticsJob_CreateAnalyticsJob_FullMethodName   = "/api.app.v1.AnalyticsJob/CreateAnalyticsJob"
	AnalyticsJob_GetAnalyticsJob_FullMethodName      = "/api.app.v1.AnalyticsJob/GetAnalyticsJob"
	AnalyticsJob_ListAnalyticsJob_FullMethodName     = "/api.app.v1.AnalyticsJob/ListAnalyticsJob"
	AnalyticsJob_CompleteAnalyticsJob_FullMethodName = "/api.app.v1.AnalyticsJob/CompleteAnalyticsJob"
	AnalyticsJob_RePushJob_FullMethodName            = "/api.app.v1.AnalyticsJob/RePushJob"
)

// AnalyticsJobClient is the client API for AnalyticsJob service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyticsJobClient interface {
	CreateAnalyticsJob(ctx context.Context, in *CreateAnalyticsJobRequest, opts ...grpc.CallOption) (*CreateAnalyticsJobReply, error)
	GetAnalyticsJob(ctx context.Context, in *GetAnalyticsJobRequest, opts ...grpc.CallOption) (*GetAnalyticsJobReply, error)
	ListAnalyticsJob(ctx context.Context, in *ListAnalyticsJobRequest, opts ...grpc.CallOption) (*ListAnalyticsJobReply, error)
	CompleteAnalyticsJob(ctx context.Context, in *CompleteAnalyticsJobRequest, opts ...grpc.CallOption) (*CompleteAnalyticsJobReply, error)
	RePushJob(ctx context.Context, in *RePushJobRequest, opts ...grpc.CallOption) (*RePushJobReply, error)
}

type analyticsJobClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyticsJobClient(cc grpc.ClientConnInterface) AnalyticsJobClient {
	return &analyticsJobClient{cc}
}

func (c *analyticsJobClient) CreateAnalyticsJob(ctx context.Context, in *CreateAnalyticsJobRequest, opts ...grpc.CallOption) (*CreateAnalyticsJobReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAnalyticsJobReply)
	err := c.cc.Invoke(ctx, AnalyticsJob_CreateAnalyticsJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsJobClient) GetAnalyticsJob(ctx context.Context, in *GetAnalyticsJobRequest, opts ...grpc.CallOption) (*GetAnalyticsJobReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAnalyticsJobReply)
	err := c.cc.Invoke(ctx, AnalyticsJob_GetAnalyticsJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsJobClient) ListAnalyticsJob(ctx context.Context, in *ListAnalyticsJobRequest, opts ...grpc.CallOption) (*ListAnalyticsJobReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListAnalyticsJobReply)
	err := c.cc.Invoke(ctx, AnalyticsJob_ListAnalyticsJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsJobClient) CompleteAnalyticsJob(ctx context.Context, in *CompleteAnalyticsJobRequest, opts ...grpc.CallOption) (*CompleteAnalyticsJobReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CompleteAnalyticsJobReply)
	err := c.cc.Invoke(ctx, AnalyticsJob_CompleteAnalyticsJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analyticsJobClient) RePushJob(ctx context.Context, in *RePushJobRequest, opts ...grpc.CallOption) (*RePushJobReply, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RePushJobReply)
	err := c.cc.Invoke(ctx, AnalyticsJob_RePushJob_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyticsJobServer is the server API for AnalyticsJob service.
// All implementations must embed UnimplementedAnalyticsJobServer
// for forward compatibility.
type AnalyticsJobServer interface {
	CreateAnalyticsJob(context.Context, *CreateAnalyticsJobRequest) (*CreateAnalyticsJobReply, error)
	GetAnalyticsJob(context.Context, *GetAnalyticsJobRequest) (*GetAnalyticsJobReply, error)
	ListAnalyticsJob(context.Context, *ListAnalyticsJobRequest) (*ListAnalyticsJobReply, error)
	CompleteAnalyticsJob(context.Context, *CompleteAnalyticsJobRequest) (*CompleteAnalyticsJobReply, error)
	RePushJob(context.Context, *RePushJobRequest) (*RePushJobReply, error)
	mustEmbedUnimplementedAnalyticsJobServer()
}

// UnimplementedAnalyticsJobServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAnalyticsJobServer struct{}

func (UnimplementedAnalyticsJobServer) CreateAnalyticsJob(context.Context, *CreateAnalyticsJobRequest) (*CreateAnalyticsJobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAnalyticsJob not implemented")
}
func (UnimplementedAnalyticsJobServer) GetAnalyticsJob(context.Context, *GetAnalyticsJobRequest) (*GetAnalyticsJobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAnalyticsJob not implemented")
}
func (UnimplementedAnalyticsJobServer) ListAnalyticsJob(context.Context, *ListAnalyticsJobRequest) (*ListAnalyticsJobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAnalyticsJob not implemented")
}
func (UnimplementedAnalyticsJobServer) CompleteAnalyticsJob(context.Context, *CompleteAnalyticsJobRequest) (*CompleteAnalyticsJobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteAnalyticsJob not implemented")
}
func (UnimplementedAnalyticsJobServer) RePushJob(context.Context, *RePushJobRequest) (*RePushJobReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RePushJob not implemented")
}
func (UnimplementedAnalyticsJobServer) mustEmbedUnimplementedAnalyticsJobServer() {}
func (UnimplementedAnalyticsJobServer) testEmbeddedByValue()                      {}

// UnsafeAnalyticsJobServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyticsJobServer will
// result in compilation errors.
type UnsafeAnalyticsJobServer interface {
	mustEmbedUnimplementedAnalyticsJobServer()
}

func RegisterAnalyticsJobServer(s grpc.ServiceRegistrar, srv AnalyticsJobServer) {
	// If the following call pancis, it indicates UnimplementedAnalyticsJobServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AnalyticsJob_ServiceDesc, srv)
}

func _AnalyticsJob_CreateAnalyticsJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAnalyticsJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsJobServer).CreateAnalyticsJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsJob_CreateAnalyticsJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsJobServer).CreateAnalyticsJob(ctx, req.(*CreateAnalyticsJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsJob_GetAnalyticsJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAnalyticsJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsJobServer).GetAnalyticsJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsJob_GetAnalyticsJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsJobServer).GetAnalyticsJob(ctx, req.(*GetAnalyticsJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsJob_ListAnalyticsJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAnalyticsJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsJobServer).ListAnalyticsJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsJob_ListAnalyticsJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsJobServer).ListAnalyticsJob(ctx, req.(*ListAnalyticsJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsJob_CompleteAnalyticsJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompleteAnalyticsJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsJobServer).CompleteAnalyticsJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsJob_CompleteAnalyticsJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsJobServer).CompleteAnalyticsJob(ctx, req.(*CompleteAnalyticsJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalyticsJob_RePushJob_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RePushJobRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsJobServer).RePushJob(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AnalyticsJob_RePushJob_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsJobServer).RePushJob(ctx, req.(*RePushJobRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AnalyticsJob_ServiceDesc is the grpc.ServiceDesc for AnalyticsJob service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AnalyticsJob_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.app.v1.AnalyticsJob",
	HandlerType: (*AnalyticsJobServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAnalyticsJob",
			Handler:    _AnalyticsJob_CreateAnalyticsJob_Handler,
		},
		{
			MethodName: "GetAnalyticsJob",
			Handler:    _AnalyticsJob_GetAnalyticsJob_Handler,
		},
		{
			MethodName: "ListAnalyticsJob",
			Handler:    _AnalyticsJob_ListAnalyticsJob_Handler,
		},
		{
			MethodName: "CompleteAnalyticsJob",
			Handler:    _AnalyticsJob_CompleteAnalyticsJob_Handler,
		},
		{
			MethodName: "RePushJob",
			Handler:    _AnalyticsJob_RePushJob_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app/v1/analytics_job.proto",
}
