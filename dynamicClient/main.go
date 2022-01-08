package main

import (
	"context"
	"fmt"
	"github.com/terloo/k8s-client/config"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"strconv"
)

// dynamicClient，处理非官方资源(CRD)所使用的client
func main() {
	dynamicClient, err := dynamic.NewForConfig(config.GenConfig())
	if err != nil {
		panic(err)
	}

	// 需要操作的CRD的apiGroup和Version，这里使用Deployment来代替CRD
	resource := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	// 查询CRD
	unstructObj, err := dynamicClient.Resource(resource).Namespace("client-go").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// 返回的是一个Unstructured(List)对象，可以将其转为指定的对象

	// 转换
	deployList := &appsv1.DeploymentList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), deployList)

	if err != nil {
		panic(err)
	}

	for _, deploy := range deployList.Items {
		fmt.Println("Deployment name: " + deploy.Name + " ReadyReplicas: " + strconv.Itoa(int(deploy.Status.ReadyReplicas)))
	}
}
