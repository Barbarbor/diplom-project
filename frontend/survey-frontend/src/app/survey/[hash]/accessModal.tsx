import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useAccessList, useAddEditAccess, useRemoveEditAccess } from '@/hooks/react-query/access';

interface AccessModalProps {
  isOpen: boolean;
  onClose: () => void;
  hash: string; // Add hash as a prop to fetch the correct access list
}

export const AccessModal = ({ isOpen, onClose, hash }: AccessModalProps) => {
  const { t } = useTranslation();
  const [newEmail, setNewEmail] = useState('');

  const { data: accessEmails = [], isLoading, error } = useAccessList(hash, {
    refetchOnWindowFocus: false,
  });

  const addMutation = useAddEditAccess();
  const removeMutation = useRemoveEditAccess();

  const handleAddEmail = () => {
    if (newEmail.trim() !== '') {
      addMutation.mutate(
        { hash, email: newEmail.trim() },
        {
          onSuccess: () => setNewEmail(''),
          onError: (err) => alert(t('survey.access.error.add', { message: err.message })),
        }
      );
    }
  };

  const handleRemoveEmail = (emailToRemove: string) => {
    removeMutation.mutate(
      { hash, email: emailToRemove },
      {
        onError: (err) => alert(t('survey.access.error.remove', { message: err.message })),
      }
    );
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
      <div className="bg-white p-4 rounded shadow-md w-1/2 h-1/2">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold">{t('survey.access.title')}</h2>
          <button
            className="px-4 py-2 bg-red-600 text-white rounded"
            onClick={onClose}
          >
            {t('survey.access.close')}
          </button>
        </div>
        {isLoading && <p>{t('survey.access.loading')}</p>}
        {error && <p className="text-red-600">{t('survey.access.error.fetch', { message: error.message })}</p>}
        <ul className="space-y-2">
          {accessEmails.map((email) => (
            <li key={email} className="flex justify-between items-center">
              <span>{email}</span>
              <button
                onClick={() => handleRemoveEmail(email)}
                className="text-red-600"
                disabled={removeMutation.isLoading}
              >
                {t('survey.access.remove')}
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
            placeholder={t('survey.access.emailPlaceholder')}
            disabled={addMutation.isLoading}
          />
          <button
            onClick={handleAddEmail}
            className="ml-2 px-4 py-2 bg-blue-600 text-white rounded"
            disabled={addMutation.isLoading}
          >
            {t('survey.access.add')}
          </button>
        </div>
      </div>
    </div>
  );
};