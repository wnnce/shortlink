package balance

import "linkhttp/pb"

const EmptyKey = ""

type HostInfo struct {
	Name string
	Addr string
}

// LoadBalance 负载均衡器
type LoadBalance interface {

	// Get 通过key获取对应的节点地址
	Get(key string) (pb.LinkServiceClient, error)

	// Add 添加主机地址
	Add(hosts ...HostInfo)

	// Delete 删除对应名称的主机
	Delete(hostname string)

	// Reset 重置主机地址
	Reset()
}
