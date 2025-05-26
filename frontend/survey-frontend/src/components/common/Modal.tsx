import React from 'react';

interface ModalProps {
  isOpen: boolean;
  title?: string;
  onConfirm: () => void;
  onCancel: () => void;
  children: React.ReactNode;
}

export default function Modal({ isOpen, title, children, onConfirm, onCancel }: ModalProps) {
  if (!isOpen) return null;
  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded shadow-lg w-full max-w-md">
        {title && <h2 className="text-xl font-bold mb-4">{title}</h2>}
        <div className="mb-6">{children}</div>
        <div className="flex justify-end space-x-3">
          <button onClick={onCancel} className="px-4 py-2 border rounded">Нет</button>
          <button onClick={onConfirm} className="px-4 py-2 bg-red-500 text-white rounded">Да</button>
        </div>
      </div>
    </div>
  );
}
