import React from 'react';
import { QuestionState } from '@/types/question';
import { OptionState } from '@/types/option';

const colorMap: Record<QuestionState, string> = {
  ACTUAL: 'bg-green-200 text-green-800',
  NEW:    'bg-blue-200 text-blue-800',
  CHANGED:'bg-yellow-200 text-yellow-800',
  DELETED:'bg-red-200 text-red-800',
};
const colorLabel = { ACTUAL: 'Актуальный',
  NEW:    'Новый',
  CHANGED:'Изменён',
  DELETED:'Удалён',}
export default function StateBadge({ state }: { state: QuestionState | OptionState }) {
  return (
    <span className={`px-2 py-1 rounded text-sm ${colorMap[state]}`}>
      {colorLabel[state]}
    </span>
  );
}
