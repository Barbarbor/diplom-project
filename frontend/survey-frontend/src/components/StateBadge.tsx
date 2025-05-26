import React from 'react';
import { QuestionState } from '@/types/question';

const colorMap: Record<QuestionState, string> = {
  ACTUAL: 'bg-green-200 text-green-800',
  NEW:    'bg-blue-200 text-blue-800',
  CHANGED:'bg-yellow-200 text-yellow-800',
  DELETED:'bg-red-200 text-red-800',
};

export default function StateBadge({ state }: { state: QuestionState }) {
  return (
    <span className={`px-2 py-1 rounded text-sm ${colorMap[state]}`}>
      {state.toLowerCase()}
    </span>
  );
}
