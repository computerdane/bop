// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: bop.proto

package bop

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Bop_List_FullMethodName = "/bop.Bop/List"
)

// BopClient is the client API for Bop service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BopClient interface {
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error)
}

type bopClient struct {
	cc grpc.ClientConnInterface
}

func NewBopClient(cc grpc.ClientConnInterface) BopClient {
	return &bopClient{cc}
}

func (c *bopClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListReply, error) {
	out := new(ListReply)
	err := c.cc.Invoke(ctx, Bop_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BopServer is the server API for Bop service.
// All implementations must embed UnimplementedBopServer
// for forward compatibility
type BopServer interface {
	List(context.Context, *ListRequest) (*ListReply, error)
	mustEmbedUnimplementedBopServer()
}

// UnimplementedBopServer must be embedded to have forward compatible implementations.
type UnimplementedBopServer struct {
}

func (UnimplementedBopServer) List(context.Context, *ListRequest) (*ListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedBopServer) mustEmbedUnimplementedBopServer() {}

// UnsafeBopServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BopServer will
// result in compilation errors.
type UnsafeBopServer interface {
	mustEmbedUnimplementedBopServer()
}

func RegisterBopServer(s grpc.ServiceRegistrar, srv BopServer) {
	s.RegisterService(&Bop_ServiceDesc, srv)
}

func _Bop_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BopServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bop_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BopServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bop_ServiceDesc is the grpc.ServiceDesc for Bop service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bop_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bop.Bop",
	HandlerType: (*BopServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Bop_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bop.proto",
}
