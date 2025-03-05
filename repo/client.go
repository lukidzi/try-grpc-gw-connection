package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "example/proto/pingpb"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// Connect to gRPC Server
func connectToServer(host, port string) (*grpc.ClientConn, pb.PingServiceClient, error) {
	dialOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                15 * time.Second,
			Timeout:             15 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port), dialOpts...)
	if err != nil {
		return nil, nil, err
	}

	client := pb.NewPingServiceClient(conn)
	return conn, client, nil
}

// Handles backoff and retries
func Start(port, host string) error {
	for generationID := uint64(1); ; generationID++ {
		errCh := make(chan error, 1)

		go func(errCh chan<- error) {
			defer close(errCh)

			// Recover from a panic
			defer func() {
				if e := recover(); e != nil {
					if err, ok := e.(error); ok {
						errCh <- errors.WithStack(err)
					} else {
						errCh <- errors.Errorf("%v", e)
					}
				}
			}()

			errCh <- startClient(port, host)
		}(errCh)

		select {
		case err := <-errCh:
			if err != nil {
				log.Printf("[Generation %d] Client failed with error: %v", generationID, err)
			}
		}

		// Exponential backoff to prevent aggressive reconnects
		delay := time.Duration(generationID) * time.Second
		if delay > 30*time.Second {
			delay = 30 * time.Second
		}
		log.Printf("Retrying in %s...", delay)
		<-time.After(delay)
	}
}

// Start gRPC Client
func startClient(port, host string) error {
	var conn *grpc.ClientConn
	var client pb.PingServiceClient
	var err error
	retryDelay := 1 * time.Second

	for {
		if conn == nil {
			log.Println("Attempting to connect to server...")
			conn, client, err = connectToServer(host, port)
			if err != nil {
				log.Printf("Connection failed: %v. Retrying in %s...", err, retryDelay)
				time.Sleep(retryDelay)
				retryDelay *= 2
				if retryDelay > 30*time.Second {
					retryDelay = 30 * time.Second
				}
				continue
			}
			retryDelay = 1 * time.Second // Reset retry delay on success
		}

		// Open a stream
		stream, err := client.PingStream(context.Background())
		if err != nil {
			log.Printf("Stream error: %v. Reconnecting...", err)
			conn.Close()
			conn = nil
			client = nil
			time.Sleep(retryDelay)
			continue
		}

		// Send messages continuously
		go func() {
			for {
				err := stream.Send(&pb.PingRequest{Message: "Hello, Server!"})
				if err != nil {
					log.Printf("Error sending message: %v", err)
					conn.Close()
					conn = nil
					client = nil
					break
				}
				time.Sleep(2 * time.Second) // Send every 2s
			}
		}()

		// Receive messages continuously
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("Server closed the stream")
				break
			}
			if err != nil {
				log.Printf("Stream error: %v. Reconnecting...", err)
				break
			}
			fmt.Println("Server Response:", resp.Message)
		}

		// Cleanup and reconnect
		conn.Close()
		conn = nil
		client = nil
	}
}

// Determines if an error should trigger a retry
func isRecoverableError(err error) bool {
	return true
}

// Cobra Command for starting client
var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start the gRPC client",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		host, _ := cmd.Flags().GetString("host")
		Start(port, host)
	},
}

func init() {
	ClientCmd.Flags().StringP("port", "p", "50051", "Port for the gRPC client")
	ClientCmd.Flags().StringP("host", "x", "grpc-server", "Host for the gRPC server")
}
