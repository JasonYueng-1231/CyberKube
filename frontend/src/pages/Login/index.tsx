import { Button, Card, Form, Input, message } from 'antd';
import { api } from '@/utils/request';

export default function Login({ onSuccess }: { onSuccess: (token: string) => void }) {
  const [form] = Form.useForm();

  const onFinish = async (v: any) => {
    try {
      const data = await api('/auth/login', { method: 'POST', body: JSON.stringify(v) });
      onSuccess(data.token);
    } catch (e: any) {
      message.error(e.message);
    }
  };

  return (
    <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100vh' }}>
      <Card title="登录 CyberKube" style={{ width: 360 }} className="cyber-card">
        <Form form={form} onFinish={onFinish} layout="vertical" initialValues={{ username: 'admin', password: 'admin123' }}>
          <Form.Item name="username" label="用户名" rules={[{ required: true }]}>
            <Input placeholder="请输入用户名" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true }]}>
            <Input.Password placeholder="请输入密码" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block>登录</Button>
        </Form>
      </Card>
    </div>
  );
}

