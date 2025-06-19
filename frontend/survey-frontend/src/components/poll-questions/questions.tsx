import React, { useState } from "react";
import { SurveyQuestionWithAnswer } from "@/types/interview";
import Radio from "@/components/common/Radio";
import Input from "../common/Input";
import Checkbox from "../common/Checkbox";
import Textarea from "../common/Textarea";

interface SingleChoiceProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const SingleChoiceQuestion = ({ question, onUpdate, error }: SingleChoiceProps) => {
  const [selected, setSelected] = useState<string | null>(question.answer || null);

  const handleChange = (value: string) => {
    setSelected(value);
    onUpdate(value);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <ul className="space-y-2">
        {question.options?.map((option) => (
          <li key={option.id} className="flex items-center space-x-2">
            <Radio
              name={`q-${question.id}`}
              value={option.id.toString()}
              checked={selected === option.id.toString()}
              onChange={() => handleChange(option.id.toString())}
            />
            <span>{option.label}</span>
          </li>
        ))}
      </ul>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface MultiChoiceProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string[]) => void;
  error?: string; // Добавляем пропс для ошибки
}export const MultipleOptionQuestion = ({ question, onUpdate, error }: MultiChoiceProps) => {
  const initialAnswers: number[] = question.answer ? JSON.parse(question.answer) : [];
  const [selected, setSelected] = useState<number[]>(initialAnswers);

  const handleChange = (value: string) => {
    const numValue = parseInt(value, 10);
    const newSelected = selected.includes(numValue)
      ? selected.filter((item) => item !== numValue)
      : [...selected, numValue];
    setSelected(newSelected);
    onUpdate(newSelected);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <p className="text-sm text-gray-500">
        {question.extra_params.minAnswersCount && question.extra_params.maxAnswersCount
          ? `Выберите не менее ${question.extra_params.minAnswersCount} и не более ${question.extra_params.maxAnswersCount} ответов`
          : question.extra_params.minAnswersCount
          ? `Выберите не менее ${question.extra_params.minAnswersCount} ответов`
          : question.extra_params.maxAnswersCount
          ? `Выберите не более ${question.extra_params.maxAnswersCount} ответов`
          : ''}
      </p>
      <ul className="space-y-2">
        {question.options?.map((option) => (
          <li key={option.id} className="flex items-center space-x-2">
            <Checkbox
              name={`q-${question.id}-${option.id}`}
              checked={selected.includes(option.id)}
              onChange={() => handleChange(option.id.toString())}
              disabled={false}
            />
            <span>{option.label}</span>
          </li>
        ))}
      </ul>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface ShortTextProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const ShortTextQuestion = ({ question, onUpdate, error }: ShortTextProps) => {
  const [value, setValue] = useState<string>(question.answer || "");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    onUpdate(e.target.value);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <Input
        type="text"
        name={`q-${question.id}-answer`}
        value={value}
        onChange={handleChange}
        disabled={false}
        errors={undefined}
      />
      <p className="text-sm text-gray-500">
        {question.extra_params.maxLength
          ? `(${value.length}/${question.extra_params.maxLength} символов)`
          : ''}
      </p>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface LongTextProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const LongTextQuestion = ({ question, onUpdate, error }: LongTextProps) => {
  const [value, setValue] = useState<string>(question.answer || '');

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    setValue(newValue);
    onUpdate(newValue);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <Textarea
        name={`q-${question.id}-answer`}
        value={value}
        onChange={handleChange}
        disabled={false}
        errors={undefined}
      />
      <p className="text-sm text-gray-500">
        {question.extra_params.maxLength
          ? `(${value.length}/${question.extra_params.maxLength} символов)`
          : ''}
      </p>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface ConsentProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: boolean) => void;
}

export const ConsentQuestion = ({ question, onUpdate }: ConsentProps) => {
  const [checked, setChecked] = useState<boolean>(question.answer === "true");

  const handleChange = () => {
    setChecked(!checked);
    onUpdate(!checked);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">{question.label}</p>
      <div className="flex items-center space-x-2">
        <Checkbox
          name={`q-${question.id}-consent`}
          checked={checked}
          onChange={handleChange}
          disabled={false}
        />
        <span>Согласен</span>
      </div>
    </div>
  );
};

interface DateProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string) => void;
  error?: string; // Добавляем пропс для ошибки
}

