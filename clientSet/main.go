package main

import (
	"context"
	"fmt"
	"github.com/terloo/k8s-client/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/pointer"
)

// clientset是一系列常用(官方)apigroup的restclient集合，方便用户操作restclient
func main() {
	kubeConfig, err := kubernetes.NewForConfig(config.GenConfig())
	if err != nil {
		panic(err)
	}
	// 使用封装好的corev1 api
	namespaces, err := kubeConfig.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, namespace := range namespaces.Items {
		fmt.Println("Name: " + namespace.Name + "\tState: " + string(namespace.Status.Phase))
	}

	// 创建一个nginx高可用服务，并使用service进行负载均衡
	// createNamespace(kubeConfig)
	// createDeployment(kubeConfig)
	createService(kubeConfig)
}

// 创建一个namespace
func createNamespace(kubeConfig *kubernetes.Clientset) {
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "client-go",
		},
	}
	result, err := kubeConfig.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Created Namespace: " + result.Name)
}

func createDeployment(kubeConfig *kubernetes.Clientset) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
			Labels: map[string]string{
				"app": "nginx",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: pointer.Int32(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "nginx",
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
							Ports: []corev1.ContainerPort{
								{
									Name:          "80port",
									ContainerPort: 80,
								},
							},
							ImagePullPolicy: "IfNotPresent",
						},
					},
				},
			},
		},
	}
	result, err := kubeConfig.AppsV1().Deployments("client-go").Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Deployment: " + result.Name)
}

func createService(kubeConfig *kubernetes.Clientset) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-service",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name: "nginx",
					Port: 80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
						StrVal: "80",
					},
					NodePort: 0,
				},
			},
			Selector: map[string]string{
				"app": "nginx",
			},
		},
	}
	result, err := kubeConfig.CoreV1().Services("client-go").Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Service: " + result.Name)
}
