import React from 'react';
import { SurveyQuestion } from '@/types/question';
import Textarea from '../common/Textarea';

interface Props {
  question: SurveyQuestion;
}

export default function LongText({ question }: Props) {
  return (
    <Textarea
      name={`q-${question.id}-answer`}
      disabled={true}
      register={() => {}}
      errors={{}}
    />
  );
}