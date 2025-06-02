import React from 'react';
import { SurveyQuestion } from '@/types/question';
import { useCreateOption } from '@/hooks/react-query/option';
import MultipleOptionItem from './MultipleOptionItem'; // Новый компонент для одной опции

interface Props {
  question: SurveyQuestion;
  hash: string;
}

export default function MultipleOption({ question, hash }: Props) {
  const createOpt = useCreateOption();

  return (
    <div>
      <ul className="space-y-2">
        {question.options?.map((opt) => (
          <MultipleOptionItem
            key={opt.id}
            hash={hash}
            questionId={question.id}
            option={opt}
          />
        ))}
      </ul>
      <button
        className="mt-2 text-sm text-blue-600"
        onClick={() => createOpt.mutate({ hash, questionId: question.id })}
      >
        Добавить опцию
      </button>
    </div>
  );
}