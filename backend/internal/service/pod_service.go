package service

import (
    "context"
    "io"
    "strconv"
    "strings"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
