package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/Sheenam3/x-tracer/api"
	"github.com/Sheenam3/x-tracer/pkg/streamserver"
	"github.com/Sheenam3/x-tracer/internal/agentmanager"
	//"google.golang.org/grpc"
	//"google.golang.org/grpc/reflection"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/describe"
	"k8s.io/kubectl/pkg/describe/versioned"
	"log"
	//"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getFieldString(e *v1.ContainerStatus, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func getHostName() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Println("Hostname : ", name)
	return name
}

var clientSet *kubernetes.Clientset

const (
	port  = ":5555"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	log.Println("Start x-tracer")

	var kubeconfig *string
	//var debug *bool
	debug := flag.Bool("kind", false, "for kind env.")
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespaces, err := clientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\n---------------------------------------------")
	for index, element := range namespaces.Items {
		fmt.Println(index, ":", element.Name)
	}

	fmt.Print("Choose namespace : ")
	var nsIndex int
	_, err = fmt.Scanf("%d", &nsIndex)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---------------------------------------------")

	fmt.Println("The namespace you chose is", namespaces.Items[nsIndex].Name)

	pods, err := clientSet.CoreV1().Pods(namespaces.Items[nsIndex].Name).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\n---------------------------------------------")
	for index, element := range pods.Items {
		fmt.Println(index, ":", element.Name)
	}

	fmt.Print("Choose pod : ")
	var podIndex int
	_, err = fmt.Scanf("%d", &podIndex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("The pod you chose is", pods.Items[podIndex].Name)

	fmt.Println("Container ID is ...")

	podObj, _ := clientSet.CoreV1().Pods(namespaces.Items[nsIndex].Name).Get(pods.Items[podIndex].Name, metav1.GetOptions{})

	podDesc := versioned.PodDescriber{Interface: clientSet}
	descStr, err := podDesc.Describe(podObj.Namespace, podObj.Name, describe.DescriberSettings{ShowEvents: false})

	descStr = strings.SplitAfter(descStr, "Node:")[1]
	descStr = strings.Split(descStr, "/")[0]
	reg := regexp.MustCompile("[^\\s]+")
	targetNode := reg.FindAllString(descStr, 1)[0]

	var containerId string
	for index := range podObj.Status.ContainerStatuses {
		containerId = getFieldString(&podObj.Status.ContainerStatuses[index], "ContainerID")
		containerId = strings.SplitAfter(containerId, "://")[1]
		fmt.Println(containerId)
	}

	var currentNode *v1.Node

	if *debug {
		currentNode, err = clientSet.CoreV1().Nodes().Get("kind-control-plane", metav1.GetOptions{})
	} else {
		currentNode, err = clientSet.CoreV1().Nodes().Get(getHostName(), metav1.GetOptions{})
	}

	if currentNode == nil {
		panic("current node can not be nil")
	}
	nodeIp := strings.Split(currentNode.Status.Addresses[0].Address, " ")[0]

	agent := agentmanager.New(containerId, targetNode, nodeIp, clientSet)
	fmt.Println("Start Agent Pod")
	agent.ApplyAgentPod()

	fmt.Println("Start Agent Service")
	agent.ApplyAgentService()

	agent.SetupCloseHandler()

	//lis, err := net.Listen("tcp", port)
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}
	//s := grpc.NewServer()
	//log.Println("Start x-agent server...")
	//pb.RegisterGreeterServer(s, &server{})
	//// Register reflection service on gRPC server.
	//reflection.Register(s)
	//if err := s.Serve(lis); err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}

	server := streamserver.New("6666")
	server.StartServer()

	// Run our program... We create a file to clean up then sleep
	for {
		fmt.Println("From x-tracer- Sleeping")
		time.Sleep(10 * time.Second)
	}


}
