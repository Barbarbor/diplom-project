import React from 'react';
import { SurveyQuestion, QuestionType } from '@/types/question';
import { useUpdateQuestionExtraParams } from '@/hooks/react-query/question';
import DatePickerParam from './extra_params/DatePickerParam';
import NumberInputParam from './extra_params/NumberInputParam';
import SelectParam from './extra_params/SelectParam';
import ToggleParam from './extra_params/ToggleParam';

interface ExtraParamsProps {
  question: SurveyQuestion;
  hash: string;
}

// Конфигурация экстра-параметров для каждого типа вопроса
const extraParamsConfig: { [key in QuestionType]?: string[] } = {
  [QuestionType.SingleChoice]: ['required'],
  [QuestionType.MultiChoice]: ['minAnswersCount', 'maxAnswersCount', 'required'],
  [QuestionType.Consent]: [],
  [QuestionType.Email]: ['required'],
  [QuestionType.Rating]: ['required', 'starsCount'],
  [QuestionType.Date]: ['required', 'minDate', 'maxDate'],
  [QuestionType.ShortText]: ['required', 'maxLength'],
  [QuestionType.LongText]: ['required', 'maxLength'],
  [QuestionType.Number]: ['required', 'minNumber', 'maxNumber'],
};

export default function ExtraParams({ question, hash }: ExtraParamsProps) {
  const updateExtraParams = useUpdateQuestionExtraParams();

  // Функция для обновления параметров
  const handleParamChange = (param: string, value: any) => {
    // Если extra_params не определено, инициализируем его как пустой объект
    const currentParams = question.extra_params || {};
    const newExtraParams = { ...currentParams, [param]: value };
    updateExtraParams.mutate({
      hash,
      questionId: question.id,
      data: newExtraParams,
    });
  };

  // Рендеринг компонента для каждого параметра
  const renderParamComponent = (param: string) => {
    switch (param) {
      case 'required':
        return (
          <ToggleParam
            label="Обязательный"
            value={question.extra_params?.required || false}
            onChange={(value) => handleParamChange('required', value)}
          />
        );
      case 'minAnswersCount':
      case 'maxAnswersCount':
      case 'maxLength':
      case 'minNumber':
      case 'maxNumber':
        return (
          <NumberInputParam
            label={param.replace(/([A-Z])/g, ' $1').trim()}
            value={question.extra_params?.[param] || 0}
            onChange={(value) => handleParamChange(param, value)}
          />
        );
      case 'starsCount':
        return (
          <SelectParam
            label="Количество звёзд"
            value={question.extra_params?.starsCount || 5}
            options={[5, 6, 7, 8, 9, 10]}
            onChange={(value) => handleParamChange('starsCount', value)}
          />
        );
      case 'minDate':
      case 'maxDate':
        return (
          <DatePickerParam
            label={param.replace(/([A-Z])/g, ' $1').trim()}
            value={question.extra_params?.[param] || null}
            onChange={(value) => handleParamChange(param, value)}
          />
        );
      default:
        return null;
    }
  };

  const params = extraParamsConfig[question.type];
  if (!params || params.length === 0) return null;

  return (
    <div>
      <h3 className="text-lg font-semibold mb-2">Дополнительные параметры</h3>
      <div className="space-y-2">
        {params.map((param) => (
          <div key={param}>
            {renderParamComponent(param)}
          </div>
        ))}
      </div>
    </div>
  );
}