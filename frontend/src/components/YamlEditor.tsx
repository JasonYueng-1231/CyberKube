import { lazy, Suspense } from 'react';
import { Spin } from 'antd';

const Monaco = lazy(() => import('@monaco-editor/react'));

export default function YamlEditor({ value, onChange, height=420 }: { value: string; onChange: (v:string)=>void; height?: number }) {
  return (
    <Suspense fallback={<Spin />}>
      <Monaco height={height} defaultLanguage="yaml" value={value} onChange={(v)=> onChange(v||'')} options={{ minimap: { enabled: false }, fontSize: 13 }} />
    </Suspense>
  );
}

