package service

import (
    "context"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/k8s"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Overview struct {
    CPUPercent float64 `json:"cpu_percent"`
    MemPercent float64 `json:"mem_percent"`
    Nodes int `json:"nodes"`
    Namespaces int `json:"namespaces"`
}

func GetOverview(cluster string) (*Overview, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }

    nodes, _ := cli.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    ns, _ := cli.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

    // 默认未知
    ov := &Overview{Nodes: len(nodes.Items), Namespaces: len(ns.Items), CPUPercent: -1, MemPercent: -1}

    // 尝试 metrics（最佳努力）
    if cfg, err := k8s.Manager.GetConfig(cluster); err == nil {
        if mc, err := metricsclient.NewForConfig(cfg); err == nil {
            // 汇总 node metrics 与 capacity 计算百分比
            nms, err1 := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
            if err1 == nil && len(nodes.Items) > 0 {
                var cpuUsageNano, memUsageBytes int64
                var cpuCapacityNano, memCapacityBytes int64
                for _, nm := range nms.Items {
                    cpuUsageNano += nm.Usage.Cpu().MilliValue() * 1_000_000
                    memUsageBytes += nm.Usage.Memory().Value()
                }
                for _, n := range nodes.Items {
                    cpuCapacityNano += n.Status.Capacity.Cpu().MilliValue() * 1_000_000
                    memCapacityBytes += n.Status.Capacity.Memory().Value()
                }
                if cpuCapacityNano > 0 { ov.CPUPercent = float64(cpuUsageNano) * 100 / float64(cpuCapacityNano) }
                if memCapacityBytes > 0 { ov.MemPercent = float64(memUsageBytes) * 100 / float64(memCapacityBytes) }
            }
        }
    }
    return ov, nil
}
