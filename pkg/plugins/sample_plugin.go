package plugins

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework"
)

// name of plugin
const Name = "sample-plugin"

type Sample struct {
	handle framework.FrameworkHandle
}

func (s *Sample) Name() string {
	return Name
}

func getNodeNameHash(s string) int {
	value := int(crc32.ChecksumIEEE([]byte(s)))
	if value >= 0 {
		return value
	} 

	if -value >= 0 {
		return -value
	}

	return 100
}


func (s *Sample) PreFilter(pc *framework.PluginContext, pod *v1.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Filter(pc *framework.PluginContext, pod *v1.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	return framework.NewStatus(framework.Success, "")
}

// score based on node name
func Score(pc *framework.PluginContext, p *v1.Pod, nodeName string) (int, *framework.Status) {
	hashValue := getNodeNameHash(nodeName) % 100
	rand.Seed(time.Now().UnixNano())
	score := (hashValue + rand.Intn(100)) % 100
	klog.V(3).Infof("node: %v, score: %d", nodeName, score)
	return score, framework.NewStatus(framework.Success, "")
}


func (s *Sample) PreBind(pc *framework.PluginContext, pod *v1.Pod, nodeName string) *framework.Status {
	if nodeInfo, ok := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]; !ok {
		return framework.NewStatus(framework.Error, fmt.Sprintf("prebind get node info error: %+v", nodeName))
	} else {
		klog.V(3).Infof("prebind node info: %+v", nodeInfo.Node())
		return framework.NewStatus(framework.Success, "")
	}
}

//type PluginFactory = func(f FrameworkHandle) (Plugin, error)
func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	return &Sample{
		handle: f,
	}, nil
}
