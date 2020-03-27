package streamclient

import (
	"context"
	"google.golang.org/grpc"
	pb "github.com/Sheenam3/x-tracer/api"
//	"github.com/Sheenam3/x-tracer/cmd/x-agent"
	pp "github.com/Sheenam3/x-tracer/parse/probeparser"
	"log"
	"time"	
//	"strconv"
)

type StreamClient struct {
	port string
	ip string
}

func New(servicePort string, masterIp string) *StreamClient{
	return &StreamClient{
		servicePort,
		masterIp}
}


func (c *StreamClient) StartClient(probename []string){  //[]pp.Log) {
	connect, err := grpc.Dial(c.ip+":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer connect.Close()

	client := pb.NewSentLogClient(connect)

	logtcpconnect := make(chan pp.Log, 1)
	
	go pp.RunTcpconnect(probename[1], logtcpconnect)
//	log.Printf("After probe")
	go func() {
//		pp.RunProbe(probename[1], logtcpconnect)
	//	for i := 0;i<100;i++ {
	//	val := <-logtcpconnect
			for val := range logtcpconnect {
//			log.Printf("%v Probe: %s, Pid: %d", val.Fulllog, val.Probe, val.Pid)
	//		for j:= range pidList {
	//			for k:= range pidList[j] {
	//				if strconv.FormatUint(uint64(val.Pid), 10) == pidList[j][k] {
						//log.Printf("PID: %d", pidList[j][k])

						err = c.startLogStream(client, &pb.Log{
						Pid:                  val.Pid,
						ProbeName:            val.Probe,
						Log:                  val.Fulllog,
						TimeStamp:            "TimeStamp",
						})
						if err!= nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}

	//				}
	//			}
			}
	//	}

	}()


	logtcptracer := make(chan pp.Log, 1)
	go pp.RunTcptracer(probename[0], logtcptracer)
	go func() {

		for val := range logtcptracer {
			log.Printf("logtcptracer")
			err = c.startLogStream(client, &pb.Log{
				Pid:                  val.Pid,
				ProbeName:            val.Probe,
				Log:                  val.Fulllog,
				TimeStamp:            "TimeStamp",
			})
			if err!= nil {
			log.Fatalf("startLogStream fail.err: %v", err)
			}
	
		}

	}()

for {
        //log.Printf("[main] Call tcptracer.Stop() in %d seconds\n", i)
        time.Sleep(time.Duration(1) * time.Second)
    }




}
	



func (c *StreamClient) startLogStream(client pb.SentLogClient, r *pb.Log) error {

log.Printf("Inside StartLog")	
	stream, err := client.RouteLog(context.Background())
	if err != nil {
		return err
	}
	
	err = stream.Send(r)
		if err != nil {
			return err
		}
	

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("Response from the Server ;) : %v", resp.Res)
	return nil

}
