import { Modal, Select, Button } from 'antd';
import { useEffect, useRef, useState } from 'react';

interface Props {
  open: boolean;
  onClose: () => void;
  cluster: string;
  namespace: string;
  pod: string;
  containers?: string[];
}

export default function LogStreamModal({ open, onClose, cluster, namespace, pod, containers = [] }: Props) {
  const [container, setContainer] = useState<string>(containers[0] || '');
  const [paused, setPaused] = useState(false);
  const [lines, setLines] = useState<string[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const boxRef = useRef<HTMLPreElement | null>(null);

  useEffect(() => { if (containers?.length) setContainer(containers[0]); }, [containers.join(',')]);

  useEffect(() => {
    if (!open) return;
    const token = localStorage.getItem('token') || '';
    const proto = location.protocol === 'https:' ? 'wss' : 'ws';
    const url = `${proto}://${location.host}/api/v1/pods/logs/stream?cluster=${encodeURIComponent(cluster)}&namespace=${encodeURIComponent(namespace)}&name=${encodeURIComponent(pod)}&container=${encodeURIComponent(container||'')}&tail=200&token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(url);
    wsRef.current = ws;
    ws.onmessage = (ev) => {
      if (paused) return;
      setLines((prev) => {
        const next = [...prev, String(ev.data)];
        if (next.length > 2000) next.shift();
        return next;
      });
      if (boxRef.current) boxRef.current.scrollTop = boxRef.current.scrollHeight;
    };
    ws.onerror = () => setLines((p)=>[...p, '[error] 日志流连接异常']);
    ws.onclose = () => setLines((p)=>[...p, '[info] 连接已关闭']);
    return () => { ws.close(); };
  }, [open, cluster, namespace, pod, container]);

  return (
    <Modal width={900} title={`实时日志: ${pod}`} open={open} onCancel={onClose} footer={null}>
      <div style={{ display:'flex', gap:8, marginBottom:8 }}>
        <Select value={container} onChange={setContainer} options={(containers||[]).map(c=>({label:c, value:c}))} style={{ minWidth:200 }} />
        <Button onClick={()=> setPaused(p=>!p)}>{paused ? '继续' : '暂停'}</Button>
        <Button onClick={()=> { navigator.clipboard.writeText(lines.join('\n')); }}>复制</Button>
      </div>
      <pre ref={boxRef} style={{ maxHeight: 500, overflow:'auto', background:'#000', color:'#0f0', padding:12 }}>{lines.join('\n')}</pre>
    </Modal>
  );
}

