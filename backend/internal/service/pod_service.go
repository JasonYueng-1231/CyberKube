package service

import (
    "context"
    "io"
    "sort"
    "strconv"
    "strings"
    "time"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

type PodItem struct {
    Name      string `json:"name"`
    Namespace string `json:"namespace"`
    Phase     string `json:"phase"`
    Ready     string `json:"ready"`
    Restarts  int32  `json:"restarts"`
    NodeName  string `json:"node_name"`
    Containers []string `json:"containers"`
}

func ListPods(cluster, namespace string) ([]PodItem, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    ctx := context.TODO()
    list, err := cli.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
    if err != nil { return nil, err }
    var out []PodItem
    for _, p := range list.Items {
        ready := 0
        restarts := int32(0)
        for _, cs := range p.Status.ContainerStatuses {
            if cs.Ready { ready++ }
            restarts += cs.RestartCount
        }
        // 容器名
        var containers []string
        for _, c := range p.Spec.Containers { containers = append(containers, c.Name) }
        out = append(out, PodItem{
            Name: p.Name, Namespace: p.Namespace, Phase: string(p.Status.Phase),
            Ready: strings.Join([]string{strconv.Itoa(ready), strconv.Itoa(len(p.Status.ContainerStatuses))}, "/"),
            Restarts: restarts, NodeName: p.Spec.NodeName, Containers: containers,
        })
    }
    return out, nil
}

func GetPodLogs(cluster, namespace, name string, tail int64) (string, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return "", err }
    ctx := context.TODO()
    opt := &corev1.PodLogOptions{}
    if tail > 0 { opt.TailLines = &tail }
    req := cli.CoreV1().Pods(namespace).GetLogs(name, opt)
    rc, err := req.Stream(ctx)
    if err != nil { return "", err }
    defer rc.Close()
    b, _ := io.ReadAll(rc)
    return string(b), nil
}

func GetPod(cluster, namespace, name string) (*corev1.Pod, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func ListPodEvents(cluster, namespace, name string) ([]corev1.Event, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.kind=Pod,involvedObject.name="+name})
    if err != nil { return nil, err }
    items := list.Items
    sort.SliceStable(items, func(i, j int) bool {
        ti := eventTimeOrLast(&items[i])
        tj := eventTimeOrLast(&items[j])
        // 按时间倒序
        return tj.Before(ti)
    })
    return items, nil
}

func eventTimeOrLast(e *corev1.Event) time.Time {
    if !e.LastTimestamp.IsZero() {
        return e.LastTimestamp.Time
    }
    if !e.EventTime.IsZero() {
        return e.EventTime.Time
    }
    if !e.CreationTimestamp.IsZero() {
        return e.CreationTimestamp.Time
    }
    return time.Time{}
<<<<<<< HEAD
}

type PodDetail struct {
    Pod    *corev1.Pod    `json:"pod"`
    Events []corev1.Event `json:"events"`
    Yaml   string         `json:"yaml"`
}

// GetPodDetail 聚合 Pod 详情、事件与 YAML，方便前端一次获取
func GetPodDetail(cluster, namespace, name string) (*PodDetail, error) {
    pod, err := GetPod(cluster, namespace, name)
    if err != nil { return nil, err }
    evs, _ := ListPodEvents(cluster, namespace, name) // 事件失败不阻断整体
    yml, _ := yaml.Marshal(pod)
    return &PodDetail{Pod: pod, Events: evs, Yaml: string(yml)}, nil
=======
>>>>>>> origin/develop
}
