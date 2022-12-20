package helpers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"regexp"
)

const hostPattern="[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\\.?"

func showLabel(key string ) bool{
	return !regexp.MustCompile(hostPattern).MatchString(key )
}

//过滤 要显示的标签
func FilterLabels(labels map[string]string) (ret []string){
	for k, v := range labels{
		if showLabel(k){
			ret=append(ret, fmt.Sprintf("%s=%s",k,v))
		}
	}
	return
}

func FilterTaints(taints []v1.Taint) (ret []string){
	for _, taint := range taints {
		if showLabel(taint.Key){
			ret = append(ret,fmt.Sprintf("%s=%s:%s", taint.Key,taint.Value,taint.Effect))
		}
	}
	return
}

// 第一个是cpu使用 第二个是内存使用
func GetNodeUsage(c *versioned.Clientset,node *v1.Node) []float64{
	nodeMetric, _ := c.MetricsV1beta1().
		NodeMetricses().Get(context.Background(),node.Name,metav1.GetOptions{})
	cpu := float64(nodeMetric.Usage.Cpu().MilliValue()) / float64(node.Status.Capacity.Cpu().MilliValue())
	memory := float64(nodeMetric.Usage.Memory().MilliValue()) / float64(node.Status.Capacity.Memory().MilliValue())
	return []float64{cpu, memory}
}
