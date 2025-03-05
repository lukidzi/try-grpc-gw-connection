package main

import (
	"io"
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

func (s *server) PingStream(stream pb.PingService_PingStreamServer) error {
	log.Println("Client connected to PingStream")

	for {
		req, err := stream.Recv() // Read from the client stream
		if err == io.EOF {
			log.Println("Client disconnected")
			return nil
		}
		if err != nil {
			log.Printf("Error receiving from client: %v", err)
			return err
		}

		log.Println("Received from client:", req.Message)

		// Send a "Pong" response
		if err := stream.Send(&pb.PingResponse{Message: "Pong"}); err != nil {
			log.Printf("Error sending response: %v", err)
			return err
		}

		time.Sleep(2 * time.Second) // Simulate processing time
	}
}

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000000),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second, // Ping clients every 10s
			Timeout: 5 * time.Second,  // Wait 5s before closing inactive clients
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             10 * time.Second,
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
