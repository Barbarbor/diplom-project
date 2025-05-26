import React, { useState, useEffect } from 'react';
import { useDebounce } from '@uidotdev/usehooks';

interface EditableLabelProps {
  initialLabel: string;
  onLabelChange: (newLabel: string) => void;
}

export default function EditableLabel({ initialLabel, onLabelChange }: EditableLabelProps) {
  const [label, setLabel] = useState(initialLabel);
  const debouncedLabel = useDebounce(label, 300); // Задержка 300 мс

  // Обновляем onLabelChange только когда debouncedLabel меняется
  useEffect(() => {
    if (debouncedLabel !== initialLabel) {
      onLabelChange(debouncedLabel);
    }
  }, [debouncedLabel, onLabelChange, initialLabel]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newLabel = e.target.value;
    setLabel(newLabel);
  };

  return (
    <input
      type="text"
      value={label}
      onChange={handleChange}
      className="block text-gray-700 font-bold mb-2 w-full border-none focus:outline-none"
    />
  );
}