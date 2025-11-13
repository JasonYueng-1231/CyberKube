package service

import (
    "context"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

func ListSecrets(cluster, namespace string) ([]corev1.Secret, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    list, err := cli.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil { return nil, err }
    return list.Items, nil
}

func GetSecret(cluster, namespace, name string) (*corev1.Secret, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func CreateSecretYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var s corev1.Secret
    if err := yaml.Unmarshal([]byte(yml), &s); err != nil { return err }
    if s.Namespace == "" { s.Namespace = namespace }
    _, err = cli.CoreV1().Secrets(s.Namespace).Create(context.TODO(), &s, metav1.CreateOptions{})
    return err
}

func UpdateSecretYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var s corev1.Secret
    if err := yaml.Unmarshal([]byte(yml), &s); err != nil { return err }
    if s.Namespace == "" { s.Namespace = namespace }
    cur, err := cli.CoreV1().Secrets(s.Namespace).Get(context.TODO(), s.Name, metav1.GetOptions{})
    if err != nil { return err }
    s.ResourceVersion = cur.ResourceVersion
    _, err = cli.CoreV1().Secrets(s.Namespace).Update(context.TODO(), &s, metav1.UpdateOptions{})
    return err
}

func DeleteSecret(cluster, namespace, name string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    return cli.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

