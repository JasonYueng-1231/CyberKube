import { lazy, Suspense } from 'react';
import { Spin } from 'antd';

const Monaco = lazy(() => import('@monaco-editor/react'));

interface Props {
  value: string;
  onChange: (v: string) => void;
  height?: number;
  readOnly?: boolean;
}

export default function YamlEditor({ value, onChange, height = 420, readOnly = false }: Props) {
  return (
    <Suspense fallback={<Spin />}>
      <Monaco
        height={height}
        defaultLanguage="yaml"
        value={value}
        onChange={(v) => onChange(v || '')}
        options={{ minimap: { enabled: false }, fontSize: 13, readOnly }}
      />
    </Suspense>
  );
}
