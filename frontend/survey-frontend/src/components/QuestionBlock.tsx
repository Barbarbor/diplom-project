'use client';

import React, { useState } from 'react';
import { SurveyQuestion, QuestionType } from '@/types/question';
import Select from '@/components/common/Select';
import StateBadge from './StateBadge';
import Modal from '@/components/common/Modal';
import {
  useUpdateQuestionType,
  useDeleteQuestion,
  useRestoreQuestion,
} from '@/hooks/react-query/question';
import QuestionBody from './QuestionBody';
import { useSurveyHash } from '@/hooks/survey';
import RestoreIcon from './common/RestoreIcon';
import { Block } from './common/Block';

interface QuestionBlockProps {
  question: SurveyQuestion;
}
 const questionTypes = {
     [QuestionType.SingleChoice]:  'Одиночный выбор',
     [QuestionType.MultiChoice]: 'Множественный выбор' ,
    [QuestionType.Consent]: 'Согласие' ,
    [QuestionType.Rating]: 'Рейтинг' ,
    [QuestionType.ShortText]: 'Короткий текст' ,
    [QuestionType.LongText]: 'Длинный текст' ,
    [QuestionType.Date]: 'Дата' ,
    [QuestionType.Number]: 'Число' ,
    [QuestionType.Email]: 'Email' ,
 };
export default function QuestionBlock({ question }: QuestionBlockProps) {
  const hash = useSurveyHash();

  const [isTypeModalOpen, setTypeModalOpen] = useState(false);
  const [pendingType, setPendingType] = useState<QuestionType>(question.type);
  const [isDeleteModalOpen, setDeleteModalOpen] = useState(false);

  const updateType = useUpdateQuestionType();
  const deleteQ = useDeleteQuestion();
  const restoreQuestion = useRestoreQuestion();

  const handleRestoreQuestion = () => {
    restoreQuestion.mutate({ hash, questionId: question.id });
  };

  const onTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setPendingType(e.target.value as QuestionType);
    setTypeModalOpen(true);
  };

  const confirmTypeChange = () => {
    setTypeModalOpen(false);
    if (pendingType !== question.type) {
      updateType.mutate({ hash, questionId: question.id, data: { newType: pendingType } });
    }
  };

  return (
    <Block>
      <div className="flex justify-between items-start">
        {/* Left: Question Details */}
        <div className="flex-1 space-y-4">
          {/* Question Type */}
          <div>
            <h3 className="text-lg font-semibold mb-2">Тип вопроса</h3>
            <Select
              label=""
              name="questionType"
              value={question.type}
              options={Object.values(QuestionType).map((t) => ({
                value: t,
                label: questionTypes[t],
              }))}
              onChange={onTypeChange}
            />
          </div>

          {/* Question Body (Label, Body, Extra Params) */}
          <QuestionBody question={question} />
        </div>

        {/* Right: Actions */}
        <div className="flex items-center space-x-2">
          <button onClick={() => setDeleteModalOpen(true)} className="text-gray-500 hover:text-gray-700">
            🗑️
          </button>
          <RestoreIcon
            onRestore={handleRestoreQuestion}
            entityType="question"
            disabled={question.question_state === 'NEW'}
          />
          <StateBadge state={question.question_state} />
        </div>
      </div>

      {/* Confirm Modals */}
      <Modal
        isOpen={isTypeModalOpen}
        title="Подтвердите смену типа"
        onConfirm={confirmTypeChange}
        onCancel={() => setTypeModalOpen(false)}
      >
        При изменении типа вопроса прошлые данные не сохранятся. Вы уверены?
      </Modal>

      <Modal
        isOpen={isDeleteModalOpen}
        title="Удалить вопрос?"
        onConfirm={() => {
          deleteQ.mutate({ hash, questionId: question.id });
          setDeleteModalOpen(false);
        }}
        onCancel={() => setDeleteModalOpen(false)}
      >
        Вы действительно хотите удалить этот вопрос?
      </Modal>
    </Block>
  );
}