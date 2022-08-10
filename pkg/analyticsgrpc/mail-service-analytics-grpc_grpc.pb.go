// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.13.0
// source: proto/mail-service-analytics-grpc.proto

package analyticsgrpc

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AnalyticsClient is the client API for Analytics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalyticsClient interface {
	StoreEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error)
}

type analyticsClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalyticsClient(cc grpc.ClientConnInterface) AnalyticsClient {
	return &analyticsClient{cc}
}

func (c *analyticsClient) StoreEvent(ctx context.Context, in *Event, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/analyticsgrpc.Analytics/StoreEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalyticsServer is the server API for Analytics service.
// All implementations must embed UnimplementedAnalyticsServer
// for forward compatibility
type AnalyticsServer interface {
	StoreEvent(context.Context, *Event) (*empty.Empty, error)
	mustEmbedUnimplementedAnalyticsServer()
}

// UnimplementedAnalyticsServer must be embedded to have forward compatible implementations.
type UnimplementedAnalyticsServer struct {
}

func (UnimplementedAnalyticsServer) StoreEvent(context.Context, *Event) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreEvent not implemented")
}
func (UnimplementedAnalyticsServer) mustEmbedUnimplementedAnalyticsServer() {}

// UnsafeAnalyticsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalyticsServer will
// result in compilation errors.
type UnsafeAnalyticsServer interface {
	mustEmbedUnimplementedAnalyticsServer()
}

func RegisterAnalyticsServer(s grpc.ServiceRegistrar, srv AnalyticsServer) {
	s.RegisterService(&Analytics_ServiceDesc, srv)
}

func _Analytics_StoreEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalyticsServer).StoreEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/analyticsgrpc.Analytics/StoreEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalyticsServer).StoreEvent(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// Analytics_ServiceDesc is the grpc.ServiceDesc for Analytics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Analytics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "analyticsgrpc.Analytics",
	HandlerType: (*AnalyticsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreEvent",
			Handler:    _Analytics_StoreEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/mail-service-analytics-grpc.proto",
}