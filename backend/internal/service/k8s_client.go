package service

import (
    "context"
    "fmt"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/config"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/database"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/k8s"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/model"
    "github.com/JasonYueng-1231/CyberKube/backend/pkg/encrypt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

// GetClientForCluster 确保并返回指定集群的 clientset
func GetClientForCluster(name string) (*kubernetes.Clientset, error) {
    if cli, err := k8s.Manager.Get(name); err == nil {
        return cli, nil
    }
    // 加载并解密 kubeconfig
    var c model.Cluster
    if err := database.DB.Where("name=?", name).First(&c).Error; err != nil {
        return nil, fmt.Errorf("cluster not found: %w", err)
    }
    aesKey := config.Load().AESKey
    kubeconfig, err := encrypt.AESDecryptGCM(c.KubeconfigEnc, aesKey)
    if err != nil { return nil, fmt.Errorf("decrypt kubeconfig: %w", err) }
    if err := k8s.Manager.Add(name, kubeconfig); err != nil { return nil, err }
    return k8s.Manager.Get(name)
}

// SimplePing 调用 APIServer 以确认可用
func SimplePing(cli *kubernetes.Clientset) (map[string]interface{}, error) {
    ver, err := cli.Discovery().ServerVersion()
    if err != nil { return nil, err }
    ctx := context.TODO()
    nodes, _ := cli.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
    ns, _ := cli.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
    return map[string]interface{}{
        "version": ver.GitVersion,
        "nodes":   len(nodes.Items),
        "namespaces": len(ns.Items),
    }, nil
}
