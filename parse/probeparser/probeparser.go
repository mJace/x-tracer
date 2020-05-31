package probeparser

import (
	"bufio"
	//"encoding/json"
	"fmt"
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


func GetNS(pid string) string {
	cmdName := "ls"
        out, err := exec.Command(cmdName, fmt.Sprintf("/proc/%s/ns/net", pid), "-al").Output()
        if err != nil {
                println(err)
        }
        ns := string(out)
        parse := strings.Fields(string(ns))
        s := strings.SplitN(parse[10], "[", 10)
        sep := strings.Split(s[1], "]")
        return sep[0]

}
func RunTcptracer(tool string, logtcptracer chan Log, pid string) {

	sep := GetNS(pid)
/*	ppid := pid
	cmdName := "ls"
	out, err := exec.Command(cmdName, fmt.Sprintf("/proc/%s/ns/net", ppid), "-al").Output()
	if err != nil {
		println(err)
	}
	ns := string(out)
	parse := strings.Fields(string(ns))
//	fmt.Printf("%q\n", strings.SplitN(parse[10], "[", 10))
	s := strings.SplitN(parse[10], "[", 10)
	sep := strings.Split(s[1], "]")
*/
	cmd := exec.Command("./tcptracer.py","-T","-t","-N" + sep)
	cmd.Dir = "/usr/share/bcc/tools/ebpf"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
	//num := 1

	for {

		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))
		//println("TCP TRACER", parsedLine[0])
		if parsedLine[0] != "Tracing" {
			if parsedLine[0] != "TIME(s)" {
				ppid, err := strconv.ParseInt(parsedLine[3], 10, 64)
				if err != nil {
					println("Tcptracer PID Error")
				}
				/*timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
				if err != nil {
					println(" Tcptracer Timestamp Error")
				}*/
				timest := 0.00
				n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
				logtcptracer <- n
				//if num > 5000 {
				//	close(logtcptracer)
				//	log.Println("Tracer has been Stopped")

				//}
				//num++

			}
		}
	}
}

func RunTcpconnect(tool string, logtcpconnect chan Log, pid string ) {

	sep := GetNS(pid)
	cmd := exec.Command("./tcpconnect.py", "-T","-t","-N" + sep)
	cmd.Dir = "/usr/share/bcc/tools/ebpf"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
//	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))
		//println(parsedLine[0])
		if parsedLine[0] != "TIME(s)" {
			ppid, err := strconv.ParseInt(parsedLine[3], 10, 64)
			if err != nil {
				println("TCPConnect PID Error")
			}
			/*timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" TCPConnect Timestamp Error")
			}*/
			timest := 0.00
			n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
			logtcpconnect <- n
/*			if num > 5000 {
				close(logtcpconnect)

			}*/
			//num++
		}
	}
}

func RunTcpaccept(tool string, logtcpaccept chan Log, pid string) {

	sep := GetNS(pid)
	cmd := exec.Command("./tcpaccept.py", "-T","-t","-N" + sep)
	cmd.Dir = "/usr/share/bcc/tools/ebpf"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
//	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))

		if parsedLine[0] != "TIME(s)" {
			ppid, err := strconv.ParseInt(parsedLine[3], 10, 64)
			if err != nil {
				println("TCPaccept PID Error")
			}
/*			timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" TCPaccept Timestamp Error")
			}*/
			timest := 0.00

			n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
			logtcpaccept <- n
/*			if num > 5000 {
				close(logtcpaccept)
			}
			num++*/

		}
	}
}


func RunTcplife(tool string, logtcplife chan Log, pid string) {

	sep := GetNS(pid)
	cmd := exec.Command("./tcplife.py", "-T","-N" + sep)
	cmd.Dir = "/usr/share/bcc/tools/ebpf"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
//	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))
		println(parsedLine[0])
		if parsedLine[0] != "TIME(s)" {
			ppid, err := strconv.ParseInt(parsedLine[2], 10, 64)
			if err != nil {
				println("TCPlife PID Error")
			}
/*			timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" TCPlife Timestamp Error")
			}*/
			timest := 0.00

			n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
			logtcplife <- n
/*			if num > 5000 {
				close(logtcpaccept)
			}
			num++*/

		}
	}
}


func RunExecsnoop(tool string, logexecsnoop chan Log, pid string) {

	sep := GetNS(pid)
	cmd := exec.Command("./execsnoop.py", "-T" ,"-t","-N" + sep)
	cmd.Dir = "/usr/share/bcc/tools/ebpf"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout)
//	num := 1

	for {
		line, _, _ := buf.ReadLine()
		parsedLine := strings.Fields(string(line))

		//if parsedLine[0] == "TIME(s)" {
		ppid, err := strconv.ParseInt(parsedLine[4], 10, 64)
		if err != nil {
			println("Execsnoop PID Error")
		}
/*			timest, err := strconv.ParseFloat(parsedLine[timestamp], 64)
			if err != nil {
				println(" Execsnoop Timestamp Error")
			}*/
		timest := 0.00

		n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
		logexecsnoop <- n
/*			if num > 5000 {
				close(logtcpaccept)
			}
			num++*/

		//}
	}
}



func RunBiosnoop(tool string, logbiosnoop chan Log, pid string) {

        sep := GetNS(pid)
        cmd := exec.Command("./biosnoop.py", "-T", "-N" + sep)
        cmd.Dir = "/usr/share/bcc/tools/ebpf"
        stdout, err := cmd.StdoutPipe()
        if err != nil {
                log.Fatal(err)
        }
        cmd.Start()
        buf := bufio.NewReader(stdout)


        for {
                line, _, _ := buf.ReadLine()
                parsedLine := strings.Fields(string(line))


                ppid, err := strconv.ParseInt(parsedLine[3], 10, 64)
                if err != nil {
                                println("Biosnoop PID Error")
                }
                timest := 0.00

                n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
                logbiosnoop <- n

	}
}	


func RunCachetop(tool string, logcachetop chan Log, pid string) {

        sep := GetNS(pid)
        cmd := exec.Command("./Cachetop.py", "-T", "-N" + sep)
        cmd.Dir = "/usr/share/bcc/tools/ebpf"
        stdout, err := cmd.StdoutPipe()
        if err != nil {
                log.Fatal(err)
        }
        cmd.Start()
        buf := bufio.NewReader(stdout)


        for {
                line, _, _ := buf.ReadLine()
                parsedLine := strings.Fields(string(line))


                ppid, err := strconv.ParseInt(parsedLine[1], 10, 64)
                if err != nil {
                     println("Cachetop PID Error")
                }
                timest := 0.00

                n := Log{Fulllog: string(line), Pid: ppid, Time: timest, Probe: tool}
                logcachetop <- n

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
