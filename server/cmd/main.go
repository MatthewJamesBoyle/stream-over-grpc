package main

import (
	"fmt"
	"log"
	"net"

	"github.com/matthewjamesboyle/stream-over-grpc/models/generated/pb"
	"github.com/matthewjamesboyle/stream-over-grpc/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterStreamingServiceServer(s, &server.GRPCServer{})

	fmt.Println("server starting and waiting")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
