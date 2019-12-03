package main

import (
	"context"
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/describe"
	"k8s.io/kubectl/pkg/describe/versioned"
	"net"

	pb "github.com/mJace/x-tracer/route"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"syscall"
	"time"
)

var clientSet *kubernetes.Clientset
var pod *v1.Pod
var svc *v1.Service

const (
	port  = ":5555"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}


func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func getHostName() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	log.Println("Hostname : ",name)
	return name
}

func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		fmt.Println("Delete agent pod and service")
		_ = clientSet.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
		_ = clientSet.CoreV1().Services(svc.Namespace).Delete(svc.Name, &metav1.DeleteOptions{})
		os.Exit(0)
	}()
}

func getFieldString(e *v1.ContainerStatus, field string) string {
	r := reflect.ValueOf(e)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func getAgentService() *v1.Service {
	return &v1.Service{
		TypeMeta:   metav1.TypeMeta{
			Kind: "service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "agent-service",
			Namespace: "default",
		},
		Spec:       v1.ServiceSpec{
			Selector: map[string]string{
				"app": "x-agent",
			},
			Ports: []v1.ServicePort{
				{
					Name: "grpc",
					Protocol: "TCP",
					Port: 5555,
				},
			},
		},
	}
}

func getAgentPodObject(containerId string, nodeId string, masterIp string) *v1.Pod {
	t := true
	var user int64 = 0
	var pathType = v1.HostPathDirectory
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "x-agent",
			Namespace: "default",
			Labels: map[string]string{
				"app" : "x-agent",
			},
		},
		Spec: v1.PodSpec{
			NodeSelector: map[string]string{

					"kubernetes.io/hostname" : nodeId,

			},
			Containers: []v1.Container{
				{
					Name:            "agent",
					Image:           "mjace/x-agent",
					Ports: []v1.ContainerPort{
						{
							Name:          "grpc",
							ContainerPort: 5555,
							Protocol:      "TCP",
						},
					},
					ImagePullPolicy: v1.PullIfNotPresent,
					SecurityContext: &v1.SecurityContext{
						Privileged:               &t,
						RunAsUser:                &user,
					},
					Env: []v1.EnvVar{
						{
							Name:  "containerId",
							Value: containerId,
						},
						{
							Name:  "tools",
							Value: "net",
						},
						{
							Name: "masterIp",
							Value: masterIp,
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							MountPath: "/lib/modules",
							Name: "kernel-modules",
						},
						{
							MountPath: "/usr/src",
							Name: "kernel-src",
						},
						{
							MountPath: "/etc/localtime",
							Name: "localtime",
						},
						//{
						//	MountPath: "/sys",
						//	Name: "sys",
						//},
						{
							MountPath: "/var/run/docker.sock",
							Name: "docker-sock",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "kernel-modules",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/lib/modules",
							Type: &pathType,
						},
					},
				},
				{
					Name: "kernel-src",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/usr/src",
							Type: &pathType,
						},
					},
				},
				{
					Name: "localtime",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/etc/localtime",
						},
					},
				},
				//{
				//	Name: "sys",
				//	VolumeSource: v1.VolumeSource{
				//		HostPath: &v1.HostPathVolumeSource{
				//			Path: "/sys",
				//			Type: &pathType,
				//		},
				//	},
				//},
				{
					Name: "docker-sock",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/var/run/docker.sock",
							Type: &pathType,
						},
					},
				},
			},
		},
	}
}

func main() {
	log.Println("Start x-tracer")

	var kubeconfig *string
	var debug *bool
	debug = flag.Bool("kind", false, "for kind env.")
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

	podDesc := versioned.PodDescriber{Interface: clientSet }
	descStr, err :=podDesc.Describe(podObj.Namespace, podObj.Name, describe.DescriberSettings{ShowEvents:false})

	descStr = strings.SplitAfter(descStr, "Node:")[1]
	descStr = strings.Split(descStr, "/")[0]
	reg := regexp.MustCompile("[^\\s]+")
	targetNode := reg.FindAllString(descStr,1)[0]


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
	nodeIp := strings.Split(currentNode.Status.Addresses[0].Address," ")[0]

	// Initial pod and service.
	agentPod := getAgentPodObject(containerId, targetNode, nodeIp)
	agentSvc := getAgentService()

	pod, err = clientSet.CoreV1().Pods(agentPod.Namespace).Create(agentPod)
	if err != nil {
		panic(err)
	}
	fmt.Println("x-tracer agent created successfully...")

	svc, err = clientSet.CoreV1().Services(agentSvc.Namespace).Create(agentSvc)
	if err != nil {
		panic(err)
	}
	fmt.Println("x-tracer agent service created successfully...")


	// Get svc cluster IP
	svcObj, err := clientSet.CoreV1().Services(agentSvc.Namespace).Get(svc.Name, metav1.GetOptions{})
	clusterIp := strings.SplitAfter(svcObj.String(), "ClusterIP:")[1]
	clusterIp = strings.Split(clusterIp, ",")[0]
	fmt.Println(clusterIp)

	SetupCloseHandler()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Println("Start x-agent server...")
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Run our program... We create a file to clean up then sleep
	for {
		fmt.Println("- Sleeping")
		time.Sleep(10 * time.Second)
	}

}