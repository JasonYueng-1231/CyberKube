package service

import (
    "errors"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/config"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/database"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/k8s"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/model"
    "github.com/JasonYueng-1231/CyberKube/backend/pkg/encrypt"
)

type ClusterCreateReq struct {
    Name       string `json:"name"`
    Alias      string `json:"alias"`
    Kubeconfig string `json:"kubeconfig"`
    APIServer  string `json:"api_server"`
    Description string `json:"description"`
}

func CreateCluster(req ClusterCreateReq) (*model.Cluster, error) {
    if req.Name == "" || req.Kubeconfig == "" { return nil, errors.New("名称与kubeconfig必填") }
    aesKey := config.Load().AESKey
    enc, err := encrypt.AESEncryptGCM(req.Kubeconfig, aesKey)
    if err != nil { return nil, err }
    c := &model.Cluster{Name: req.Name, Alias: req.Alias, KubeconfigEnc: enc, APIServer: req.APIServer, Description: req.Description, Status: 1}
    if err := database.DB.Create(c).Error; err != nil { return nil, err }
    // 试图加入内存客户端（不影响创建）
    _ = k8s.Manager.Add(req.Name, req.Kubeconfig)
    return c, nil
}

func ListClusters() ([]model.Cluster, error) {
    var list []model.Cluster
    if err := database.DB.Order("id desc").Find(&list).Error; err != nil { return nil, err }
    return list, nil
}

func GetClusterByName(name string) (*model.Cluster, error) {
    var c model.Cluster
    if err := database.DB.Where("name=?", name).First(&c).Error; err != nil { return nil, err }
    return &c, nil
}

func DeleteCluster(name string) error {
    return database.DB.Where("name=?", name).Delete(&model.Cluster{}).Error
}

