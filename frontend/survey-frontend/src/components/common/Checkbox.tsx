import React from 'react';

interface CheckboxProps {
  name: string;
}

export default function Checkbox({ name }: CheckboxProps) {
  return (
    <div className="mb-4">
      <input type="checkbox" name={name} disabled className="mr-2" />
    </div>
  );
}