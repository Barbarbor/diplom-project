import React from 'react';
import { SurveyQuestion } from '@/types/question';

interface Props {
  question: SurveyQuestion;
}

export default function DatePicker({ question }: Props) {
  return (
    <input
      type="date"
      name={`q-${question.id}-date`}
      disabled
      className="w-full px-4 py-2 border border-gray-300 rounded-md bg-gray-100 focus:outline-none"
    />
  );
}