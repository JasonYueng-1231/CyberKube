import { Modal, message } from 'antd';
import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import YamlEditor from './YamlEditor';

interface Props { open: boolean; onClose: () => void; cluster: string; namespace: string; name?: string; kind: 'service'|'configmap'|'secret'|'deployment'|'ingress' }

export default function GenericYamlModal({ open, onClose, cluster, namespace, name, kind }: Props) {
  const [yaml, setYaml] = useState('');
  const isEdit = !!name;

  const base = kind === 'service'
    ? 'services'
    : kind === 'configmap'
    ? 'configmaps'
    : kind === 'secret'
    ? 'secrets'
    : kind === 'ingress'
    ? 'ingresses'
    : 'deployments';

  useEffect(() => { (async () => {
    if (!open) return;
    if (isEdit) {
      try { const d = await api(`/${base}/detail?cluster=${cluster}&namespace=${namespace}&name=${name}`); setYaml(JSON.stringify(d, null, 2)); } catch (e:any) { message.error(e.message); }
    } else {
      setYaml(sample(kind, namespace));
    }
  })(); }, [open, cluster, namespace, name, kind]);

  const onOk = async () => {
    try {
      const body = JSON.stringify({ cluster, namespace, yaml });
      if (isEdit) await api(`/${base}/yaml`, { method: 'PUT', body });
      else await api(`/${base}/yaml`, { method: 'POST', body });
      message.success('已提交'); onClose();
    } catch (e:any) { message.error(e.message); }
  };

  return (
    <Modal width={900} open={open} onCancel={onClose} onOk={onOk} title={`${isEdit?'编辑':'新建'} ${kind.toUpperCase()}`} okText={isEdit ? '更新' : '创建'}>
      <YamlEditor value={yaml} onChange={setYaml} />
    </Modal>
  );
}

function sample(kind: string, ns: string) {
  if (kind === 'deployment') return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-deployment
  namespace: ${ns}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: demo
  template:
    metadata:
      labels:
        app: demo
    spec:
      containers:
        - name: nginx
          image: nginx:1.25-alpine
          ports:
            - containerPort: 80
`;
  if (kind === 'service') return `apiVersion: v1
kind: Service
metadata:
  name: demo-svc
  namespace: ${ns}
spec:
  type: ClusterIP
  selector:
    app: demo
  ports:
    - name: http
      port: 80
      targetPort: 80
`;
  if (kind === 'configmap') return `apiVersion: v1
kind: ConfigMap
metadata:
  name: demo-cm
  namespace: ${ns}
data:
  key: value
`;
  if (kind === 'ingress') return `apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: demo-ing
  namespace: ${ns}
spec:
  rules:
    - host: demo.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: demo-svc
                port:
                  number: 80
`;
  return `apiVersion: v1
kind: Secret
metadata:
  name: demo-secret
  namespace: ${ns}
type: Opaque
stringData:
  username: admin
  password: pass
`;
}
