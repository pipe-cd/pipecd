// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pkg/plugin/api/v1alpha1/livestate/api.proto

package livestate

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

// LivestateServiceClient is the client API for LivestateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LivestateServiceClient interface {
	// GetLivestate gets the application live state for the give app id.
	GetLivestate(ctx context.Context, in *GetLivestateRequest, opts ...grpc.CallOption) (*GetLivestateResponse, error)
}

type livestateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLivestateServiceClient(cc grpc.ClientConnInterface) LivestateServiceClient {
	return &livestateServiceClient{cc}
}

func (c *livestateServiceClient) GetLivestate(ctx context.Context, in *GetLivestateRequest, opts ...grpc.CallOption) (*GetLivestateResponse, error) {
	out := new(GetLivestateResponse)
	err := c.cc.Invoke(ctx, "/grpc.plugin.livestateapi.v1alpha1.LivestateService/GetLivestate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LivestateServiceServer is the server API for LivestateService service.
// All implementations must embed UnimplementedLivestateServiceServer
// for forward compatibility
type LivestateServiceServer interface {
	// GetLivestate gets the application live state for the give app id.
	GetLivestate(context.Context, *GetLivestateRequest) (*GetLivestateResponse, error)
	mustEmbedUnimplementedLivestateServiceServer()
}

// UnimplementedLivestateServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLivestateServiceServer struct {
}

func (UnimplementedLivestateServiceServer) GetLivestate(context.Context, *GetLivestateRequest) (*GetLivestateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLivestate not implemented")
}
func (UnimplementedLivestateServiceServer) mustEmbedUnimplementedLivestateServiceServer() {}

// UnsafeLivestateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LivestateServiceServer will
// result in compilation errors.
type UnsafeLivestateServiceServer interface {
	mustEmbedUnimplementedLivestateServiceServer()
}

func RegisterLivestateServiceServer(s grpc.ServiceRegistrar, srv LivestateServiceServer) {
	s.RegisterService(&LivestateService_ServiceDesc, srv)
}

func _LivestateService_GetLivestate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLivestateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LivestateServiceServer).GetLivestate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpc.plugin.livestateapi.v1alpha1.LivestateService/GetLivestate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LivestateServiceServer).GetLivestate(ctx, req.(*GetLivestateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LivestateService_ServiceDesc is the grpc.ServiceDesc for LivestateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LivestateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.plugin.livestateapi.v1alpha1.LivestateService",
	HandlerType: (*LivestateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLivestate",
			Handler:    _LivestateService_GetLivestate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/plugin/api/v1alpha1/livestate/api.proto",
}
