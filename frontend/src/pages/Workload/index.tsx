import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import { Card, Tabs, Select, Button, Table, Modal, Form, message } from 'antd';
import LogStreamModal from '@/components/LogStreamModal';
import TerminalModal from '@/components/TerminalModal';
import DeploymentYamlModal from '@/components/DeploymentYamlModal';

export default function Workloads() {
  const [clusters, setClusters] = useState<any[]>([]);
  const [cluster, setCluster] = useState<string>('');
  const [namespace, setNamespace] = useState<string>('default');
  const [namespaces, setNamespaces] = useState<string[]>(['default']);

  const [deps, setDeps] = useState<any[]>([]);
  const [pods, setPods] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [logOpen, setLogOpen] = useState(false);
  const [logs, setLogs] = useState('');
  const [form] = Form.useForm();
  const [logCtx, setLogCtx] = useState<any>(null);
  const [termOpen, setTermOpen] = useState(false);
  const [termCtx, setTermCtx] = useState<any>(null);
  const [depYamlOpen, setDepYamlOpen] = useState(false);
  const [depEdit, setDepEdit] = useState<string | undefined>(undefined);

  useEffect(() => { (async () => {
    const data = await api('/clusters');
    setClusters(data.items || []);
    if (!cluster && data.items?.length) setCluster(data.items[0].name);
  })(); }, []);

  // 加载命名空间
  useEffect(() => { (async () => {
    if (!cluster) return;
    try { const data = await api(`/namespaces?cluster=${cluster}`); setNamespaces(data.items || ['default']); if (!data.items?.includes(namespace)) setNamespace('default'); } catch {}
  })(); }, [cluster]);

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
  const viewLogs = (r: any) => {
    setLogCtx({ pod: r.name, containers: r.containers });
    setLogOpen(true);
  };
  const openShell = (r: any) => {
    setTermCtx({ pod: r.name, containers: r.containers });
    setTermOpen(true);
  };

  return (
    <Card className="cyber-card" title="工作负载">
      <div style={{ marginBottom: 12, display:'flex', gap:8, alignItems:'center' }}>
        <Select value={cluster} onChange={setCluster} style={{ minWidth: 200 }} placeholder="选择集群"
          options={(clusters||[]).map((c:any)=> ({label:c.name, value:c.name}))}
        />
        <Select value={namespace} onChange={setNamespace} style={{ minWidth: 200 }} placeholder="命名空间"
          options={(namespaces||[]).map(n=> ({label:n, value:n}))}
        />
        <Button onClick={()=>{loadDeps();loadPods();}}>刷新</Button>
        <Button type="primary" onClick={()=> { setDepEdit(undefined); setDepYamlOpen(true); }}>新建 Deployment</Button>
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
                <Button size="small" onClick={()=> restart(r.name)} style={{marginRight:8}}>重启</Button>
                <Button size="small" onClick={()=> { setDepEdit(r.name); setDepYamlOpen(true); }}>编辑(YAML)</Button>
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
              { title:'操作', render: (_:any, r:any)=> <>
                <Button size="small" onClick={()=> viewLogs(r)} style={{marginRight:8}}>日志(流)</Button>
                <Button size="small" onClick={()=> openShell(r)}>Shell</Button>
              </> },
            ]}
          />
        ) }
      ]} />

      <LogStreamModal open={logOpen} onClose={()=> setLogOpen(false)} cluster={cluster} namespace={namespace} pod={logCtx?.pod} containers={logCtx?.containers} />
      <TerminalModal open={termOpen} onClose={()=> setTermOpen(false)} cluster={cluster} namespace={namespace} pod={termCtx?.pod} containers={termCtx?.containers} />
      <DeploymentYamlModal open={depYamlOpen} onClose={()=> { setDepYamlOpen(false); loadDeps(); }} cluster={cluster} namespace={namespace} name={depEdit} />
    </Card>
  );
}
