import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import { Card, Select, Button, Table, message, Modal } from 'antd';
import GenericYamlModal from '@/components/GenericYamlModal';

export default function ServiceList() {
  const [clusters, setClusters] = useState<any[]>([]);
  const [cluster, setCluster] = useState('');
  const [namespaces, setNamespaces] = useState<string[]>(['default']);
  const [namespace, setNamespace] = useState('default');
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const [editName, setEditName] = useState<string | undefined>(undefined);

  useEffect(() => { (async () => {
    const d = await api('/clusters'); setClusters(d.items||[]); if (d.items?.length) setCluster(d.items[0].name);
  })(); }, []);
  useEffect(() => { (async () => { if (!cluster) return; const d = await api(`/namespaces?cluster=${cluster}`); setNamespaces(d.items||['default']); if (!d.items?.includes(namespace)) setNamespace('default'); })(); }, [cluster]);

  const load = async () => { if (!cluster) return; setLoading(true); try { const d = await api(`/services?cluster=${cluster}&namespace=${namespace}`); setItems(d.items||[]); } catch (e:any) { message.error(e.message); } finally { setLoading(false); } };
  useEffect(() => { load(); }, [cluster, namespace]);

  const del = async (name:string) => { Modal.confirm({ title:`删除 ${name}?`, onOk: async ()=>{ await api(`/services?cluster=${cluster}&namespace=${namespace}&name=${name}`, { method:'DELETE' }); message.success('已删除'); load(); } }); };

  return (
    <Card className="cyber-card" title="服务列表" extra={<Button type="primary" onClick={()=> { setEditName(undefined); setOpen(true); }}>新建 Service</Button>}>
      <div style={{ display:'flex', gap:8, marginBottom:12 }}>
        <Select value={cluster} onChange={setCluster} options={(clusters||[]).map((c:any)=>({label:c.name, value:c.name}))} style={{ minWidth:200 }} />
        <Select value={namespace} onChange={setNamespace} options={(namespaces||[]).map((n)=>({label:n, value:n}))} style={{ minWidth:200 }} />
        <Button onClick={load}>刷新</Button>
      </div>
      <Table rowKey={(r:any)=> r.metadata?.namespace+"/"+r.metadata?.name} loading={loading} dataSource={items} pagination={false}
        columns={[
          { title:'名称', dataIndex:['metadata','name'] },
          { title:'命名空间', dataIndex:['metadata','namespace'] },
          { title:'类型', dataIndex:['spec','type'] },
          { title:'ClusterIP', dataIndex:['spec','clusterIP'] },
          { title:'端口', render: (_:any, r:any)=> (r.spec?.ports||[]).map((p:any)=> `${p.port}:${p.targetPort}`).join(', ') },
          { title:'操作', render: (_:any,r:any)=> <>
            <Button size="small" style={{marginRight:8}} onClick={()=> { setEditName(r.metadata.name); setOpen(true); }}>编辑(YAML)</Button>
            <Button size="small" danger onClick={()=> del(r.metadata.name)}>删除</Button>
          </> }
        ]}
      />
      <GenericYamlModal open={open} onClose={()=> { setOpen(false); load(); }} cluster={cluster} namespace={namespace} name={editName} kind='service' />
    </Card>
  );
}

