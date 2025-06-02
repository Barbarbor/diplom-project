'use client';

import React from 'react';
import { Block } from '@/components/common/Block';

interface SurveyDistributionProps {
  hash: string;

}

export const SurveyDistribution = ({ hash }: SurveyDistributionProps) => {
  const surveyLink = `http://localhost:3000/poll/${hash}`;

  const handleCopyLink = () => {
    navigator.clipboard.writeText(surveyLink).then(() => {
      alert('Ссылка скопирована в буфер обмена!');
    }).catch(err => {
      console.error('Ошибка копирования:', err);
    });
  };

  return (
    <Block>
      <h3 className="text-lg font-semibold mb-2">Распространение опроса</h3>
      <p className="mb-2">Опрос будет доступен по ссылке:</p>
      <div className="flex items-center space-x-2">
        <input
          type="text"
          value={surveyLink}
          readOnly
          disabled
          className="flex-1 p-2 border border-gray-300 rounded bg-gray-100 text-gray-700"
        />
        <button
          onClick={handleCopyLink}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"

        >
          Копировать
        </button>
      </div>
    </Block>
  );
};