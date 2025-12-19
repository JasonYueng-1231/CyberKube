package service

import (
    "context"

    networkingv1 "k8s.io/api/networking/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

func ListIngresses(cluster, namespace string) ([]networkingv1.Ingress, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil { return nil, err }
    return list.Items, nil
}

func GetIngress(cluster, namespace, name string) (*networkingv1.Ingress, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateIngressYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var ing networkingv1.Ingress
    if err := yaml.Unmarshal([]byte(yml), &ing); err != nil { return err }
    if ing.Namespace == "" { ing.Namespace = namespace }
    _, err = cli.NetworkingV1().Ingresses(ing.Namespace).Create(context.TODO(), &ing, metav1.CreateOptions{})
    return err
}

func UpdateIngressYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var ing networkingv1.Ingress
    if err := yaml.Unmarshal([]byte(yml), &ing); err != nil { return err }
    if ing.Namespace == "" { ing.Namespace = namespace }
    cur, err := cli.NetworkingV1().Ingresses(ing.Namespace).Get(context.TODO(), ing.Name, metav1.GetOptions{})
    if err != nil { return err }
    ing.ResourceVersion = cur.ResourceVersion
    _, err = cli.NetworkingV1().Ingresses(ing.Namespace).Update(context.TODO(), &ing, metav1.UpdateOptions{})
    return err
}

func DeleteIngress(cluster, namespace, name string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    return cli.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
