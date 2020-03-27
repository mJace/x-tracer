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
	Fulllog   string
	Pid   int64
	Time  float64
	Probe string
}

/*type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64
	NumGC        uint32
	NumGoroutine int
}
*/

var IsTracerDoneSig = make(chan bool, 1)

const (
	timestamp int = 0
)

/*func NewMonitor(duration int) {
	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)

		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		// Number of goroutines
		m.NumGoroutine = runtime.NumGoroutine()

		// Misc memory stats
		m.Alloc = rtm.Alloc
		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		m.Mallocs = rtm.Mallocs
		m.Frees = rtm.Frees

		// Live objects = Mallocs - Frees
		m.LiveObjects = m.Mallocs - m.Frees

		// GC Stats
		m.PauseTotalNs = rtm.PauseTotalNs
		m.NumGC = rtm.NumGC

		// Just encode to json and print
		b, _ := json.Marshal(m)
		fmt.Println(string(b))
	}
}*/

func RunProbe(tool string, logtcpconnect chan Log) {


	//quit := make(chan bool)
/*
	if tool == "tcptracer" {
		var Tcplog []Log
		cmd := exec.Command("./tcptracer", "-t")
		cmd.Dir = "/usr/share/bcc/tools"
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		cmd.Start()
		probePid := cmd.Process.Pid
		log.Printf("pid: %d", cmd.Process.Pid)
		probeName, err := ps.FindProcess(probePid)
		if err != nil {
			fmt.Println("Error : ", err)
			os.Exit(-1)
		}
		log.Printf("Probe Name: %v", probeName.Executable())
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

					pn := probeName.Executable()
					n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: pn}
					Tcplog = append(Tcplog, n)

				}
			}

			if num > 10 {
				for i := 0; i < 9; i++ {
					return Tcplog
					/*fmt.Printf("Struct %d  includes: %v\n", i, tcplog[i])
					fmt.Printf("Output %d: %v\n PID:%v \t| TimeStamp:%v \t | ProbeName:%v \n", i, tcplog[i].log, tcplog[i].pid, tcplog[i].time, tcplog[i].probe)
				}
				//quit <- true
			}
			num += 1

		}
	}*/

	if tool == "tcpconnect" {
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
						log.Println("Tracer has been Stopped")					
						IsTracerDoneSig <- true
						
				}
				num++
				//Logging = append(Logging, n)

			}
/*
			if num > 10 {
				for i := 0; i < 9; i++ {
					return Logging
					/*fmt.Printf("Struct %d  includes: %v\n", i, Logging[i])
					fmt.Printf("Output %d: %v\n PID:%v \t| TimeStamp:%v \t | ProbeName:%v \n", i, Logging[i].Log, Logging[i].Pid, Logging[i].Time, Logging[i].Probe)
				}

			}
			num += 1*/

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
