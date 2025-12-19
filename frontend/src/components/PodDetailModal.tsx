<<<<<<< HEAD
import { useEffect, useMemo, useState } from 'react';
import { Descriptions, Modal, Tabs, Tag, Typography, message } from 'antd';
import { api } from '@/utils/request';
import YamlEditor from './YamlEditor';

const { Text } = Typography;
=======
import { Modal, Tabs, Descriptions, Table, Tag, Spin, Empty, Button, message } from 'antd';
import { useEffect, useState } from 'react';
import { api } from '@/utils/request';
import YamlEditor from './YamlEditor';
import LogStreamModal from './LogStreamModal';
import TerminalModal from './TerminalModal';
>>>>>>> origin/develop

interface Props {
  open: boolean;
  onClose: () => void;
  cluster: string;
  namespace: string;
  podName?: string;
}

export default function PodDetailModal({ open, onClose, cluster, namespace, podName }: Props) {
  const [loading, setLoading] = useState(false);
<<<<<<< HEAD
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    if (!open || !podName) return;
    (async () => {
      setLoading(true);
      try {
        const d = await api(
=======
  const [detail, setDetail] = useState<any>(null);
  const [logOpen, setLogOpen] = useState(false);
  const [logCtx, setLogCtx] = useState<any>(null);
  const [termOpen, setTermOpen] = useState(false);
  const [termCtx, setTermCtx] = useState<any>(null);

  useEffect(() => {
    (async () => {
      if (!open || !cluster || !namespace || !podName) return;
      setLoading(true);
      try {
        const data = await api(
>>>>>>> origin/develop
          `/pods/detail?cluster=${encodeURIComponent(cluster)}&namespace=${encodeURIComponent(
            namespace,
          )}&name=${encodeURIComponent(podName)}`,
        );
<<<<<<< HEAD
        setData(d);
=======
        setDetail(data);
>>>>>>> origin/develop
      } catch (e: any) {
        message.error(e.message || '获取 Pod 详情失败');
      } finally {
        setLoading(false);
      }
    })();
  }, [open, cluster, namespace, podName]);

<<<<<<< HEAD
  const containers = useMemo(() => data?.pod?.spec?.containers || [], [data]);
  const initContainers = useMemo(() => data?.pod?.spec?.initContainers || [], [data]);

  return (
    <Modal width={960} open={open} onCancel={onClose} title={`Pod 详情: ${podName || ''}`} footer={null}>
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
=======
  const pod = detail?.pod;
  const events = detail?.events || [];
  const yaml = detail?.yaml || '';

  const containers: string[] = (pod?.spec?.containers || []).map((c: any) => c.name).filter(Boolean);

  const openLogs = () => {
    if (!pod) return;
    setLogCtx({ pod: pod.metadata?.name, containers });
    setLogOpen(true);
  };

  const openShell = () => {
    if (!pod) return;
    setTermCtx({ pod: pod.metadata?.name, containers });
    setTermOpen(true);
  };

  const renderLabels = (obj: Record<string, string> | undefined) => {
    if (!obj || !Object.keys(obj).length) return '-';
    return (
      <>
        {Object.entries(obj).map(([k, v]) => (
          <Tag key={k}>{`${k}=${v}`}</Tag>
        ))}
      </>
    );
  };

  const renderEventTime = (ev: any) => ev.lastTimestamp || ev.eventTime || ev.metadata?.creationTimestamp || '';

  return (
    <>
      <Modal
        width={900}
        open={open}
        onCancel={onClose}
        footer={null}
        title={`Pod 详情: ${podName || ''}`}
      >
        {loading ? (
          <div style={{ textAlign: 'center', padding: 40 }}>
            <Spin />
          </div>
        ) : !pod ? (
          <Empty description="暂无数据" />
        ) : (
          <Tabs
            items={[
              {
                key: 'basic',
                label: '基本信息',
                children: (
                  <>
                    <div style={{ marginBottom: 12 }}>
                      <Button size="small" onClick={openLogs} style={{ marginRight: 8 }} disabled={!containers.length}>
                        日志(流)
                      </Button>
                      <Button size="small" onClick={openShell} disabled={!containers.length}>
                        Shell
                      </Button>
                    </div>
                    <Descriptions column={2} bordered size="small">
                      <Descriptions.Item label="名称">{pod.metadata?.name}</Descriptions.Item>
                      <Descriptions.Item label="命名空间">{pod.metadata?.namespace}</Descriptions.Item>
                      <Descriptions.Item label="状态">{pod.status?.phase}</Descriptions.Item>
                      <Descriptions.Item label="节点">{pod.spec?.nodeName}</Descriptions.Item>
                      <Descriptions.Item label="Pod IP">{pod.status?.podIP}</Descriptions.Item>
                      <Descriptions.Item label="主机 IP">{pod.status?.hostIP}</Descriptions.Item>
                      <Descriptions.Item label="重启次数">
                        {(pod.status?.containerStatuses || []).reduce(
                          (sum: number, cs: any) => sum + (cs.restartCount || 0),
                          0,
                        )}
                      </Descriptions.Item>
                      <Descriptions.Item label="创建时间">
                        {pod.metadata?.creationTimestamp || ''}
                      </Descriptions.Item>
                      <Descriptions.Item label="标签" span={2}>
                        {renderLabels(pod.metadata?.labels)}
                      </Descriptions.Item>
                      <Descriptions.Item label="注解" span={2}>
                        {renderLabels(pod.metadata?.annotations)}
                      </Descriptions.Item>
                    </Descriptions>
                  </>
                ),
              },
              {
                key: 'containers',
                label: '容器',
                children: (
                  <Table
                    rowKey={(r: any) => r.name}
                    dataSource={(pod.spec?.containers || []) as any[]}
                    pagination={false}
                    columns={[
                      { title: '名称', dataIndex: 'name' },
                      { title: '镜像', dataIndex: 'image' },
                      {
                        title: '状态',
                        render: (_: any, r: any) => {
                          const cs =
                            (pod.status?.containerStatuses || []).find((it: any) => it.name === r.name) || {};
                          const ready = cs.ready;
                          const state = cs.state?.waiting?.reason || cs.state?.terminated?.reason || 'Running';
                          return (
                            <>
                              <Tag color={ready ? 'green' : 'red'}>{ready ? '就绪' : '未就绪'}</Tag>
                              <span>{state}</span>
                            </>
                          );
                        },
                      },
                      {
                        title: '重启',
                        render: (_: any, r: any) => {
                          const cs =
                            (pod.status?.containerStatuses || []).find((it: any) => it.name === r.name) || {};
                          return cs.restartCount ?? 0;
                        },
                      },
                      {
                        title: '端口',
                        render: (_: any, r: any) =>
                          (r.ports || [])
                            .map((p: any) => `${p.name || ''}${p.name ? ':' : ''}${p.containerPort}/${p.protocol}`)
                            .join(', ') || '-',
                      },
                    ]}
                  />
                ),
              },
              {
                key: 'events',
                label: '事件',
                children: events.length ? (
                  <Table
                    rowKey={(r: any) => r.metadata?.uid || `${r.metadata?.name}-${renderEventTime(r)}`}
                    dataSource={events as any[]}
                    pagination={false}
                    columns={[
                      { title: '类型', dataIndex: 'type' },
                      { title: '原因', dataIndex: 'reason' },
                      { title: '消息', dataIndex: 'message', ellipsis: true },
                      {
                        title: '时间',
                        render: (_: any, r: any) => renderEventTime(r),
                      },
                      { title: '次数', dataIndex: 'count' },
                    ]}
                  />
                ) : (
                  <Empty description="暂无事件" />
                ),
              },
              {
                key: 'yaml',
                label: 'YAML',
                children: (
                  <YamlEditor
                    value={yaml}
                    onChange={() => {}}
                    readOnly
                    height={420}
                  />
                ),
              },
            ]}
          />
        )}
      </Modal>

      <LogStreamModal
        open={logOpen}
        onClose={() => setLogOpen(false)}
        cluster={cluster}
        namespace={namespace}
        pod={logCtx?.pod}
        containers={logCtx?.containers}
      />
      <TerminalModal
        open={termOpen}
        onClose={() => setTermOpen(false)}
        cluster={cluster}
        namespace={namespace}
        pod={termCtx?.pod}
        containers={termCtx?.containers}
      />
    </>
  );
}

>>>>>>> origin/develop
