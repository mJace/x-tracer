package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	pb "github.com/mJace/x-tracer/api"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main (){
	log.Println("Start api...")

	containerId := os.Getenv("containerId")
	if containerId == "" {
		containerId = "ec9515bb14a2"
	}

	serverIp := os.Getenv("masterIp")
	if containerId == "" {
		containerId = "ec9515bb14a2"
	}

	endPoint := serverIp+":5555"

	conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "hello jace"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	topResult, err := cli.ContainerTop(context.Background(), containerId, []string{"o","pid"})
	if err != nil {
		panic(err)
	}
	fmt.Println(topResult.Processes)

	for {
		fmt.Println("- Sleeping")
		time.Sleep(10 * time.Second)
	}

}