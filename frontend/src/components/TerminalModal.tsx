import { Modal, Select } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { Terminal } from 'xterm';
import 'xterm/css/xterm.css';

interface Props {
  open: boolean;
  onClose: () => void;
  cluster: string;
  namespace: string;
  pod: string;
  containers?: string[];
}

export default function TerminalModal({ open, onClose, cluster, namespace, pod, containers = [] }: Props) {
  const [container, setContainer] = useState<string>(containers[0] || '');
  const termRef = useRef<HTMLDivElement | null>(null);
  const socketRef = useRef<WebSocket | null>(null);
  const termObj = useRef<Terminal | null>(null);

  useEffect(() => { if (containers?.length) setContainer(containers[0]); }, [containers.join(',')]);

  useEffect(() => {
    if (!open) return;
    const term = new Terminal({cols: 100, rows: 28, convertEol: true, cursorBlink: true});
    termObj.current = term;
    term.open(termRef.current!);

    const token = localStorage.getItem('token') || '';
    const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    const url = `${proto}://${location.host}/api/v1/pods/shell?cluster=${encodeURIComponent(cluster)}&namespace=${encodeURIComponent(namespace)}&pod=${encodeURIComponent(pod)}&container=${encodeURIComponent(container||'')}&token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(url);
    socketRef.current = ws;
    ws.onmessage = (ev) => term.write(String(ev.data));
    ws.onclose = () => term.write('\r\n[连接关闭]\r\n');
    ws.onerror = () => term.write('\r\n[错误] 连接失败\r\n');
    term.onData((d) => ws.readyState === WebSocket.OPEN && ws.send(d));
    return () => { ws.close(); term.dispose(); };
  }, [open, cluster, namespace, pod, container]);

  return (
    <Modal width={900} title={`WebShell: ${pod}`} open={open} onCancel={onClose} footer={null}>
      <div style={{ marginBottom:8 }}>
        <Select value={container} onChange={setContainer} options={(containers||[]).map(c=>({label:c, value:c}))} style={{ minWidth:200 }} />
      </div>
      <div ref={termRef} style={{ height: 520, background: '#000' }} />
    </Modal>
  );
}

