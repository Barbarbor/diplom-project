"use client";

import { useState } from "react";
import {
  useGetSurveyWithAnswers,
  useUpdateQuestionAnswer,
  useFinishInterview,
} from "@/hooks/react-query/interview";
import {
  SingleChoiceQuestion,
  MultipleOptionQuestion,
  ShortTextQuestion,
  LongTextQuestion,
  ConsentQuestion,
  DateQuestion,
  RatingQuestion,
  NumberQuestion,
  EmailQuestion,
} from "@/components/poll-questions/questions";
import { QuestionType } from "@/types/question";
import { useGetInterviewId } from "@/hooks/interview";
import { useParams, useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import { Block } from "@/components/common/Block";

export default function PollPage() {
  const router = useRouter();
  const params = useParams();
  const searchParams = useSearchParams();
  const isDemo = searchParams.get('isDemo');

  const hash = params.hash as string;
  const {
    interviewId,
    loading: idLoading,
    error: idError,
  } = useGetInterviewId(hash, isDemo);
  const {
    data,
    isLoading: surveyLoading,
    error: surveyError,
  } = useGetSurveyWithAnswers(hash, interviewId || "");
  const { mutate: updateAnswer } = useUpdateQuestionAnswer();
  const { mutate: finishInterview } = useFinishInterview();

  // Состояние для хранения ошибок валидации
  const [validationErrors, setValidationErrors] = useState<Record<number, string>>({});

  const handleUpdateAnswer = (questionId: number, answer: any) => {
    const formattedAnswer =
      typeof answer === "object" ? JSON.stringify(answer) : String(answer);
    updateAnswer({
      interviewId,
      hash,
      questionId,
      data: { answer: formattedAnswer },
    });
  };

  // Функция валидации вопросов
  const validateQuestions = () => {
    const errors: Record<number, string> = {};

    data?.questions.forEach((question) => {
      const answer = question.answer;

      // Проверка обязательности
      if (question.extra_params.required && (!answer || answer === "")) {
        errors[question.id] = "Необходимо ответить на этот вопрос";
        return;
      }

      // Пропускаем дальнейшую валидацию, если ответа нет и вопрос необязательный
      if (!answer) return;

      switch (question.type) {
        case QuestionType.SingleChoice:
          // Для single_choice дополнительная валидация не требуется
          break;

        case QuestionType.MultiChoice:
          const selectedAnswers: number[] = answer ? JSON.parse(answer) : [];
          if (
            question.extra_params.minAnswersCount &&
            selectedAnswers.length < question.extra_params.minAnswersCount
          ) {
            errors[question.id] = `Необходимо выбрать не менее ${question.extra_params.minAnswersCount} ответов`;
          } else if (
            question.extra_params.maxAnswersCount &&
            selectedAnswers.length > question.extra_params.maxAnswersCount
          ) {
            errors[question.id] = `Необходимо выбрать не более ${question.extra_params.maxAnswersCount} ответов`;
          }
          break;

        case QuestionType.ShortText:
        case QuestionType.LongText:
          if (
            question.extra_params.maxLength &&
            answer.length > question.extra_params.maxLength
          ) {
            errors[question.id] = `Максимальная длина ответа ${question.extra_params.maxLength} символов`;
          }
          break;

        case QuestionType.Consent:
          // Для consent дополнительная валидация не требуется (сервер уже проверяет true/false)
          break;

        case QuestionType.Date:
          // Для date валидация minDate/maxDate уже встроена в <input type="date">
          break;

        case QuestionType.Rating:
          // Для rating дополнительная валидация не требуется
          break;

        case QuestionType.Number:
          const numValue = parseFloat(answer);
          if (
            question.extra_params.minNumber &&
            numValue < question.extra_params.minNumber
          ) {
            errors[question.id] = `Число не должно быть меньше ${question.extra_params.minNumber}`;
          } else if (
            question.extra_params.maxNumber &&
            numValue > question.extra_params.maxNumber
          ) {
            errors[question.id] = `Число не должно превышать ${question.extra_params.maxNumber}`;
          }
          break;

        case QuestionType.Email:
          // Простая проверка формата email
          const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
          if (!emailRegex.test(answer)) {
            errors[question.id] = "Введите корректный email";
          }
          break;

        default:
          break;
      }
    });

    return errors;
  };

  const handleFinish = () => {
    const errors = validateQuestions();
    setValidationErrors(errors);

    // Если ошибок нет, завершаем интервью и перенаправляем
    if (Object.keys(errors).length === 0 && interviewId) {
      finishInterview(
        { hash, interviewId },
        {
          onSuccess: () => {
            router.push(`/poll/${hash}/success`);
          },
          onError: (error) => {
            console.error("Ошибка при завершении интервью:", error);
            // Можно добавить отображение ошибки пользователю, если нужно
          },
        }
      );
    }
  };

  if (idLoading || surveyLoading) return <div>Загрузка...</div>;
  if (idError || surveyError)
    return <div>Ошибка: {idError || surveyError}</div>;
  if (!data) return <div>Данные опроса недоступны</div>;

  return (

    <div className="max-w-2xl mx-auto p-4">
      {data.questions.map((question) => {
        const error = validationErrors[question.id];
        switch (question.type) {
          case QuestionType.SingleChoice:
            return (
              <SingleChoiceQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.MultiChoice:
            return (
              <MultipleOptionQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.ShortText:
            return (
              <ShortTextQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.LongText:
            return (
              <LongTextQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.Consent:
            return (
              <ConsentQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.Date:
            return (
              <DateQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.Rating:
            return (
              <RatingQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.Number:
            return (
              <NumberQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          case QuestionType.Email:
            return (
              <EmailQuestion
                key={question.id}
                question={question}
                onUpdate={(answer) => handleUpdateAnswer(question.id, answer)}
                error={error}
              />
            );
          default:
            return null;
        }
      })}
      <button
        className="mt-4 px-4 py-2 bg-blue-600 text-white rounded"
        onClick={handleFinish}
      >
        Завершить опрос
      </button>
    </div>
  );
}