package service

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

func ListServices(cluster, namespace string) ([]corev1.Service, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil { return nil, err }
    return list.Items, nil
}

func GetService(cluster, namespace, name string) (*corev1.Service, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateServiceYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var s corev1.Service
    if err := yaml.Unmarshal([]byte(yml), &s); err != nil { return err }
    if s.Namespace == "" { s.Namespace = namespace }
    _, err = cli.CoreV1().Services(s.Namespace).Create(context.TODO(), &s, metav1.CreateOptions{})
    return err
}

func UpdateServiceYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var s corev1.Service
    if err := yaml.Unmarshal([]byte(yml), &s); err != nil { return err }
    if s.Namespace == "" { s.Namespace = namespace }
    cur, err := cli.CoreV1().Services(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{})
    if err != nil { return err }
    s.ResourceVersion = cur.ResourceVersion
    _, err = cli.CoreV1().Services(s.Namespace).Update(context.TODO(), &s, metav1.UpdateOptions{})
    return err
}

func DeleteService(cluster, namespace, name string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    return cli.CoreV1().Services(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

