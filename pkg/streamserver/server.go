package streamserver

import (
	"strings"
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
//		fmt.Println("\n", r.Log)

	        parse := strings.Fields(string(r.Log))
//		fmt.Println("PID:",r.Pid)
		if r.ProbeName == "tcptracer"{

		//fmt.Println("ProbeName:",r.ProbeName)
                //fmt.Printf("{%s}\n", r.Log)
                fmt.Printf("{Probe:%s |Sys_Time: %s |T: %s | PID:%s | PNAME:%s |IP->%s | SADDR:%s | DADDR:%s | SPORT:%s | DPORT:%s \n",r.ProbeName,parse[0],parse[1],parse[3],parse[4],parse[5],parse[6],parse[7],parse[8],parse[9])

                }else if r.ProbeName == "tcpaccept"{

                //fmt.Println("ProbeName:",r.ProbeName)
		//fmt.Printf("{%s}\n", r.Log)
                fmt.Printf("{Probe:%s |Sys_Time: %s |T: %s | PID:%s | PNAME:%s | IP:%s | RADDR:%s | RPORT:%s | LADDR:%s | LPORT:%s \n",r.ProbeName,parse[0],parse[1],parse[3],parse[4],parse[5],parse[6],parse[7],parse[8],parse[9])

                }else if r.ProbeName == "tcplife"{

		fmt.Printf("{Probe:%s |Sys_Time: %s |PID:%s | PNAME:%s | LADDRR:%s | LPORT:%s | RADDR:%s | RPORT:%s | TX_KB:%s | RX_KB:%s | MS: %s \n",r.ProbeName,parse[0],parse[2],parse[3],parse[4],parse[5],parse[6],parse[7],parse[8],parse[9],parse[10])

		}else if r.ProbeName == "execsnoop"{
		fmt.Printf("{Probe:%s |Sys_Time: %s | T:%s | PNAME: %s | PID:%s | PPID:%s | RET:%s | ARGS:%s \n",r.ProbeName,parse[0],parse[1],parse[3],parse[4],parse[5],parse[6],parse[7])

		}else if r.ProbeName == "biosnoop"{

		fmt.Printf("{Probe:%s |Sys_Time: %s |T: %s |PNAME: %s | PID:%s | DISK:%s | R/W:%s | SECTOR:%s |BYTES: %s | Lat(ms): %s | \n",r.ProbeName,parse[0],parse[1],parse[2],parse[3],parse[4],parse[5],parse[6],parse[7],parse[9])

		}else if r.ProbeName == "cachetop"{

		fmt.Printf("{Probe:%s |Sys_Time: %s | PID:%s | UID:%s | CMD:%s | HITS:%s | MISS:%s | DIRTIES: %s| READ_HIT%:%s | W_HIT%: %s | \n",r.ProbeName,parse[0],parse[1],parse[2],parse[3],parse[5],parse[6],parse[7],parse[8], parse[9])

		}else{

                fmt.Printf("{Probe:%s |Sys_Time: %s |T: %s | PID:%s | PNAME:%s | IP:%s | SADDR:%s | DADDR:%s | DPORT:%s \n",r.ProbeName,parse[0],parse[1],parse[3],parse[4],parse[5],parse[6],parse[7],parse[8])
                }
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

