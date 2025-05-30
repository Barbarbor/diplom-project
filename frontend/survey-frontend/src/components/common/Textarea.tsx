import React from 'react';

interface TextareaProps {
  name: string;
  value?: string; // Добавляем value для управляемого ввода
  onChange?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void; // Добавляем onChange
  register?: (e: React.ChangeEvent<HTMLTextAreaElement>) => void; // Делаем register необязательным
  disabled?: boolean; // Делаем disabled необязательным
  errors?: any; // Делаем errors необязательным
}

export default function Textarea({
  name,
  value,
  onChange,
  register,
  disabled = false, // По умолчанию не отключаем
  errors,
}: TextareaProps) {
  return (
    <div className="mb-4">
      <textarea
        name={name}
        value={value}
        onChange={(e) => {
          // Если передан onChange (управляемый режим), вызываем его
          if (onChange) onChange(e);
          // Если передан register (неуправляемый режим), вызываем его
          if (register) register(e);
        }}
        disabled={disabled}
        className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none disabled:bg-gray-100"
      />
      {errors?.[name] && <p className="text-red-500 text-sm">{errors[name]?.message}</p>}
    </div>
  );
}