// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pingpb

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

// PingServiceClient is the client API for PingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PingServiceClient interface {
	PingStream(ctx context.Context, opts ...grpc.CallOption) (PingService_PingStreamClient, error)
}

type pingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPingServiceClient(cc grpc.ClientConnInterface) PingServiceClient {
	return &pingServiceClient{cc}
}

func (c *pingServiceClient) PingStream(ctx context.Context, opts ...grpc.CallOption) (PingService_PingStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &PingService_ServiceDesc.Streams[0], "/ping.PingService/PingStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &pingServicePingStreamClient{stream}
	return x, nil
}

type PingService_PingStreamClient interface {
	Send(*PingRequest) error
	Recv() (*PingResponse, error)
	grpc.ClientStream
}

type pingServicePingStreamClient struct {
	grpc.ClientStream
}

func (x *pingServicePingStreamClient) Send(m *PingRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *pingServicePingStreamClient) Recv() (*PingResponse, error) {
	m := new(PingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PingServiceServer is the server API for PingService service.
// All implementations must embed UnimplementedPingServiceServer
// for forward compatibility
type PingServiceServer interface {
	PingStream(PingService_PingStreamServer) error
	mustEmbedUnimplementedPingServiceServer()
}

// UnimplementedPingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPingServiceServer struct {
}

func (UnimplementedPingServiceServer) PingStream(PingService_PingStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method PingStream not implemented")
}
func (UnimplementedPingServiceServer) mustEmbedUnimplementedPingServiceServer() {}

// UnsafePingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PingServiceServer will
// result in compilation errors.
type UnsafePingServiceServer interface {
	mustEmbedUnimplementedPingServiceServer()
}

func RegisterPingServiceServer(s grpc.ServiceRegistrar, srv PingServiceServer) {
	s.RegisterService(&PingService_ServiceDesc, srv)
}

func _PingService_PingStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PingServiceServer).PingStream(&pingServicePingStreamServer{stream})
}

type PingService_PingStreamServer interface {
	Send(*PingResponse) error
	Recv() (*PingRequest, error)
	grpc.ServerStream
}

type pingServicePingStreamServer struct {
	grpc.ServerStream
}

func (x *pingServicePingStreamServer) Send(m *PingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *pingServicePingStreamServer) Recv() (*PingRequest, error) {
	m := new(PingRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PingService_ServiceDesc is the grpc.ServiceDesc for PingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ping.PingService",
	HandlerType: (*PingServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PingStream",
			Handler:       _PingService_PingStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "ping.proto",
}
