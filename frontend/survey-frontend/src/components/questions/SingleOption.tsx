import React, { useState } from 'react';
import { Option } from '@/types/option';
import { useUpdateOptionLabel, useDeleteOption } from '@/hooks/react-query/option';

import Radio from '../common/Radio';
import EditableLabel from './EditableLabel';

interface Props {
  hash: string;
  questionId: number;
  option: Option;
}

export default function SingleOption({ hash, questionId, option }: Props) {
  const [label, setLabel] = useState(option.label || '');
  const updateLabel = useUpdateOptionLabel();
  const deleteOpt = useDeleteOption();

  const handleLabelChange = (newLabel: string) => {
    setLabel(newLabel);
    updateLabel.mutate({
      hash,
      questionId,
      optionId: option.id,
      data: { label: newLabel },
    });
  };

  return (
    <li className="flex items-center space-x-2">
      <Radio name={`option-${questionId}`} value={option.id.toString()} disabled />
      <div className="flex-1">
        <EditableLabel initialLabel={label} onLabelChange={handleLabelChange} />
      </div>
      <button onClick={() => deleteOpt.mutate({ hash, questionId, optionId: option.id })}>
        🗑️
      </button>
    </li>
  );
}