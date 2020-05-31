package streamclient

import (
	"context"
	pb "github.com/Sheenam3/x-tracer/api"
	pp "github.com/Sheenam3/x-tracer/parse/probeparser"
	"google.golang.org/grpc"
	"log"
//	"strconv"
	"time"
)

type StreamClient struct {
	port string
	ip   string
}

func New(servicePort string, masterIp string) *StreamClient {
	return &StreamClient{
		servicePort,
		masterIp}
}

func (c *StreamClient) StartClient(probename []string, pidList [][]string) { //[]pp.Log) {
	connect, err := grpc.Dial(c.ip+":"+c.port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}

	defer connect.Close()

	client := pb.NewSentLogClient(connect)

	logtcpconnect := make(chan pp.Log, 1)

	go pp.RunTcpconnect(probename[1], logtcpconnect, pidList[0][0])

	go func() {

		for val := range logtcpconnect {

			/*for j := range pidList {
				for k := range pidList[j] {
					if strconv.FormatUint(uint64(val.Pid), 10) == pidList[j][k] {*/
						//log.Printf("PID: %d", pidList[j][k])

						err = c.startLogStream(client, &pb.Log{
							Pid:       val.Pid,
							ProbeName: val.Probe,
							Log:       val.Fulllog,
							TimeStamp: "TimeStamp",
						})
						if err != nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}

					//}
				//}
			//}
		}

	}()

	logtcptracer := make(chan pp.Log, 1)
	go pp.RunTcptracer(probename[0], logtcptracer, pidList[0][0])
	go func() {

		for val := range logtcptracer {
			log.Printf("logtcptracer")
		/*	for j := range pidList {
				for k := range pidList[j] {
					if strconv.FormatUint(uint64(val.Pid), 10) == pidList[j][k] {*/
						err = c.startLogStream(client, &pb.Log{
							Pid:       val.Pid,
							ProbeName: val.Probe,
							Log:       val.Fulllog,
							TimeStamp: "TimeStamp",
						})
						if err != nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}
					//}
				//}
			//}
		}

	}()

	logtcpaccept := make(chan pp.Log, 1)
	go pp.RunTcpaccept(probename[2], logtcpaccept, pidList[0][0])
	go func() {

		for val := range logtcpaccept {
		/*	for j := range pidList {
				for k := range pidList[j] /
					if strconv.FormatUint(uint64(val.Pid), 10) == pidList[j][k] {*/
						err = c.startLogStream(client, &pb.Log{
							Pid:       val.Pid,
							ProbeName: val.Probe,
							Log:       val.Fulllog,
							TimeStamp: "TimeStamp",
						})
						if err != nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}
					//}
				//}
			//}
		}

	}()


	logtcplife := make(chan pp.Log, 1)
	go pp.RunTcplife(probename[3], logtcplife, pidList[0][0])
	go func() {

		for val := range logtcplife {
		/*	for j := range pidList {
				for k := range pidList[j] /
					if strconv.FormatUint(uint64(val.Pid), 10) == pidList[j][k] {*/
						err = c.startLogStream(client, &pb.Log{
							Pid:       val.Pid,
							ProbeName: val.Probe,
							Log:       val.Fulllog,
							TimeStamp: "TimeStamp",
						})
						if err != nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}
					//}
				//}
			//}
		}

	}()

	logexecsnoop := make(chan pp.Log, 1)
	go pp.RunExecsnoop(probename[4], logexecsnoop, pidList[0][0])
	go func() {

		for val := range logexecsnoop {
						err = c.startLogStream(client, &pb.Log{
							Pid:       val.Pid,
							ProbeName: val.Probe,
							Log:       val.Fulllog,
							TimeStamp: "TimeStamp",
						})
                                                println("chala gaya") 
						if err != nil {
							log.Fatalf("startLogStream fail.err: %v", err)
						}
		}

	}()


	logbiosnoop := make(chan pp.Log, 1)
        go pp.RunBiosnoop(probename[5], logbiosnoop, pidList[0][0])
        go func() {

                for val := range logbiosnoop {
                                                err = c.startLogStream(client, &pb.Log{
                                                        Pid:       val.Pid,
                                                        ProbeName: val.Probe,
                                                        Log:       val.Fulllog,
                                                        TimeStamp: "TimeStamp",
                                                })
                                                println("chala gaya") 
                                                if err != nil {
                                                        log.Fatalf("startLogStream fail.err: %v", err)
                                                }
                }

        }()

	logcachetop := make(chan pp.Log, 1)
        go pp.RunCachetop(probename[6], logcachetop, pidList[0][0])
        go func() {

                for val := range logcachetop {
                                                err = c.startLogStream(client, &pb.Log{
                                                        Pid:       val.Pid,
                                                        ProbeName: val.Probe,
                                                        Log:       val.Fulllog,
                                                        TimeStamp: "TimeStamp",
                                                })
                                                println("chala gaya") 
                                                if err != nil {
                                                        log.Fatalf("startLogStream fail.err: %v", err)
                                                }
                }

        }()


for {

		time.Sleep(time.Duration(1) * time.Second)
}

}

func (c *StreamClient) startLogStream(client pb.SentLogClient, r *pb.Log) error {

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
