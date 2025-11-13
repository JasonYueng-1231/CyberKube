import { Layout, Menu, theme, Button } from 'antd';
import {
  AppstoreOutlined,
  ClusterOutlined,
  DeploymentUnitOutlined,
  SettingOutlined,
  ApiOutlined,
  ProfileOutlined,
} from '@ant-design/icons';
import Dashboard from './pages/Dashboard';
import ClusterList from './pages/Cluster/ClusterList';
import Workloads from './pages/Workload';
import ServiceList from './pages/Network/ServiceList';
import ConfigMapList from './pages/Config/ConfigMapList';
import SecretList from './pages/Config/SecretList';
import Login from './pages/Login';
import { useEffect, useState } from 'react';
import './styles/global.css';

const { Header, Sider, Content } = Layout;

export default function App() {
  const [active, setActive] = useState('dashboard');
  const [token, setToken] = useState<string | null>(() => {
    const t = (localStorage.getItem('token') || '').trim();
    return t ? t : null;
  });
  useEffect(() => {
    const onLogout = () => setToken(null);
    window.addEventListener('auth-logout', onLogout as any);
    return () => window.removeEventListener('auth-logout', onLogout as any);
  }, []);
  // 交由各 API 的 401 统一处理，不在此处强制拉取校验，避免首次登录时闪退
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  if (!token) {
    return <Login onSuccess={(t) => { localStorage.setItem('token', t); setToken(t); }} />;
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider width={220} theme="dark" style={{ background: 'var(--bg-sidebar)' }}>
        <div className="logo">CYBERKUBE</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[active]}
          onClick={(e) => setActive(e.key)}
          items={[
            { key: 'dashboard', icon: <AppstoreOutlined />, label: '仪表盘' },
            { key: 'clusters', icon: <ClusterOutlined />, label: '集群' },
            { key: 'workloads', icon: <DeploymentUnitOutlined />, label: '工作负载' },
            { key: 'services', icon: <ApiOutlined />, label: '服务' },
            { key: 'configmaps', icon: <ProfileOutlined />, label: 'ConfigMap' },
            { key: 'secrets', icon: <ProfileOutlined />, label: 'Secret' },
            { key: 'settings', icon: <SettingOutlined />, label: '设置' },
          ]}
          style={{ background: 'transparent' }}
        />
      </Sider>
      <Layout>
        <Header style={{ background: colorBgContainer, backgroundColor: 'var(--bg-primary)', display:'flex', alignItems:'center', justifyContent:'space-between' }}>
          <div className="topbar">赛博朋克风格 Demo</div>
          <Button onClick={()=> { localStorage.removeItem('token'); setToken(null); try { window.dispatchEvent(new CustomEvent('auth-logout')); } catch {} }}>
            退出登录
          </Button>
        </Header>
        <Content style={{ margin: 16 }}>
          {active === 'dashboard' && <Dashboard />}
          {active === 'clusters' && <ClusterList />}
          {active === 'workloads' && <Workloads />}
          {active === 'services' && <ServiceList />}
          {active === 'configmaps' && <ConfigMapList />}
          {active === 'secrets' && <SecretList />}
        </Content>
      </Layout>
    </Layout>
  );
}
