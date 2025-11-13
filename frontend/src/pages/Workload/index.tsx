import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import { Card, Tabs, Select, Input, Button, Table, Modal, Form, message } from 'antd';

export default function Workloads() {
  const [clusters, setClusters] = useState<any[]>([]);
  const [cluster, setCluster] = useState<string>('');
  const [namespace, setNamespace] = useState<string>('default');

  const [deps, setDeps] = useState<any[]>([]);
  const [pods, setPods] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [logOpen, setLogOpen] = useState(false);
  const [logs, setLogs] = useState('');
  const [form] = Form.useForm();

  useEffect(() => { (async () => {
    const data = await api('/clusters');
    setClusters(data.items || []);
    if (!cluster && data.items?.length) setCluster(data.items[0].name);
  })(); }, []);

  const loadDeps = async () => {
    if (!cluster) return; setLoading(true);
    try { const data = await api(`/deployments?cluster=${cluster}&namespace=${namespace}`); setDeps(data.items || []); } finally { setLoading(false); }
  };
  const loadPods = async () => {
    if (!cluster) return; setLoading(true);
    try { const data = await api(`/pods?cluster=${cluster}&namespace=${namespace}`); setPods(data.items || []); } finally { setLoading(false); }
  };

  useEffect(() => { loadDeps(); loadPods(); }, [cluster, namespace]);

  const scale = async (name: string, replicas: number) => {
    try {
      await api('/deployments/scale', { method:'POST', body: JSON.stringify({ cluster, namespace, name, replicas }) });
      message.success('伸缩已提交'); loadDeps();
    } catch(e:any) { message.error(e.message); }
  };
  const restart = async (name: string) => {
    try {
      await api('/deployments/restart', { method:'POST', body: JSON.stringify({ cluster, namespace, name }) });
      message.success('重启已提交');
    } catch(e:any) { message.error(e.message); }
  };
  const viewLogs = async (name: string) => {
    setLogOpen(true); setLogs('加载中...');
    try {
      const res = await fetch(`/api/v1/pods/logs?cluster=${cluster}&namespace=${namespace}&name=${encodeURIComponent(name)}&tail=200`, { headers: { Authorization: 'Bearer '+ localStorage.getItem('token') } });
      const txt = await res.text(); setLogs(txt);
    } catch { setLogs('获取日志失败'); }
  };

  return (
    <Card className="cyber-card" title="工作负载">
      <div style={{ marginBottom: 12, display:'flex', gap:8 }}>
        <Select value={cluster} onChange={setCluster} style={{ minWidth: 200 }} placeholder="选择集群">
          {clusters.map((c:any)=> <Select.Option key={c.name} value={c.name}>{c.name}</Select.Option>)}
        </Select>
        <Input value={namespace} onChange={(e)=> setNamespace(e.target.value)} style={{ width: 200 }} placeholder="命名空间，如 default"/>
        <Button onClick={()=>{loadDeps();loadPods();}}>刷新</Button>
      </div>

      <Tabs items={[
        { key: 'dep', label: 'Deployment', children: (
          <Table rowKey={(r:any)=> r.namespace+"/"+r.name} dataSource={deps} loading={loading}
            columns={[
              { title:'名称', dataIndex:'name' },
              { title:'副本', render: (_:any,r:any)=> `${r.available}/${r.replicas}` },
              { title:'已更新', dataIndex:'updated' },
              { title:'操作', render: (_:any, r:any)=> <>
                <Button size="small" onClick={()=> scale(r.name, r.replicas+1)} style={{marginRight:8}}>扩容</Button>
                <Button size="small" onClick={()=> scale(r.name, Math.max(0, r.replicas-1))} style={{marginRight:8}}>缩容</Button>
                <Button size="small" onClick={()=> restart(r.name)}>重启</Button>
              </> }
            ]}
          />
        ) },
        { key: 'pod', label: 'Pod', children: (
          <Table rowKey={(r:any)=> r.namespace+"/"+r.name} dataSource={pods} loading={loading}
            columns={[
              { title:'名称', dataIndex:'name' },
              { title:'状态', dataIndex:'phase' },
              { title:'就绪', dataIndex:'ready' },
              { title:'重启', dataIndex:'restarts' },
              { title:'节点', dataIndex:'node_name' },
              { title:'操作', render: (_:any, r:any)=> <Button size="small" onClick={()=> viewLogs(r.name)}>日志</Button> },
            ]}
          />
        ) }
      ]} />

      <Modal width={900} title="Pod 日志" open={logOpen} onCancel={()=> setLogOpen(false)} footer={null}>
        <pre style={{ maxHeight: 500, overflow:'auto' }}>{logs}</pre>
      </Modal>
    </Card>
  );
}

