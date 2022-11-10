// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package rpc_impl

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

// ITaskClient is the client API for ITask service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ITaskClient interface {
	List(ctx context.Context, in *Task_ListReq, opts ...grpc.CallOption) (*Task_ListRes, error)
	Create(ctx context.Context, in *Task_CreateReq, opts ...grpc.CallOption) (*Task_CreateRes, error)
	Modify(ctx context.Context, in *Task_ModifyReq, opts ...grpc.CallOption) (*Task_ModifyRes, error)
}

type iTaskClient struct {
	cc grpc.ClientConnInterface
}

func NewITaskClient(cc grpc.ClientConnInterface) ITaskClient {
	return &iTaskClient{cc}
}

func (c *iTaskClient) List(ctx context.Context, in *Task_ListReq, opts ...grpc.CallOption) (*Task_ListRes, error) {
	out := new(Task_ListRes)
	err := c.cc.Invoke(ctx, "/task.ITask/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iTaskClient) Create(ctx context.Context, in *Task_CreateReq, opts ...grpc.CallOption) (*Task_CreateRes, error) {
	out := new(Task_CreateRes)
	err := c.cc.Invoke(ctx, "/task.ITask/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *iTaskClient) Modify(ctx context.Context, in *Task_ModifyReq, opts ...grpc.CallOption) (*Task_ModifyRes, error) {
	out := new(Task_ModifyRes)
	err := c.cc.Invoke(ctx, "/task.ITask/Modify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ITaskServer is the server API for ITask service.
// All implementations must embed UnimplementedITaskServer
// for forward compatibility
type ITaskServer interface {
	List(context.Context, *Task_ListReq) (*Task_ListRes, error)
	Create(context.Context, *Task_CreateReq) (*Task_CreateRes, error)
	Modify(context.Context, *Task_ModifyReq) (*Task_ModifyRes, error)
	mustEmbedUnimplementedITaskServer()
}

// UnimplementedITaskServer must be embedded to have forward compatible implementations.
type UnimplementedITaskServer struct {
}

func (UnimplementedITaskServer) List(context.Context, *Task_ListReq) (*Task_ListRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedITaskServer) Create(context.Context, *Task_CreateReq) (*Task_CreateRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedITaskServer) Modify(context.Context, *Task_ModifyReq) (*Task_ModifyRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Modify not implemented")
}
func (UnimplementedITaskServer) mustEmbedUnimplementedITaskServer() {}

// UnsafeITaskServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ITaskServer will
// result in compilation errors.
type UnsafeITaskServer interface {
	mustEmbedUnimplementedITaskServer()
}

func RegisterITaskServer(s grpc.ServiceRegistrar, srv ITaskServer) {
	s.RegisterService(&ITask_ServiceDesc, srv)
}

func _ITask_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_ListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ITaskServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.ITask/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ITaskServer).List(ctx, req.(*Task_ListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ITask_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_CreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ITaskServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.ITask/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ITaskServer).Create(ctx, req.(*Task_CreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ITask_Modify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Task_ModifyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ITaskServer).Modify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/task.ITask/Modify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ITaskServer).Modify(ctx, req.(*Task_ModifyReq))
	}
	return interceptor(ctx, in, info, handler)
}

// ITask_ServiceDesc is the grpc.ServiceDesc for ITask service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ITask_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "task.ITask",
	HandlerType: (*ITaskServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _ITask_List_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ITask_Create_Handler,
		},
		{
			MethodName: "Modify",
			Handler:    _ITask_Modify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "4_task.proto",
}