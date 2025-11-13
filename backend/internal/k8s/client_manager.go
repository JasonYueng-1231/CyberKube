package k8s

import (
    "fmt"
    "sync"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

type ClusterManager struct {
    mu sync.RWMutex
    clients map[string]*kubernetes.Clientset
}

var Manager = &ClusterManager{clients: map[string]*kubernetes.Clientset{}}

func (m *ClusterManager) Add(name, kubeconfig string) error {
    cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
    if err != nil { return err }
    cli, err := kubernetes.NewForConfig(cfg)
    if err != nil { return err }
    m.mu.Lock(); defer m.mu.Unlock()
    m.clients[name] = cli
    return nil
}

func (m *ClusterManager) Get(name string) (*kubernetes.Clientset, error) {
    m.mu.RLock(); defer m.mu.RUnlock()
    if c, ok := m.clients[name]; ok { return c, nil }
    return nil, fmt.Errorf("client not found: %s", name)
}

