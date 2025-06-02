import React, { useState } from 'react';

interface AccessModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const AccessModal = ({ isOpen, onClose }: AccessModalProps) => {
  const [accessEmails, setAccessEmails] = useState(["user1@example.com", "user2@example.com"]);
  const [newEmail, setNewEmail] = useState("");

  const handleRemoveEmail = (emailToRemove: string) => {
    setAccessEmails(accessEmails.filter((email) => email !== emailToRemove));
  };

  const handleAddEmail = () => {
    if (newEmail.trim() !== "") {
      setAccessEmails([...accessEmails, newEmail.trim()]);
      setNewEmail("");
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-4 rounded shadow-md w-1/2 h-1/2">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">Управление доступом к опросу</h2>
          <button
            className="px-4 py-2 bg-red-600 text-white rounded"
            onClick={onClose}
          >
            Закрыть
          </button>
        </div>
        <ul className="space-y-2">
          {accessEmails.map((email) => (
            <li key={email} className="flex justify-between items-center">
              <span>{email}</span>
              <button onClick={() => handleRemoveEmail(email)} className="text-red-600">
                ✖
              </button>
            </li>
          ))}
        </ul>
        <div className="flex mt-4">
          <input
            type="text"
            value={newEmail}
            onChange={(e) => setNewEmail(e.target.value)}
            className="flex-1 p-2 border border-gray-300 rounded"
            placeholder="Введите email"
          />
          <button
            onClick={handleAddEmail}
            className="ml-2 px-4 py-2 bg-blue-600 text-white rounded"
          >
            Добавить
          </button>
        </div>
      </div>
    </div>
  );
};