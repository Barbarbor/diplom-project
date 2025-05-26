import React from 'react';
import { SurveyQuestion } from '@/types/question';
import Input from '@/components/common/Input';

interface Props {
  question: SurveyQuestion;
}

export default function Email({ question }: Props) {
  return (
    <Input
      type="email"
      name={`q-${question.id}-answer`}
      register={() => {}}
      errors={{}}
      disabled
    />
  );
}