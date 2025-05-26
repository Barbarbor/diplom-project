import React from 'react';
import { SurveyQuestion } from '@/types/question';
import Input from '@/components/common/Input';

interface Props {
  question: SurveyQuestion;
}

export default function ShortText({ question }: Props) {
  return (
    <Input
      type="text"
      name={`q-${question.id}-answer`}
      register={() => {}}
      errors={{}}
      disabled
    />
  );
}