import { Card, Row, Col, Button, Select, message } from 'antd';
import { RocketOutlined, ApiOutlined, AppstoreAddOutlined, CodeOutlined } from '@ant-design/icons';
import StatsCard from '@/shared/StatsCard/StatsCard';
import '@/shared/StatsCard/StatsCard.css';
import { useEffect, useState } from 'react';
import { api } from '@/utils/request';

export default function Dashboard() {
  const [clusters, setClusters] = useState<any[]>([]);
  const [cluster, setCluster] = useState<string>('');
  const [overview, setOverview] = useState<any>({ cpu_percent: 62, mem_percent: 48 });

  useEffect(() => { (async () => {
    try {
      const data = await api('/clusters');
      setClusters(data.items || []);
      const c = (data.items?.[0]?.name) || '';
      setCluster(c);
    } catch (e:any) { message.error(e.message); }
  })(); }, []);

  useEffect(() => { (async () => {
    if (!cluster) return;
    try {
      const data = await api(`/metrics/overview?cluster=${cluster}`);
      if (data.cpu_percent >= 0) setOverview((o:any)=> ({...o, cpu_percent: data.cpu_percent}));
      if (data.mem_percent >= 0) setOverview((o:any)=> ({...o, mem_percent: data.mem_percent}));
    } catch {}
  })(); }, [cluster]);
  return (
    <div>
      <div style={{ display:'flex', justifyContent:'flex-end', marginBottom: 12 }}>
        <Select value={cluster} onChange={setCluster} style={{ minWidth: 200 }} placeholder="选择集群"
          options={(clusters||[]).map((c:any)=> ({label:c.name, value:c.name}))}
        />
      </div>
      <Row gutter={[16, 16]}>
        <Col xs={24} md={12} lg={6}>
          <StatsCard color="cyan" icon={<ApiOutlined />} value={8} label="节点" status="Ready 8 / NotReady 0" />
        </Col>
        <Col xs={24} md={12} lg={6}>
          <StatsCard color="purple" icon={<RocketOutlined />} value={42} label="容器组" status="Running 39 / Pending 3" />
        </Col>
        <Col xs={24} md={12} lg={6}>
          <StatsCard color="green" icon={<AppstoreAddOutlined />} value={12} label="命名空间" status="活跃 9 / 冻结 3" />
        </Col>
        <Col xs={24} md={12} lg={6}>
          <StatsCard color="cyan" icon={<CodeOutlined />} value={26} label="服务" status="LB 4 / NodePort 6" />
        </Col>
      </Row>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <Card className="cyber-card" title="CPU 使用率">
            <div className="progress-bar">
              <div className="progress progress-requests" style={{ width: `${Math.min(100, Math.max(0, Math.round(overview.cpu_percent||0)))}%` }} />
              <div className="progress-label">使用率 {Math.round(overview.cpu_percent||0)}%</div>
            </div>
            <div className="progress-bar">
              <div className="progress progress-limits" style={{ width: '200%' }} />
              <div className="progress-label">Limits 200%</div>
            </div>
          </Card>
          <Card className="cyber-card" title="最近事件" style={{ marginTop: 16 }}>
            <ul className="event-list">
              <li className="ok">Scaled up replica set nginx to 3</li>
              <li className="warn">OOMKilled on pod api-7d9f...</li>
              <li className="ok">Created service lb/payment</li>
              <li className="error">Probe failed on pod cache-0</li>
            </ul>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card className="cyber-card" title="内存 使用率">
            <div className="progress-bar">
              <div className="progress progress-requests" style={{ width: `${Math.min(100, Math.max(0, Math.round(overview.mem_percent||0)))}%` }} />
              <div className="progress-label">使用率 {Math.round(overview.mem_percent||0)}%</div>
            </div>
            <div className="progress-bar">
              <div className="progress progress-limits" style={{ width: '120%' }} />
              <div className="progress-label">Limits 120%</div>
            </div>
          </Card>
          <Card className="cyber-card" title="快速操作" style={{ marginTop: 16 }}>
            <Button type="primary" style={{ marginRight: 8 }}>新建 Deployment</Button>
            <Button style={{ marginRight: 8 }}>扩容</Button>
            <Button style={{ marginRight: 8 }}>重启</Button>
            <Button danger>YAML 编辑</Button>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
