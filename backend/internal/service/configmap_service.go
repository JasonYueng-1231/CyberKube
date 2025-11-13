package service

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

func ListConfigMaps(cluster, namespace string) ([]corev1.ConfigMap, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil { return nil, err }
    return list.Items, nil
}

func GetConfigMap(cluster, namespace, name string) (*corev1.ConfigMap, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateConfigMapYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var cm corev1.ConfigMap
    if err := yaml.Unmarshal([]byte(yml), &cm); err != nil { return err }
    if cm.Namespace == "" { cm.Namespace = namespace }
    _, err = cli.CoreV1().ConfigMaps(cm.Namespace).Create(context.TODO(), &cm, metav1.CreateOptions{})
    return err
}

func UpdateConfigMapYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var cm corev1.ConfigMap
    if err := yaml.Unmarshal([]byte(yml), &cm); err != nil { return err }
    if cm.Namespace == "" { cm.Namespace = namespace }
    cur, err := cli.CoreV1().ConfigMaps(cm.Namespace).Get(context.TODO(), cm.Name, metav1.GetOptions{})
    if err != nil { return err }
    cm.ResourceVersion = cur.ResourceVersion
    _, err = cli.CoreV1().ConfigMaps(cm.Namespace).Update(context.TODO(), &cm, metav1.UpdateOptions{})
    return err
}

func DeleteConfigMap(cluster, namespace, name string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    return cli.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

