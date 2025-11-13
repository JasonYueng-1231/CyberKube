import { Layout, Menu, theme } from 'antd';
import {
  AppstoreOutlined,
  ClusterOutlined,
  DeploymentUnitOutlined,
  SettingOutlined,
} from '@ant-design/icons';
import Dashboard from './pages/Dashboard';
import './styles/global.css';

const { Header, Sider, Content } = Layout;

export default function App() {
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={220} theme="dark" style={{ background: 'var(--bg-sidebar)' }}>
        <div className="logo">CYBERKUBE</div>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={["dashboard"]}
          items={[
            { key: 'dashboard', icon: <AppstoreOutlined />, label: '仪表盘' },
            { key: 'clusters', icon: <ClusterOutlined />, label: '集群' },
            { key: 'workloads', icon: <DeploymentUnitOutlined />, label: '工作负载' },
            { key: 'settings', icon: <SettingOutlined />, label: '配置' },
          ]}
          style={{ background: 'transparent' }}
        />
      </Sider>
      <Layout>
        <Header style={{ background: colorBgContainer, backgroundColor: 'var(--bg-primary)' }}>
          <div className="topbar">赛博朋克风格 Demo</div>
        </Header>
        <Content style={{ margin: 16 }}>
          <Dashboard />
        </Content>
      </Layout>
    </Layout>
  );
}

