package agentmanager

import (
//	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type agent struct {
	targetContainerId string
	targetNodeId      string
	masterIp          string
	clientSet *kubernetes.Clientset
}

var podObj *v1.Pod
var svcObj *v1.Service

func New(containerId string, nodeId string, masterIp string, clientSet *kubernetes.Clientset) *agent {
	return &agent{
		containerId,
		nodeId,
		masterIp,
		clientSet}
}

func (a *agent) getAgentService() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "agent-service",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": "x-agent",
			},
			Ports: []v1.ServicePort{
				{
					Name:     "grpc",
					Protocol: "TCP",
					Port:     5555,
				},
			},
		},
	}
}

func (a *agent) getAgentPodObject() *v1.Pod {
	t := true
	var user int64 = 0
	net := "tcptracer,tcpconnect,tcpaccept,tcplife,execsnoop,biosnoop,cachetop"
	//var pathType = v1.HostPathDirectory
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "x-agent",
			Namespace: "default",
			Labels: map[string]string{
				"app": "x-agent",
			},
		},
		Spec: v1.PodSpec{
			NodeSelector: map[string]string{

				"kubernetes.io/hostname": a.targetNodeId,
			},
			ShareProcessNamespace: &t,
			Containers: []v1.Container{
				{
					Name:  "agent",
					Image: "sheenam3/x-agent",
					/*Command: []string{
						 "sleep",
						 "9900000" },*/
					Ports: []v1.ContainerPort{
						{
							Name:          "grpc",
							ContainerPort: 5555,
							Protocol:      "TCP",
						},
					},
					ImagePullPolicy: v1.PullAlways,
					SecurityContext: &v1.SecurityContext{
						Privileged: &t,
						RunAsUser:  &user,
					},
					Env: []v1.EnvVar{
						{
							Name:  "containerId",
							Value: a.targetContainerId,
						},
						{
							Name:  "tools",
							Value: net,
						},
						{
							Name:  "masterIp",
							Value: a.masterIp,
						},
					},
					VolumeMounts: []v1.VolumeMount{
						{
							MountPath: "/proc",
							Name: "host-proc",

						},
						{
							MountPath: "/lib/modules",
							Name:      "kernel-modules",
						},
						{
							MountPath: "/usr/src",
							Name:      "kernel-src",
						},
						{
							MountPath: "/etc/localtime",
							Name:      "localtime",
						},
						{
							MountPath: "/var/run/docker.sock",
							Name:      "docker-sock",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "host-proc",
					VolumeSource:v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/proc",
						},
					},
				},


				{
					Name: "kernel-modules",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/lib/modules",
						},
					},
				},
				{
					Name: "kernel-src",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/usr/src",
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
				{
					Name: "docker-sock",
					VolumeSource: v1.VolumeSource{
						HostPath: &v1.HostPathVolumeSource{
							Path: "/var/run/docker.sock",
						},
					},
				},
			},
		},
	}
}

func (a *agent) ApplyAgentPod(){
	agentPod := a.getAgentPodObject()
	podObj, _ = a.clientSet.CoreV1().Pods(agentPod.Namespace).Create(agentPod)
}

func (a *agent) ApplyAgentService() {
	agentSvc := a.getAgentService()
	svcObj, _ = a.clientSet.CoreV1().Services(agentSvc.Namespace).Create(agentSvc)
}

func (a *agent) SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		_ = a.clientSet.CoreV1().Pods(podObj.Namespace).Delete(podObj.Name, &metav1.DeleteOptions{})
		_ = a.clientSet.CoreV1().Services(svcObj.Namespace).Delete(svcObj.Name, &metav1.DeleteOptions{})
		os.Exit(0)
	}()
}

func (a *agent) GetServiceClusterIp() string {
	clusterIp := strings.SplitAfter(svcObj.String(), "ClusterIP:")[1]
	clusterIp = strings.Split(clusterIp, ",")[0]
	return clusterIp
}
