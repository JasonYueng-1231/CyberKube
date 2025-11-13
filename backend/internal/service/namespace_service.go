package service

import (
    "context"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ListNamespaces(cluster string) ([]string, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
    if err != nil { return nil, err }
    out := make([]string, 0, len(list.Items))
    for _, n := range list.Items { out = append(out, n.Name) }
    return out, nil
}

