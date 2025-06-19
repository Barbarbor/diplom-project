import React from 'react';
import { SurveyQuestion, QuestionType } from '@/types/question';
import { useUpdateQuestionExtraParams } from '@/hooks/react-query/question';
import DatePickerParam from './extra_params/DatePickerParam';
import NumberInputParam from './extra_params/NumberInputParam';
import SelectParam from './extra_params/SelectParam';
import ToggleParam from './extra_params/ToggleParam';
import { useTranslation } from 'next-i18next';

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
  const { t } = useTranslation('translation', { keyPrefix: 'survey.question.extraParams' });
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
            label={t('required')}
            value={question.extra_params?.required || false}
            onChange={(value) => handleParamChange('required', value)}
          />
        );
      case 'minAnswersCount':
        return (
          <NumberInputParam
            label={t('minAnswersCount')}
            value={question.extra_params?.minAnswersCount || 0}
            onChange={(value) => handleParamChange('minAnswersCount', value)}
          />
        );
      case 'maxAnswersCount':
        return (
          <NumberInputParam
            label={t('maxAnswersCount')}
            value={question.extra_params?.maxAnswersCount || 0}
            onChange={(value) => handleParamChange('maxAnswersCount', value)}
          />
        );
      case 'maxLength':
        return (
          <NumberInputParam
            label={t('maxLength')}
            value={question.extra_params?.maxLength || 0}
            onChange={(value) => handleParamChange('maxLength', value)}
          />
        );
      case 'minNumber':
        return (
          <NumberInputParam
            label={t('minNumber')}
            value={question.extra_params?.minNumber || 0}
            onChange={(value) => handleParamChange('minNumber', value)}
          />
        );
      case 'maxNumber':
        return (
          <NumberInputParam
            label={t('maxNumber')}
            value={question.extra_params?.maxNumber || 0}
            onChange={(value) => handleParamChange('maxNumber', value)}
          />
        );
      case 'starsCount':
        return (
          <SelectParam
            label={t('starsCount')}
            value={question.extra_params?.starsCount || 5}
            options={[5, 6, 7, 8, 9, 10]}
            onChange={(value) => handleParamChange('starsCount', value)}
          />
        );
      case 'minDate':
        return (
          <DatePickerParam
            label={t('minDate')}
            value={question.extra_params?.minDate || null}
            onChange={(value) => handleParamChange('minDate', value)}
          />
        );
      case 'maxDate':
        return (
          <DatePickerParam
            label={t('maxDate')}
            value={question.extra_params?.maxDate || null}
            onChange={(value) => handleParamChange('maxDate', value)}
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
      <h3 className="text-lg font-semibold mb-2">{t('title')}</h3>
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