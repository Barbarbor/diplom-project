import React, { useState } from 'react';
import { SurveyQuestion, QuestionType } from '@/types/question';
import { useUpdateQuestionLabel } from '@/hooks/react-query/question';
import { UpdateQuestionLabelRequest } from '@/types/question';
import EditableLabel from './questions/EditableLabel';
import SingleChoice from './questions/SingleChoice';
import MultipleOption from './questions/MultipleOption';
import ShortText from './questions/ShortText';
import LongText from './questions/LongText';
import Email from './questions/Email';
import NumberInput from './questions/NumberInput';
import Consent from './questions/Consent';
import DatePicker from './questions/DatePicker';
import Rating from './questions/Rating';
import { useSurveyHash } from '@/hooks/survey';
import ExtraParams from './questions/ExtraParams';

interface Props {
  question: SurveyQuestion;
}

export default function QuestionBody({ question }: Props) {
  const hash = useSurveyHash();
  const updateLabel = useUpdateQuestionLabel();
  const [label, setLabel] = useState(question.label || 'Default Label');

  const handleLabelChange = (newLabel: string) => {
    setLabel(newLabel);
    const payload: UpdateQuestionLabelRequest = { label: newLabel };
    updateLabel.mutate({
      hash,
      questionId: question.id,
      data: payload,
    });
  };

  const renderQuestionComponent = () => {
    switch (question.type) {
      case QuestionType.SingleChoice:
        return <SingleChoice question={question} hash={hash} />;
      case QuestionType.MultiChoice:
        return <MultipleOption question={question} hash={hash} />;
      case QuestionType.ShortText:
        return <ShortText question={question} />;
      case QuestionType.LongText:
        return <LongText question={question} />;
      case QuestionType.Email:
        return <Email question={question} />;
      case QuestionType.Number:
        return <NumberInput question={question} />;
      case QuestionType.Consent:
        return <Consent question={question} />;
      case QuestionType.Date:
        return <DatePicker question={question} />;
      case QuestionType.Rating:
        return <Rating question={question} />;
      default:
        return null;
    }
  };

  return (
    <div className="space-y-4">
      {/* Label Section */}
      <div>
        <EditableLabel initialLabel={label} onLabelChange={handleLabelChange} />
      </div>

      {/* Question Body Section */}
      <div>
        {renderQuestionComponent()}
      </div>


      {/* Extra Params Section */}
      <ExtraParams question={question} hash={hash} />
    </div>
  );
}