import React from 'react';

interface InputProps {
  type: string;
  name: string;
  register?: any;
  errors: any;
  label?: string;
  value?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (e: React.FocusEvent<HTMLInputElement>) => void;
  defaultValue?: string;
  disabled?: boolean;
}

export default function Input({
  type,
  name,
  register,
  errors,
  label,
  value,
  onChange,
  onBlur,
  defaultValue,
  disabled,
}: InputProps) {
  return (
    <div className="mb-4">
      {label && (
        <label htmlFor={name} className="block text-sm font-medium text-gray-700 mb-1">
          {label}
        </label>
      )}
      <input
        type={type}
        id={name}
        {...(register ? register(name) : { value, onChange })}
        onBlur={onBlur}
        defaultValue={defaultValue}
        disabled={disabled}
        className={`w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 ${errors[name] ? 'border-red-500' : ''} ${disabled ? 'bg-gray-100' : ''}`}
      />
      {errors[name] && <p className="text-red-500 text-sm">{errors[name]?.message}</p>}
    </div>
  );
}