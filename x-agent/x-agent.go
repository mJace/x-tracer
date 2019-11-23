package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "github.com/mJace/x-tracer/x-agent/route"
	"log"
	"net"
	"os"
)

const (
	port  = ":5555"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main (){
	log.Println("Start route...")

	containerId := os.Getenv("containerId")
	if containerId == "" {
		containerId = "ec9515bb14a2"
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Println("Start x-agent server...")
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	topResult, err := cli.ContainerTop(context.Background(), containerId, []string{"o","pid"})
	if err != nil {
		panic(err)
	}
	fmt.Println(topResult.Processes)
}