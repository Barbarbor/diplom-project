// src/components/ToggleParam.tsx
import Toggle from '@/components/common/Toggle';
import React from 'react';


interface ToggleParamProps {
  label: string;
  value: boolean;
  onChange: (value: boolean) => void;
}

export default function ToggleParam({ label, value, onChange }: ToggleParamProps) {
  return (
    <div className="flex items-center space-x-2">
      <label>{label}</label>
      <Toggle checked={value} onChange={onChange} />
    </div>
  );
}