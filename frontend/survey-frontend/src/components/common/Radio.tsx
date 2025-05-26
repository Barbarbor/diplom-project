// /components/common/Radio.tsx
import React from 'react';

interface RadioProps {
  name: string;
  value: string;
  disabled?: boolean;
}

export default function Radio({ name, value, disabled = false }: RadioProps) {
  return (
    <div className="mb-4">
      <input type="radio" name={name} value={value} disabled={disabled} className="mr-2" />
    </div>
  );
}