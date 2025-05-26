import React from 'react';
import { SurveyQuestion } from '@/types/question';
import Checkbox from '../common/Checkbox';

interface Props {
  question: SurveyQuestion;
}

export default function Consent({ question }: Props) {
  return <Checkbox name={`q-${question.id}-consent`} />;
}