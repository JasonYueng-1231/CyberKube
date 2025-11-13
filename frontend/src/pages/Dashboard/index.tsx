import { Card, Row, Col, Button } from 'antd';
import { RocketOutlined, ApiOutlined, AppstoreAddOutlined, CodeOutlined } from '@ant-design/icons';
import StatsCard from '@/shared/StatsCard/StatsCard';
import '@/shared/StatsCard/StatsCard.css';

export default function Dashboard() {
  return (
    <div>
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
              <div className="progress progress-requests" style={{ width: '62%' }} />
              <div className="progress-label">Requests 62%</div>
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
              <div className="progress progress-requests" style={{ width: '48%' }} />
              <div className="progress-label">Requests 48%</div>
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

