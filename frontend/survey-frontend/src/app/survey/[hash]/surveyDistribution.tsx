'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Block } from '@/components/common/Block';

interface SurveyDistributionProps {
  hash: string;
}

export const SurveyDistribution = ({ hash }: SurveyDistributionProps) => {
  const surveyLink = `http://localhost:3000/poll/${hash}`;
  const [showTooltip, setShowTooltip] = useState(false);
  const tooltipTimeout = useRef<number>(null);

  const handleCopyLink = () => {
    navigator.clipboard.writeText(surveyLink)
      .then(() => {
        setShowTooltip(true);
        // hide after 1.5 seconds
        tooltipTimeout.current = window.setTimeout(() => {
          setShowTooltip(false);
        }, 1500);
      })
      .catch(err => {
        console.error('Ошибка копирования:', err);
      });
  };

  // clear timeout on unmount
  useEffect(() => {
    return () => {
      if (tooltipTimeout.current) {
        clearTimeout(tooltipTimeout.current);
      }
    };
  }, []);

  return (
    <Block>
      <h3 className="text-lg font-semibold mb-2">Распространение опроса</h3>
      <p className="mb-2">Опрос будет доступен по ссылке:</p>
      <div className="flex items-center space-x-2 relative">
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
        {showTooltip && (
          <div className="absolute -top-8 left-1/2 transform -translate-x-1/2 px-3 py-1 bg-black text-white text-sm rounded">
            Ссылка скопирована!
          </div>
        )}
      </div>
    </Block>
  );
};
