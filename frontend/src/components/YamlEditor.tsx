import { lazy, Suspense, useMemo } from 'react';
import { Spin } from 'antd';

const Monaco = lazy(() => import('@monaco-editor/react'));

type Props = {
  value: string;
  onChange: (v: string) => void;
  height?: number;
  readOnly?: boolean;
};

export default function YamlEditor({ value, onChange, height = 420, readOnly = false }: Props) {
  const options = useMemo(
    () => ({
      minimap: { enabled: false },
      fontSize: 13,
      readOnly,
      scrollBeyondLastLine: false,
      wordWrap: 'on',
    }),
    [readOnly],
  );

  return (
    <Suspense fallback={<Spin />}>
      <Monaco
        height={height}
        defaultLanguage="yaml"
        value={value}
        onChange={(v) => onChange(v || '')}
        options={options}
      />
    </Suspense>
  );
}
