import React from 'react';
import { SurveyQuestion } from '@/types/question';

interface Props {
  question: SurveyQuestion;
}

export default function Rating({ question }: Props) {
  return (
    <div className="flex">
      {[...Array(5)].map((_, i) => (
        <span key={i} className="text-yellow-500 text-xl">★</span>
      ))}
    </div>
  );
}