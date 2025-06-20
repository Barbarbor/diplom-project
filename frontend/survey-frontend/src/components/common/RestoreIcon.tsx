import React, { useState } from 'react';

interface RestoreIconProps {
  onRestore: () => void;
  entityType: 'survey' | 'question';
  disabled?: boolean; // Добавляем пропс для дизейбла
}

export default function RestoreIcon({ onRestore, entityType, disabled = false }: RestoreIconProps) {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleRestore = () => {
    setIsModalOpen(false);
    onRestore();
  };

  return (
    <>
      <button
        onClick={() => !disabled && setIsModalOpen(true)} // Открываем модал только если не disabled
        className={`text-gray-500 hover:text-gray-700 cursor-pointer z-[0] ${
          disabled ? 'opacity-50 cursor-not-allowed' : ''
        }`}
        title={`Restore ${entityType}`}
        disabled={disabled}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth="1.5"
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M9 15L3 9m0 0l6-6M3 9h12a6 6 0 010 12h-3"
          />
        </svg>
      </button>

      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white p-6 rounded shadow-lg w-96">
            <h3 className="text-lg font-semibold mb-4">
              Confirm Restore {entityType}
            </h3>
            <p className="mb-4">
              При восстановлении все предыдущие изменения не сохранятся. Продолжить?
            </p>
            <div className="flex justify-end space-x-4">
              <button
                onClick={() => setIsModalOpen(false)}
                className="px-4 py-2 bg-gray-300 text-black rounded"
              >
                Cancel
              </button>
              <button
                onClick={handleRestore}
                className="px-4 py-2 bg-blue-600 text-white rounded"
              >
                Restore
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}