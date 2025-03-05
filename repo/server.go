package main

import (
	"log"
	"net"
	"time"

	pb "example/proto/pingpb"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type server struct {
	pb.UnimplementedPingServiceServer
}

func (s *server) PingStream(req *pb.PingRequest, stream pb.PingService_PingStreamServer) error {
	log.Println("Received Ping request:", req.Message)

	for i := 0; i < 9999999; i++ { // Send 5 pings
		if err := stream.Send(&pb.PingResponse{Message: "Pong"}); err != nil {
			return err
		}
		time.Sleep(5 * time.Second) // 5-second interval
	}

	return nil
}

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    15 * time.Second,
			Timeout: 15 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             15 * time.Second,
			PermitWithoutStream: true,
		}),
	}
	grpcServer := grpc.NewServer(grpcOptions...)
	pb.RegisterPingServiceServer(grpcServer, &server{})

	log.Printf("gRPC Server listening on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the gRPC server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		startServer(port)
	},
}

func init() {
	ServerCmd.Flags().StringP("port", "p", "50051", "Port for the gRPC server")
}
