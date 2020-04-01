package main

import (
	"context"
	"fmt"
	"github.com/Sheenam3/x-tracer/pkg/streamclient"
	"github.com/docker/docker/client"
	//probeparser "github.com/sheenam3/tcptracer-goparser/parser"
//	"github.com/Sheenam3/x-tracer/parse/probeparser"
	"log"
	"os"
	"strings"
	"time"
)


//var Probelogtracer []probeparser.Log
//var Probelogconnect []probeparser.Log
//var Probename [2]string

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

	probeName := os.Getenv("tools")
	//parse the probenames

	Probe := strings.Split(probeName, ",")
//       log.Printf("Name of Probe Tool  %s", Probe)
/*
	for i:=0;i<2;i++ {

		log.Printf("Name of Probe Tool %d: %s", i+1, Probe[i])
	}

		Probelogtracer = probeparser.RunProbe(Probe[0])
		for n := 0; n < 10; n++ {
			//fmt.Printf("Struct %d  includes: %v\n", i, tcplog[i])
			log.Printf("Logs of %s : %v \n", Probe[0], Probelogtracer[n] )
	}

		Probelogconnect = probeparser.RunProbe(Probe[1])
		for n := 0; n < 10; n++ {
                        //fmt.Printf("Struct %d  includes: %v\n", i, tcplog[i])
                        log.Printf("Logs of %s : %v \n", Probe[1], Probelogtracer[n] )
        }



*/
	//endPoint := serverIp+":5555"
	//
	//conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	//defer conn.Close()
	//c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	//name := "hello jace"
	//if len(os.Args) > 1 {
	//	name = os.Args[1]
	//}
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	//if err != nil {
	//	log.Fatalf("could not greet: %v", err)
	//}
	//log.Printf("Greeting: %s", r.Message)

	cli, err := client.NewClientWithOpts(client.WithHost("unix:///var/run/docker.sock"), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	topResult, err := cli.ContainerTop(context.Background(), containerId, []string{"o","pid"})
	if err != nil {
		panic(err)
	}
	fmt.Println(topResult.Processes)


	log.Printf("Start new client")

	testClient := streamclient.New("6666",serverIp)
	testClient.StartClient(Probe,topResult.Processes)

	for {
		fmt.Println("x-agent - Sleeping")
		time.Sleep(10 * time.Second)
	}

}
