import { useEffect, useMemo, useState } from 'react';
import { Descriptions, Modal, Tabs, Tag, Typography, message } from 'antd';
import YamlEditor from './YamlEditor';
import { api } from '@/utils/request';

const { Text } = Typography;

interface Props {
  open: boolean;
  onClose: () => void;
  cluster: string;
  namespace: string;
  pod: string;
}

export default function PodDetailModal({ open, onClose, cluster, namespace, pod }: Props) {
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    if (!open || !pod) return;
    (async () => {
      setLoading(true);
      try {
        const d = await api(`/pods/detail?cluster=${cluster}&namespace=${namespace}&name=${pod}`);
        setData(d);
      } catch (e: any) {
        message.error(e.message || '获取 Pod 详情失败');
      } finally {
        setLoading(false);
      }
    })();
  }, [open, cluster, namespace, pod]);

  const containers = useMemo(() => data?.pod?.spec?.containers || [], [data]);
  const initContainers = useMemo(() => data?.pod?.spec?.initContainers || [], [data]);

  return (
    <Modal width={960} open={open} onCancel={onClose} title={`Pod 详情: ${pod}`} footer={null}>
      <Tabs
        items={[
          {
            key: 'info',
            label: '基本信息',
            children: (
              <Descriptions column={2} bordered size="small" labelStyle={{ width: 140 }}>
                <Descriptions.Item label="名称">{data?.pod?.metadata?.name}</Descriptions.Item>
                <Descriptions.Item label="命名空间">{data?.pod?.metadata?.namespace}</Descriptions.Item>
                <Descriptions.Item label="节点">{data?.pod?.spec?.nodeName}</Descriptions.Item>
                <Descriptions.Item label="Pod IP">{data?.pod?.status?.podIP || '-'}</Descriptions.Item>
                <Descriptions.Item label="状态">
                  <Tag color={data?.pod?.status?.phase === 'Running' ? 'green' : 'orange'}>
                    {data?.pod?.status?.phase || '-'}
                  </Tag>
                </Descriptions.Item>
                <Descriptions.Item label="重启次数">
                  {data?.pod?.status?.containerStatuses?.reduce((s: number, cs: any) => s + (cs?.restartCount || 0), 0) || 0}
                </Descriptions.Item>
                <Descriptions.Item label="创建时间" span={2}>
                  {data?.pod?.metadata?.creationTimestamp || '-'}
                </Descriptions.Item>
              </Descriptions>
            ),
          },
          {
            key: 'containers',
            label: `容器 (${containers.length})`,
            children: (
              <div style={{ display: 'grid', gap: 12 }}>
                {containers.map((c: any) => (
                  <div key={c.name} className="cyber-card" style={{ padding: 12 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
                      <Text strong>{c.name}</Text>
                      <Tag>{c.image}</Tag>
                    </div>
                    <div style={{ color: 'var(--text-secondary)' }}>
                      端口: {(c.ports || []).map((p: any) => `${p.containerPort}/${p.protocol || 'TCP'}`).join(', ') || '无'}
                    </div>
                    <div style={{ color: 'var(--text-secondary)' }}>
                      环境变量: {(c.env || []).length} 项
                    </div>
                  </div>
                ))}
                {containers.length === 0 && <Text type="secondary">无容器</Text>}
              </div>
            ),
          },
          {
            key: 'init',
            label: `Init 容器 (${initContainers.length})`,
            children: (
              <div style={{ display: 'grid', gap: 12 }}>
                {initContainers.map((c: any) => (
                  <div key={c.name} className="cyber-card" style={{ padding: 12 }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 6 }}>
                      <Text strong>{c.name}</Text>
                      <Tag>{c.image}</Tag>
                    </div>
                    <div style={{ color: 'var(--text-secondary)' }}>
                      命令: {(c.command || []).join(' ') || '无'}
                    </div>
                  </div>
                ))}
                {initContainers.length === 0 && <Text type="secondary">无 Init 容器</Text>}
              </div>
            ),
          },
          {
            key: 'events',
            label: `事件 (${data?.events?.length || 0})`,
            children: (
              <div style={{ maxHeight: 320, overflow: 'auto' }}>
                {(data?.events || []).map((e: any, idx: number) => (
                  <div key={idx} style={{ padding: '8px 0', borderBottom: '1px dashed var(--border-normal)' }}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 4 }}>
                      <Text>{e.message}</Text>
                      <Tag color={e.type === 'Warning' ? 'orange' : 'blue'}>{e.type}</Tag>
                    </div>
                    <div style={{ color: 'var(--text-secondary)', fontSize: 12 }}>
                      {e.reason} · {e.firstTimestamp || e.eventTime || '-'}
                    </div>
                  </div>
                ))}
                {(!data?.events || data.events.length === 0) && <Text type="secondary">暂无事件</Text>}
              </div>
            ),
          },
          {
            key: 'yaml',
            label: 'YAML',
            children: <YamlEditor value={data?.yaml || ''} onChange={() => {}} readOnly height={360} />,
          },
        ]}
        loading={loading}
      />
    </Modal>
  );
}
