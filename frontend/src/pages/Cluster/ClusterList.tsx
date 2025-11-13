import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import { Button, Card, Modal, Form, Input, Table, message } from 'antd';

export default function ClusterList() {
  const [loading, setLoading] = useState(false);
  const [items, setItems] = useState<any[]>([]);
  const [open, setOpen] = useState(false);
  const [form] = Form.useForm();

  const load = async () => {
    setLoading(true);
    try {
      const data = await api('/clusters');
      setItems(data.items || []);
    } finally { setLoading(false); }
  };

  useEffect(() => { load(); }, []);

  const create = async () => {
    const v = await form.validateFields();
    try {
      await api('/clusters', { method: 'POST', body: JSON.stringify(v) });
      message.success('创建成功'); setOpen(false); form.resetFields(); load();
    } catch (e: any) { message.error(e.message); }
  };

  const remove = async (name: string) => {
    Modal.confirm({ title: `删除集群 ${name}?`, onOk: async () => {
      await api(`/clusters/${name}`, { method: 'DELETE' });
      message.success('已删除'); load();
    }});
  };

  return (
    <Card className="cyber-card" title="集群管理" extra={<Button type="primary" onClick={() => setOpen(true)}>新增集群</Button>}>
      <Table rowKey="id" dataSource={items} loading={loading} pagination={false}
        columns={[
          { title: '名称', dataIndex: 'name' },
          { title: '别名', dataIndex: 'alias' },
          { title: 'API Server', dataIndex: 'api_server' },
          { title: '版本', dataIndex: 'version' },
          { title: '状态', dataIndex: 'status', render: (v) => (v === 1 ? '正常' : '异常') },
          { title: '操作', render: (_:any, r:any) => <Button danger onClick={() => remove(r.name)}>删除</Button> },
        ]}
      />

      <Modal title="新增集群" open={open} onCancel={() => setOpen(false)} onOk={create} okText="创建">
        <Form form={form} layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}>
            <Input placeholder="唯一英文名，如 prod-cluster" />
          </Form.Item>
          <Form.Item name="alias" label="别名">
            <Input placeholder="便于识别的中文名" />
          </Form.Item>
          <Form.Item name="api_server" label="API Server">
            <Input placeholder="例如 10.0.0.100:6443" />
          </Form.Item>
          <Form.Item name="kubeconfig" label="kubeconfig" rules={[{ required: true }]}>
            <Input.TextArea rows={6} placeholder="粘贴 kubeconfig 内容" />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
}

