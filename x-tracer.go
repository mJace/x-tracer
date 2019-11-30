package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/describe/versioned"
	"k8s.io/kubectl/pkg/describe"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	pb "github.com/mJace/x-tracer/x-agent/route"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

func getAgentPodObject(containerId string, nodeId string) *v1.Pod {
	t := true
	var user int64 = 0
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
					Name:            "busybox",
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
					},
				},
			},
		},
	}
}

func main() {
	log.Println("Start x-tracer")

	var kubeconfig *string
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

	 clientSet, err := kubernetes.NewForConfig(config)
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

	podDesc := versioned.PodDescriber{clientSet }
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


	agentPod := getAgentPodObject(containerId, targetNode)
	agentSvc := getAgentService()

	pod, err := clientSet.CoreV1().Pods(agentPod.Namespace).Create(agentPod)
	if err != nil {
		panic(err)
	}
	fmt.Println("x-tracer agent created successfully...")

	svc, err := clientSet.CoreV1().Services(agentSvc.Namespace).Create(agentSvc)
	if err != nil {
		panic(err)
	}
	fmt.Println("x-tracer agent service created successfully...")



	defer func() {
		deletePolicy := metav1.DeletePropagationForeground
		fmt.Println("delete pod and svc")
		err = clientSet.CoreV1().Pods(agentPod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		err = clientSet.CoreV1().Services(agentSvc.Namespace).Delete(svc.Name, &metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
	}()

	fmt.Println("sleep 20 second...")
	time.Sleep(10* time.Second)

	// Get svc cluster IP
	svcObj, err := clientSet.CoreV1().Services(agentSvc.Namespace).Get(svc.Name, metav1.GetOptions{})
	clusterIp := strings.SplitAfter(svcObj.String(), "ClusterIP:")[1]
	clusterIp = strings.Split(clusterIp, ",")[0]
	fmt.Println(clusterIp)

	endPoint := clusterIp+":5555"

	conn, err := grpc.Dial(endPoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "hello jace"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)


	for {
		time.Sleep(10* time.Second)
	}
	//TODO gRPC client



	// TODO pod destroyer

}