// src/components/SelectParam.tsx
import React from 'react';

interface SelectParamProps {
  label: string;
  value: number;
  options: number[];
  onChange: (value: number) => void;
}

export default function SelectParam({ label, value, options, onChange }: SelectParamProps) {
  return (
    <div className="flex items-center space-x-2">
      <label>{label}</label>
      <select
        value={value}
        onChange={(e) => onChange(parseInt(e.target.value, 10))}
        className="border border-gray-300 rounded px-2 py-1"
      >
        {options.map((opt) => (
          <option key={opt} value={opt}>
            {opt}
          </option>
        ))}
      </select>
    </div>
  );
}