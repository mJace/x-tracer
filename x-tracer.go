package main

import (
	"flag"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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

func getAgentPodObject(containerId string) *v1.Pod {
	return &v1.Pod{
		TypeMeta:   metav1.TypeMeta{
			Kind: "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Namespace: "default",
		},
		Spec:       v1.PodSpec{
			Containers: []v1.Container {
				{
					Name: "busybox",
					Image: "ubuntu",
					ImagePullPolicy: v1.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
					Env: []v1.EnvVar {
						{
							Name: "containerId",
							Value: containerId,
						},
						{
							Name: "tools",
							Value: "net",
						},
					},
				},
			},
		},
	}
}

func main(){
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
		fmt.Println(index,":",element.Name)
	}

	fmt.Print("Choose namespace : ")
	var nsIndex int
	_, err = fmt.Scanf("%d", &nsIndex)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---------------------------------------------")

	fmt.Println("The namespace you chose is",namespaces.Items[nsIndex].Name)

	pods, err := clientSet.CoreV1().Pods(namespaces.Items[nsIndex].Name).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("\n---------------------------------------------")
	for index, element := range pods.Items {
		fmt.Println(index,":",element.Name)
	}

	fmt.Print("Choose pod : ")
	var podIndex int
	_, err = fmt.Scanf("%d", &podIndex)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("The pod you chose is",pods.Items[podIndex].Name)

	fmt.Println("Container ID is ...")

	podObj, _ := clientSet.CoreV1().Pods(namespaces.Items[nsIndex].Name).Get(pods.Items[podIndex].Name, metav1.GetOptions{})

	var containerId string
	for index := range podObj.Status.ContainerStatuses {
		containerId = getFieldString(&podObj.Status.ContainerStatuses[index], "ContainerID")
		containerId = strings.SplitAfter(containerId,"://")[1]
		fmt.Println(containerId)
	}


	agentPod := getAgentPodObject(containerId)

	_, err = clientSet.CoreV1().Pods(agentPod.Namespace).Create(agentPod)
	if err != nil {
		panic(err)
	}
	fmt.Println("x-tracer agent created successfully...")

	//TODO gRPC server

	// TODO pod destroyer

}