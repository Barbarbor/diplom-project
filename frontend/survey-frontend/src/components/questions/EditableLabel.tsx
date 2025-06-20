import React, { useState, useEffect } from 'react';

interface EditableLabelProps {
  initialLabel: string;
  onLabelChange: (newLabel: string) => void;
}

export default function EditableLabel({ initialLabel, onLabelChange }: EditableLabelProps) {

  

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onLabelChange(e.target.value)
  };

  return (
    <input
      type="text"
      value={initialLabel}
      onChange={handleChange}
      className="block text-gray-700 font-bold p-1 w-full border-none focus:outline-none"
    />
  );
}