export const DateQuestion = ({ question, onUpdate, error }: DateProps) => {
  const [value, setValue] = useState<string>(question.answer || "");


  const formatDateToYMD = (isoDate: string): string => {
    if (!isoDate) return "";
    const date = new Date(isoDate);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Месяцы начинаются с 0
    const day = String(date.getDate()).padStart(2, "0");
    return `${year}-${month}-${day}`;
  };

// Преобразуем minDate и maxDate в формат yyyy-mm-dd
  const minDateFormatted = formatDateToYMD(question.extra_params.minDate);
  const maxDateFormatted = formatDateToYMD(question.extra_params.maxDate);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    onUpdate(e.target.value);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <input
        type="date"
        name={`q-${question.id}-date`}
        value={value}
        onChange={handleChange}
min={minDateFormatted}
        max={maxDateFormatted}
        className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none"
      />
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};
interface RatingProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: number) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const RatingQuestion = ({ question, onUpdate, error }: RatingProps) => {
  const [rating, setRating] = useState<number>(
    question.answer ? parseInt(question.answer, 10) : 0
  );

  // Получаем starsCount из extra_params или устанавливаем значение по умолчанию 5
  const starsCount = question.extra_params?.starsCount
    ? parseInt(question.extra_params.starsCount, 10)
    : 5;

  const handleClick = (value: number) => {
    setRating(value);
    onUpdate(value);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params?.required && <span className="text-red-500">*</span>}
      </p>
      <div className="flex">
        {Array.from({ length: starsCount }, (_, index) => index + 1).map((star) => (
          <span
            key={star}
            className={`text-2xl cursor-pointer ${
              star <= rating ? "text-yellow-500" : "text-gray-300"
            }`}
            onClick={() => handleClick(star)}
          >
            ★
          </span>
        ))}
      </div>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface NumberProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: number) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const NumberQuestion = ({ question, onUpdate, error }: NumberProps) => {
  const [value, setValue] = useState<string>(question.answer || "");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const numValue = e.target.value;
    setValue(numValue);
    onUpdate(parseFloat(numValue));
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <Input
        type="number"
        name={`q-${question.id}-answer`}
        value={value}
        onChange={handleChange}
        disabled={false}
        errors={undefined}
      />
      <p className="text-sm text-gray-500">
        {question.extra_params.minNumber && question.extra_params.maxNumber
          ? `Введите число от ${question.extra_params.minNumber} до ${question.extra_params.maxNumber}`
          : question.extra_params.minNumber
          ? `Введите число не меньше ${question.extra_params.minNumber}`
          : question.extra_params.maxNumber
          ? `Введите число не больше ${question.extra_params.maxNumber}`
          : ''}
      </p>
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};

interface EmailProps {
  question: SurveyQuestionWithAnswer;
  onUpdate: (answer: string) => void;
  error?: string; // Добавляем пропс для ошибки
}
export const EmailQuestion = ({ question, onUpdate, error }: EmailProps) => {
  const [value, setValue] = useState<string>(question.answer || "");

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    onUpdate(e.target.value);
  };

  return (
    <div className="mb-4">
      <p className="font-medium">
        {question.label}
        {question.extra_params.required && <span className="text-red-500">*</span>}
      </p>
      <Input
        type="email"
        name={`q-${question.id}-answer`}
        value={value}
        onChange={handleChange}
        placeholder="sample@gmail.com"
        disabled={false}
        errors={undefined}
      />
      {error && <p className="text-red-500 text-sm mt-1">{error}</p>}
    </div>
  );
};
