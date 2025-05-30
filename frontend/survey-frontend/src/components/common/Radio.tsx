// /components/common/Radio.tsx
import React from 'react';

interface RadioProps {
  name: string;
  value: string;
  disabled?: boolean;
  checked?: boolean;
  onChange?: () => void;
}

export default function Radio({onChange = ()=> {}, name, value, disabled = false, checked=false }: RadioProps) {
  return (
    <div className="mb-4">
      <input onChange={onChange} type="radio" name={name} value={value} disabled={disabled} className="mr-2" checked={checked} />
    </div>
  );
}