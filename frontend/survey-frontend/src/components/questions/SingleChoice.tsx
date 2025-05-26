import React from 'react';
import { SurveyQuestion } from '@/types/question';
import { useCreateOption } from '@/hooks/react-query/option';
import SingleOption from './SingleOption';

interface Props {
  question: SurveyQuestion;
  hash: string;
}

export default function SingleChoice({ question, hash }: Props) {
  const createOpt = useCreateOption();

  return (
    <div>
      <ul className="space-y-2">
        {question.options?.map(opt => (
          <SingleOption
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
        + Add option
      </button>
    </div>
  );
}