import React from 'react';

interface TextareaProps {
  name: string;
  register: (e: React.ChangeEvent<HTMLTextAreaElement>) => void;
  errors: any;
}

export default function Textarea({name, register, errors }: TextareaProps) {
  return (
    <div className="mb-4">
      <textarea
        name={name}
        onChange={register}
        disabled
        className="w-full px-4 py-2 border border-gray-300 rounded-md bg-gray-100 focus:outline-none"
      />
      {errors[name] && <p className="text-red-500 text-sm">{errors[name]?.message}</p>}
    </div>
  );
}