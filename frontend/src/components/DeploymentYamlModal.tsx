import { Modal, Tabs, Input, message } from 'antd';
import { useEffect, useState } from 'react';
import { api } from '@/utils/request';

interface Props { open: boolean; onClose: () => void; cluster: string; namespace: string; name?: string }

export default function DeploymentYamlModal({ open, onClose, cluster, namespace, name }: Props) {
  const [yaml, setYaml] = useState('');
  const isEdit = !!name;
  useEffect(() => { (async () => {
    if (open && isEdit) {
      try { const d = await api(`/deployments/detail?cluster=${cluster}&namespace=${namespace}&name=${name}`); setYaml(yamlStringify(d)); } catch (e:any) { message.error(e.message); }
    } else if (open) {
      setYaml(defaultYaml(namespace));
    }
  })(); }, [open, cluster, namespace, name]);

  const onOk = async () => {
    try {
      const body = JSON.stringify({ cluster, namespace, yaml });
      if (isEdit) await api('/deployments/yaml', { method: 'PUT', body });
      else await api('/deployments/yaml', { method: 'POST', body });
      message.success('已提交'); onClose();
    } catch (e:any) { message.error(e.message); }
  };

  return (
    <Modal width={900} open={open} onCancel={onClose} onOk={onOk} title={isEdit ? `编辑 Deployment: ${name}` : '新建 Deployment'} okText={isEdit ? '更新' : '创建'}>
      <Input.TextArea rows={24} value={yaml} onChange={(e)=> setYaml(e.target.value)} spellCheck={false} style={{ fontFamily: 'monospace' }} />
    </Modal>
  );
}

function defaultYaml(ns: string) {
  return `apiVersion: apps/v1
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
}

function yamlStringify(obj: any) {
  // 简化：直接输出 JSON 风格的 YAML（占位，后续可替换为真正 YAML 库）
  try { return JSON.stringify(obj, null, 2); } catch { return '' }
}

