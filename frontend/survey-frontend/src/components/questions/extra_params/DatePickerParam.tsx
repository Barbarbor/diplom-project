// src/components/DatePickerParam.tsx
import React from 'react';
import DatePicker from 'react-datepicker';
import 'react-datepicker/dist/react-datepicker.css';

interface DatePickerParamProps {
  label: string;
  value: Date | null;
  onChange: (date: Date | null) => void;
}

export default function DatePickerParam({ label, value, onChange }: DatePickerParamProps) {
  return (
    <div className="flex items-center space-x-2">
      <label>{label}</label>
      <DatePicker
        selected={value}
        onChange={onChange}
        className="border border-gray-300 rounded px-2 py-1"
      />
    </div>
  );
}