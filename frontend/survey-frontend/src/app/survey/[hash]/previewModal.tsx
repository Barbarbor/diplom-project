import React from 'react';

interface PreviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  hash: string;
}

export const PreviewModal = ({ isOpen, onClose, hash }: PreviewModalProps) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-4 rounded shadow-md w-11/12 h-5/6 flex flex-col">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Предпросмотр опроса</h2>
          <button
            className="px-4 py-2 bg-red-600 text-white rounded"
            onClick={onClose}
          >
            Закрыть
          </button>
        </div>
        <div className="flex-grow overflow-hidden">
          <iframe
            src={`/poll/${hash}?isDemo=true`}
            className="w-full h-full border-0"
            title="Survey Preview"
            style={{ display: 'block' }} // Устраняем возможные отступы
          />
        </div>
      </div>
    </div>
  );
};