# go语言k8s客户端---client-go


## 客户端类型
1. RESTClient：最基础的客户端类型，其余三种均是该类型的封装。调用k8s api时，需要指定api路径，apiGroup和apiVersion
2. ClientSet：所有可能的RESTClient的集合，提供了所有可能的api路径，apiGroup和apiVersion组合而成的客户端。
3. dynamicClient：用于操作CRD等非官方定义的资源
4. DiscoveryClient：用于获取集群的Group、Version、Resource有关资源