// src/components/NumberInputParam.tsx
import React from 'react';

interface NumberInputParamProps {
  label: string;
  value: number;
  onChange: (value: number) => void;
}

export default function NumberInputParam({ label, value, onChange }: NumberInputParamProps) {
  return (
    <div className="flex items-center space-x-2">
      <label>{label}</label>
      <input
        type="number"
        value={value}
        onChange={(e) => onChange(parseInt(e.target.value, 10) || 0)}
        className="border border-gray-300 rounded px-2 py-1 w-24"
      />
    </div>
  );
}