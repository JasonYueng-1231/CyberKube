package k8s

import (
    "fmt"
    "sync"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/config"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

type ClusterManager struct {
    mu sync.RWMutex
    clients map[string]*kubernetes.Clientset
    configs map[string]*rest.Config
}

var Manager = &ClusterManager{clients: map[string]*kubernetes.Clientset{}, configs: map[string]*rest.Config{}}

func (m *ClusterManager) Add(name, kubeconfig string) error {
    cfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
    if err != nil { return err }
    // TLS 验证可选关闭（用于排障/内网证书）
    if config.Load().SkipTLSVerify {
        cfg.TLSClientConfig.Insecure = true
    }
    cfg.QPS = 50
    cfg.Burst = 100
    cli, err := kubernetes.NewForConfig(cfg)
    if err != nil { return err }
    m.mu.Lock(); defer m.mu.Unlock()
    m.clients[name] = cli
    m.configs[name] = cfg
    return nil
}

func (m *ClusterManager) Get(name string) (*kubernetes.Clientset, error) {
    m.mu.RLock(); defer m.mu.RUnlock()
    if c, ok := m.clients[name]; ok { return c, nil }
    return nil, fmt.Errorf("client not found: %s", name)
}

func (m *ClusterManager) GetConfig(name string) (*rest.Config, error) {
    m.mu.RLock(); defer m.mu.RUnlock()
    if c, ok := m.configs[name]; ok { return c, nil }
    return nil, fmt.Errorf("config not found: %s", name)
}
