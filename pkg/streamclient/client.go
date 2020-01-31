package streamclient

import (
	"context"
	"google.golang.org/grpc"
	pb "github.com/mJace/x-tracer/api"
	"log"
)

type StreamClient struct {
	port string
}

func New(servicePort string) *StreamClient{
	return &StreamClient{
		servicePort}
}


func (c *StreamClient) StartClient () {
	connect, err := grpc.Dial(":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer connect.Close()

	client := pb.NewSentLogClient(connect)
	err = c.startLogStream(client, &pb.Log{
		Pid:                  3422,
		ProbeName:            "net",
		Log:                  "test log 123",
		TimeStamp:            "local current time",
	})


}

func (c *StreamClient) startLogStream(client pb.SentLogClient, r *pb.Log) error {
	stream, err := client.RouteLog(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n<6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("Response: %v", resp.Res)
	return nil

}