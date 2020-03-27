package streamserver

import (
	"fmt"
	pb "github.com/mJace/x-tracer/api"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)


type StreamServer struct {
	port string
}

func (s *StreamServer) RouteLog(stream pb.SentLog_RouteLogServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.Response{
				Res:                  "Stream closed",
			})
		}
		if err != nil {
			return err
		}
		fmt.Println("\n", r.Log)
		fmt.Println("PID:",r.Pid)
		fmt.Println("ProbeName: ",r.ProbeName)
		//fmt.Println(r.TimeStamp, "\n")
	}
}

func New(servicePort string) *StreamServer{
	return &StreamServer{
		servicePort}
}

func (s *StreamServer) StartServer(){
	server := grpc.NewServer()
	pb.RegisterSentLogServer(server, &StreamServer{})

	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalln("net.Listen error:", err)
	}

	_ = server.Serve(lis)
}

