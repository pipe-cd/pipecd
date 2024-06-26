// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/plugin/api/v1alpha1/platform/api.proto

package platform

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

// PlannerServiceClient is the client API for PlannerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlannerServiceClient interface {
	// DetermineStrategy determines which strategy should be used for the given deployment.
	DetermineStrategy(ctx context.Context, in *DetermineStrategyRequest, opts ...grpc.CallOption) (*DetermineStrategyResponse, error)
	// QuickSyncPlan builds plan for the given deployment using quick sync strategy.
	QuickSyncPlan(ctx context.Context, in *QuickSyncPlanRequest, opts ...grpc.CallOption) (*QuickSyncPlanResponse, error)
	// PipelineSyncPlan builds plan for the given deployment using pipeline sync strategy.
	PipelineSyncPlan(ctx context.Context, in *PipelineSyncPlanRequest, opts ...grpc.CallOption) (*PipelineSyncPlanResponse, error)
}

type plannerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlannerServiceClient(cc grpc.ClientConnInterface) PlannerServiceClient {
	return &plannerServiceClient{cc}
}

func (c *plannerServiceClient) DetermineStrategy(ctx context.Context, in *DetermineStrategyRequest, opts ...grpc.CallOption) (*DetermineStrategyResponse, error) {
	out := new(DetermineStrategyResponse)
	err := c.cc.Invoke(ctx, "/grpc.plugin.platformapi.v1alpha1.PlannerService/DetermineStrategy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plannerServiceClient) QuickSyncPlan(ctx context.Context, in *QuickSyncPlanRequest, opts ...grpc.CallOption) (*QuickSyncPlanResponse, error) {
	out := new(QuickSyncPlanResponse)
	err := c.cc.Invoke(ctx, "/grpc.plugin.platformapi.v1alpha1.PlannerService/QuickSyncPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *plannerServiceClient) PipelineSyncPlan(ctx context.Context, in *PipelineSyncPlanRequest, opts ...grpc.CallOption) (*PipelineSyncPlanResponse, error) {
	out := new(PipelineSyncPlanResponse)
	err := c.cc.Invoke(ctx, "/grpc.plugin.platformapi.v1alpha1.PlannerService/PipelineSyncPlan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlannerServiceServer is the server API for PlannerService service.
// All implementations must embed UnimplementedPlannerServiceServer
// for forward compatibility
type PlannerServiceServer interface {
	// DetermineStrategy determines which strategy should be used for the given deployment.
	DetermineStrategy(context.Context, *DetermineStrategyRequest) (*DetermineStrategyResponse, error)
	// QuickSyncPlan builds plan for the given deployment using quick sync strategy.
	QuickSyncPlan(context.Context, *QuickSyncPlanRequest) (*QuickSyncPlanResponse, error)
	// PipelineSyncPlan builds plan for the given deployment using pipeline sync strategy.
	PipelineSyncPlan(context.Context, *PipelineSyncPlanRequest) (*PipelineSyncPlanResponse, error)
	mustEmbedUnimplementedPlannerServiceServer()
}

// UnimplementedPlannerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPlannerServiceServer struct {
}

func (UnimplementedPlannerServiceServer) DetermineStrategy(context.Context, *DetermineStrategyRequest) (*DetermineStrategyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetermineStrategy not implemented")
}
func (UnimplementedPlannerServiceServer) QuickSyncPlan(context.Context, *QuickSyncPlanRequest) (*QuickSyncPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuickSyncPlan not implemented")
}
func (UnimplementedPlannerServiceServer) PipelineSyncPlan(context.Context, *PipelineSyncPlanRequest) (*PipelineSyncPlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PipelineSyncPlan not implemented")
}
func (UnimplementedPlannerServiceServer) mustEmbedUnimplementedPlannerServiceServer() {}

// UnsafePlannerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlannerServiceServer will
// result in compilation errors.
type UnsafePlannerServiceServer interface {
	mustEmbedUnimplementedPlannerServiceServer()
}

func RegisterPlannerServiceServer(s grpc.ServiceRegistrar, srv PlannerServiceServer) {
	s.RegisterService(&PlannerService_ServiceDesc, srv)
}

func _PlannerService_DetermineStrategy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DetermineStrategyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlannerServiceServer).DetermineStrategy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.plugin.platformapi.v1alpha1.PlannerService/DetermineStrategy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlannerServiceServer).DetermineStrategy(ctx, req.(*DetermineStrategyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlannerService_QuickSyncPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuickSyncPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlannerServiceServer).QuickSyncPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.plugin.platformapi.v1alpha1.PlannerService/QuickSyncPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlannerServiceServer).QuickSyncPlan(ctx, req.(*QuickSyncPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlannerService_PipelineSyncPlan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PipelineSyncPlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlannerServiceServer).PipelineSyncPlan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.plugin.platformapi.v1alpha1.PlannerService/PipelineSyncPlan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlannerServiceServer).PipelineSyncPlan(ctx, req.(*PipelineSyncPlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlannerService_ServiceDesc is the grpc.ServiceDesc for PlannerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlannerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.plugin.platformapi.v1alpha1.PlannerService",
	HandlerType: (*PlannerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DetermineStrategy",
			Handler:    _PlannerService_DetermineStrategy_Handler,
		},
		{
			MethodName: "QuickSyncPlan",
			Handler:    _PlannerService_QuickSyncPlan_Handler,
		},
		{
			MethodName: "PipelineSyncPlan",
			Handler:    _PlannerService_PipelineSyncPlan_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/plugin/api/v1alpha1/platform/api.proto",
}

// ExecutorServiceClient is the client API for ExecutorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExecutorServiceClient interface {
	// Execute executes the given stage of the deployment plan.
	ExecuteStage(ctx context.Context, in *ExecuteStageRequest, opts ...grpc.CallOption) (ExecutorService_ExecuteStageClient, error)
}

type executorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExecutorServiceClient(cc grpc.ClientConnInterface) ExecutorServiceClient {
	return &executorServiceClient{cc}
}

func (c *executorServiceClient) ExecuteStage(ctx context.Context, in *ExecuteStageRequest, opts ...grpc.CallOption) (ExecutorService_ExecuteStageClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExecutorService_ServiceDesc.Streams[0], "/grpc.plugin.platformapi.v1alpha1.ExecutorService/ExecuteStage", opts...)
	if err != nil {
		return nil, err
	}
	x := &executorServiceExecuteStageClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExecutorService_ExecuteStageClient interface {
	Recv() (*ExecuteStageResponse, error)
	grpc.ClientStream
}

type executorServiceExecuteStageClient struct {
	grpc.ClientStream
}

func (x *executorServiceExecuteStageClient) Recv() (*ExecuteStageResponse, error) {
	m := new(ExecuteStageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExecutorServiceServer is the server API for ExecutorService service.
// All implementations must embed UnimplementedExecutorServiceServer
// for forward compatibility
type ExecutorServiceServer interface {
	// Execute executes the given stage of the deployment plan.
	ExecuteStage(*ExecuteStageRequest, ExecutorService_ExecuteStageServer) error
	mustEmbedUnimplementedExecutorServiceServer()
}

// UnimplementedExecutorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExecutorServiceServer struct {
}

func (UnimplementedExecutorServiceServer) ExecuteStage(*ExecuteStageRequest, ExecutorService_ExecuteStageServer) error {
	return status.Errorf(codes.Unimplemented, "method ExecuteStage not implemented")
}
func (UnimplementedExecutorServiceServer) mustEmbedUnimplementedExecutorServiceServer() {}

// UnsafeExecutorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExecutorServiceServer will
// result in compilation errors.
type UnsafeExecutorServiceServer interface {
	mustEmbedUnimplementedExecutorServiceServer()
}

func RegisterExecutorServiceServer(s grpc.ServiceRegistrar, srv ExecutorServiceServer) {
	s.RegisterService(&ExecutorService_ServiceDesc, srv)
}

func _ExecutorService_ExecuteStage_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExecuteStageRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExecutorServiceServer).ExecuteStage(m, &executorServiceExecuteStageServer{stream})
}

type ExecutorService_ExecuteStageServer interface {
	Send(*ExecuteStageResponse) error
	grpc.ServerStream
}

type executorServiceExecuteStageServer struct {
	grpc.ServerStream
}

func (x *executorServiceExecuteStageServer) Send(m *ExecuteStageResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ExecutorService_ServiceDesc is the grpc.ServiceDesc for ExecutorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExecutorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.plugin.platformapi.v1alpha1.ExecutorService",
	HandlerType: (*ExecutorServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecuteStage",
			Handler:       _ExecutorService_ExecuteStage_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/plugin/api/v1alpha1/platform/api.proto",
}
