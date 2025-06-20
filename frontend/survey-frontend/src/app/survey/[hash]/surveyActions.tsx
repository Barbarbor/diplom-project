'use client';

import React, { useState } from 'react';
import { usePublishSurvey } from '@/hooks/react-query/survey';
import RestoreIcon from '@/components/common/RestoreIcon';
import { QuestionType } from '@/api-client/question';
import { SurveyState } from '@/types/survey';

interface SurveyActionsProps {
  state: SurveyState;
  hash: string;
  questionCount: number; // Новый проп для количества вопросов
  onAddQuestion: (type: QuestionType) => void;
  onOpenAccessModal: () => void;
  onOpenPreviewModal: () => void;
  onRestoreSurvey: () => void;
}

export const SurveyActions = ({
  state,
  hash,
  questionCount,
  onAddQuestion,
  onOpenAccessModal,
  onOpenPreviewModal,
  onRestoreSurvey,
}: SurveyActionsProps) => {
  const publishSurvey = usePublishSurvey();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);

  const handlePublish = () => {
    publishSurvey.mutate(hash);
  };

  const questionTypes = [
    { type: QuestionType.SingleChoice, label: 'Одиночный выбор' },
    { type: QuestionType.MultiChoice, label: 'Множественный выбор' },
    { type: QuestionType.Consent, label: 'Согласие' },
    { type: QuestionType.Rating, label: 'Рейтинг' },
    { type: QuestionType.ShortText, label: 'Короткий текст' },
    { type: QuestionType.LongText, label: 'Длинный текст' },
    { type: QuestionType.Date, label: 'Дата' },
    { type: QuestionType.Number, label: 'Число' },
    { type: QuestionType.Email, label: 'Email' },
  ];

  const handleSelectQuestionType = (type: QuestionType) => {
    onAddQuestion(type);
    setIsDropdownOpen(false);
  };

  return (
    <div className="flex space-x-4">
      <div className="fixed top-60 right-6 w-48 p-4 bg-white border border-gray-300 rounded shadow-md">
        <div className="relative">
          <button
            className="w-full mb-2 px-4 py-2 bg-blue-600 text-white rounded"
            onClick={() => setIsDropdownOpen(!isDropdownOpen)}
          >
            Добавить вопрос
          </button>
          {isDropdownOpen && (
            <ul className="absolute top-full left-0 w-full bg-white border border-gray-300 rounded shadow-md z-10">
              {questionTypes.map(({ type, label }) => (
                <li
                  key={type}
                  className="px-4 py-2 hover:bg-gray-100 cursor-pointer z-10"
                  onClick={() => handleSelectQuestionType(type)}
                >
                  {label}
                </li>
              ))}
            </ul>
          )}
        </div>
        <button
          className={`w-full px-4 py-2 bg-gray-500 text-white rounded ${questionCount === 0 ? 'opacity-50 cursor-not-allowed' : ''}`}
          onClick={questionCount === 0 ? undefined : onOpenPreviewModal}
          disabled={questionCount === 0}
        >
          Предпросмотр
        </button>
      </div>
      <div className="mt-4 flex space-x-4">
        <button
          className="px-4 py-2 bg-green-600 text-white rounded"
          onClick={handlePublish}
        >
          Опубликовать опрос
        </button>
        <button
          className="px-4 py-2 bg-blue-600 text-white rounded"
          onClick={onOpenAccessModal}
        >
          Доступ к опросу
        </button>
        <RestoreIcon disabled={state === SurveyState.Draft} onRestore={onRestoreSurvey} entityType="survey" />
      </div>
    </div>
  );
};