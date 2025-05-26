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


interface QuestionBlockProps {
  question: SurveyQuestion;
}

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
    <div className="border p-4 rounded mb-4 bg-white shadow-sm">
      <div className="flex justify-between items-start">
        {/* Middle: Form fields */}
        <div className="flex-1">
          {/* Type selector */}
          <Select
            label="Тип вопроса"
            name="questionType"
            value={question.type}
            options={Object.values(QuestionType).map((t) => ({
              value: t,
              label: t.replace('_', ' '),
            }))}
            onChange={onTypeChange}
          />
          <QuestionBody question={question} />
        </div>

        {/* Right: Actions */}
        <div className="space-y-2 text-right">
          <StateBadge state={question.question_state} />

          {question.question_state !== 'NEW' && (
            <RestoreIcon onRestore={handleRestoreQuestion} entityType="question" />
          )}

          <button onClick={() => setDeleteModalOpen(true)}>🗑️</button>
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
    </div>
  );
}