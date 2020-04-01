package probeparser

import (
	"bufio"
	//"encoding/json"
	//"fmt"
	"log"
	//"os"
	"os/exec"
	//"runtime"
	"strconv"
	"strings"
	//"time"
)

type Log struct {
	Fulllog string
	Pid     int64
	Time    float64
	Probe   string
}

const (
	timestamp int = 0
)

func RunTcptracer(tool string, logtcptracer chan Log) {

	cmd := exec.Command("./tcptracer", "-t")
	cmd.Dir = "/usr/share/bcc/tools"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
	num := 1

	for {

		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))
		//println("TCP TRACER", parsedLine[0])
		if parsedLine[0] != "Tracing" {
			if parsedLine[0] != "TIME(s)" {
				ppid, err := strconv.ParseInt(parsedLine[2], 10, 64)
				if err != nil {
					println("Tcptracer PID Error")
				}
				timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
				if err != nil {
					println(" Tcptracer Timestamp Error")
				}
				n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
				logtcptracer <- n
				if num > 1000 {
					close(logtcptracer)
					log.Println("Tracer has been Stopped")

				}
				num++

			}
		}
	}
}

func RunTcpconnect(tool string, logtcpconnect chan Log) {

	cmd := exec.Command("./tcpconnect", "-t")
	cmd.Dir = "/usr/share/bcc/tools"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))
		//println(parsedLine[0])
		if parsedLine[0] != "TIME(s)" {
			ppid, err := strconv.ParseInt(parsedLine[1], 10, 64)
			if err != nil {
				println("TCPConnect PID Error")
			}
			timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" TCPConnect Timestamp Error")
			}

			n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
			logtcpconnect <- n
			if num > 300 {
				close(logtcpconnect)

			}
			num++
		}
	}
}

func RunTcpaccept(tool string, logtcpaccept chan Log) {

	cmd := exec.Command("./tcpaccept", "-t")
	cmd.Dir = "/usr/share/bcc/tools"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))

		if parsedLine[0] != "TIME(s)" {
			ppid, err := strconv.ParseInt(parsedLine[1], 10, 64)
			if err != nil {
				println("TCPaccept PID Error")
			}
			timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" TCPaccept Timestamp Error")
			}

			n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
			logtcpaccept <- n
			if num > 300 {
				close(logtcpaccept)
			}
			num++

		}
	}
}

/*func main() {
	//go RunTCP("tcptracer")

	logtcpconnect := make(chan Log, 1)

	go RunProbe("tcpconnect", logtcpconnect)
	for val := range logtcpconnect {
	log.Printf("%v Probe: %s, Pid: %d", val.Fulllog, val.Probe, val.Pid)

	}

	for
	{

		time.Sleep(10 * time.Second)
	}
}*/