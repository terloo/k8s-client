package main

import (
	"fmt"
	"github.com/terloo/k8s-client/config"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
)

func main() {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config.GenConfig())
	if err != nil {
		panic(err)
	}

	// 获取所有apiGroup和Resource
	APIGroup, APIResourceListSlice, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		panic(err)
	}

	fmt.Printf("APIGroup: \n\n %v\n\n\n\n", APIGroup)

	for _, singleAPIResourceList := range APIResourceListSlice {

		// groupVersion是一个字符串，如"apps/v1"
		groupVersion := singleAPIResourceList.GroupVersion

		// 将字符串解析为groupVersion结构体
		gv, err := schema.ParseGroupVersion(groupVersion)
		if err != nil {
			panic(err)
		}

		fmt.Println("*****************************************************************")
		fmt.Printf("GV string [%v]\nGV struct [%#v]\nresources :\n\n", groupVersion, gv)

		// APIResources是个切片，保存着该groupVersion下所有Resouce
		for _, singleAPIResource := range singleAPIResourceList.APIResources {
			fmt.Printf("%v\n", singleAPIResource.Name)
		}

	}

}
