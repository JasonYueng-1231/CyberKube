package service

import (
    "context"
    "time"

    appsv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "sigs.k8s.io/yaml"
)

type DeploymentItem struct {
    Name      string `json:"name"`
    Namespace string `json:"namespace"`
    Replicas  int32  `json:"replicas"`
    Available int32  `json:"available"`
    Updated   int32  `json:"updated"`
    Age       string `json:"age"`
}

func ListDeployments(cluster, namespace string) ([]DeploymentItem, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    ctx := context.TODO()
    list, err := cli.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
    if err != nil { return nil, err }
    var out []DeploymentItem
    for _, d := range list.Items {
        out = append(out, DeploymentItem{
            Name: d.Name, Namespace: d.Namespace,
            Replicas: getInt32(d.Spec.Replicas),
            Available: d.Status.AvailableReplicas,
            Updated: d.Status.UpdatedReplicas,
            Age: time.Since(d.CreationTimestamp.Time).Round(time.Minute).String(),
        })
    }
    return out, nil
}

func ScaleDeployment(cluster, namespace, name string, replicas int32) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    ctx := context.TODO()
    dep, err := cli.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil { return err }
    dep.Spec.Replicas = &replicas
    _, err = cli.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
    return err
}

func RestartDeployment(cluster, namespace, name string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    ctx := context.TODO()
    dep, err := cli.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
    if err != nil { return err }
    if dep.Spec.Template.Annotations == nil { dep.Spec.Template.Annotations = map[string]string{} }
    dep.Spec.Template.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
    _, err = cli.AppsV1().Deployments(namespace).Update(ctx, dep, metav1.UpdateOptions{})
    return err
}

func getInt32(p *int32) int32 { if p == nil { return 0 }; return *p }

// CreateDeploymentFromYAML 基于 YAML 创建 Deployment
func CreateDeploymentFromYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var dep appsv1.Deployment
    if err := yaml.Unmarshal([]byte(yml), &dep); err != nil { return err }
    if dep.Namespace == "" { dep.Namespace = namespace }
    _, err = cli.AppsV1().Deployments(dep.Namespace).Create(context.TODO(), &dep, metav1.CreateOptions{})
    return err
}

// UpdateDeploymentFromYAML 基于 YAML 更新 Deployment（保留 ResourceVersion）
func UpdateDeploymentFromYAML(cluster, namespace, yml string) error {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return err }
    var dep appsv1.Deployment
    if err := yaml.Unmarshal([]byte(yml), &dep); err != nil { return err }
    if dep.Namespace == "" { dep.Namespace = namespace }
    cur, err := cli.AppsV1().Deployments(dep.Namespace).Get(context.TODO(), dep.Name, metav1.GetOptions{})
    if err != nil { return err }
    dep.ResourceVersion = cur.ResourceVersion
    _, err = cli.AppsV1().Deployments(dep.Namespace).Update(context.TODO(), &dep, metav1.UpdateOptions{})
    return err
}

func GetDeployment(cluster, namespace, name string) (*appsv1.Deployment, error) {
    cli, err := GetClientForCluster(cluster)
    if err != nil { return nil, err }
    return cli.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
