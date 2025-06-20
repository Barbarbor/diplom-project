import React from 'react';

interface CheckboxProps {
  name: string;
  checked?: boolean
  onChange?: ()=> void;
  disabled?: boolean;
}

export default function Checkbox({ onChange = ()=> {}, name, checked = false, disabled = true }: CheckboxProps) {
  return (
    <div >
      <input onChange={onChange} type="checkbox" name={name} disabled={disabled} className="mr-2" checked={checked} />
    </div>
  );
}