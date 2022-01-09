package main

import (
	"context"
	"fmt"

	"github.com/terloo/k8s-client/config"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

// 操作最基础的RestClient
func main() {
	// 生成访问配置文件
	kubeConfig := config.GenConfig()

	// 配置api相关配置
	// 由于此处需要访问核心组的api，所以APIPath设置为"api"，GroupVersion设置为"v1"
	kubeConfig.APIPath = "api"
	kubeConfig.GroupVersion = &schema.GroupVersion{
		Group:   "",
		Version: "v1",
	}
	// kubeConfig.GroupVersion = &corev1.SchemeGroupVersion  // 也可以使用包装好的变量代替
	// 配置序列化工具
	kubeConfig.NegotiatedSerializer = scheme.Codecs
	// 根据config生成restClient
	restClient, err := rest.RESTClientFor(kubeConfig)
	if err != nil {
		panic(err.Error())
	}
	result := &corev1.PodList{}
	// url为/api/v1/namespaces/{namespace}/pods
	restClient.Get().Namespace("default").Resource("pods").
		// 指定大小限制和序列化工具
		VersionedParams(&metav1.ListOptions{Limit: 100}, scheme.ParameterCodec).
		// 请求
		Do(context.TODO()).
		// 结果存入result
		Into(result)

	fmt.Printf("namespace\t status\t\t name\n")

	for _, d := range result.Items {
		fmt.Printf("%v\t %v\t %v\n",
			d.Namespace,
			d.Status.Phase,
			d.Name)
	}

}
