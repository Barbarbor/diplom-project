import React from 'react';
import { SurveyQuestion } from '@/types/question';

interface Props {
  question: SurveyQuestion;
}

export default function Rating({ question }: Props) {
  // Получаем starsCount из extra_params или устанавливаем значение по умолчанию 5
  const starsCount = question.extra_params?.starsCount
    ? parseInt(question.extra_params.starsCount, 10)
    : 5;

  return (
    <div className="flex">
      {[...Array(starsCount)].map((_, i) => (
        <span key={i} className="text-yellow-500 text-xl">★</span>
      ))}
    </div>
  );
}