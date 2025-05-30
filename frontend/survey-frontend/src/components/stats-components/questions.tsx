import React, { useState } from 'react';

import { QuestionStats } from '@/types/stats';

// Компонент для single_choice
export const SingleChoiceStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const answerCounts = question.answers.reduce((acc, ans) => {
    const numAns = Number(ans); // Парсим строку в число
    acc[numAns] = (acc[numAns] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      {question.options.map((option) => {
        const count = answerCounts[option.id] || 0;
        const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
        return (
          <div key={option.id} className="flex items-center space-x-3 mb-2">
            <span className="w-32">{option.label}</span>
            <div className="flex-1 bg-gray-200 rounded-full h-6">
              <div
                className="bg-blue-500 h-6 rounded-full text-xs text-white text-center leading-6"
                style={{ width: `${percentage}%` }}
              >
                {percentage.toFixed(1)}%
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

// Компонент для multi_choice
export const MultiChoiceStats = ({ question }: { question: QuestionStats }) => {
  const totalInterviews = question.answers.length;
  const answerCounts = question.answers.reduce((acc, ans) => {
    const parsedAns = JSON.parse(ans) as number[]; // Парсим строковый массив в числа
    parsedAns.forEach((id) => {
      acc[id] = (acc[id] || 0) + 1;
    });
    return acc;
  }, {} as Record<number, number>);

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      {question.options.map((option) => {
        const count = answerCounts[option.id] || 0;
        const percentage = totalInterviews > 0 ? (count / totalInterviews) * 100 : 0;
        return (
          <div key={option.id} className="flex items-center space-x-3 mb-2">
            <span className="w-32">{option.label}</span>
            <div className="flex-1 bg-gray-200 rounded-full h-6">
              <div
                className="bg-green-500 h-6 rounded-full text-xs text-white text-center leading-6"
                style={{ width: `${percentage}%` }}
              >
                {percentage.toFixed(1)}%
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};

// Компонент для consent
export const ConsentStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const trueCount = question.answers.filter((ans) => ans === 'true').length;
  const falseCount = totalAnswers - trueCount;
  const truePercentage = totalAnswers > 0 ? (trueCount / totalAnswers) * 100 : 0;
  const falsePercentage = totalAnswers > 0 ? (falseCount / totalAnswers) * 100 : 0;

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <div className="flex items-center space-x-3 mb-2">
        <span className="w-32">Согласны</span>
        <div className="flex-1 bg-gray-200 rounded-full h-6">
          <div
            className="bg-purple-500 h-6 rounded-full text-xs text-white text-center leading-6"
            style={{ width: `${truePercentage}%` }}
          >
            {truePercentage.toFixed(1)}%
          </div>
        </div>
      </div>
      <div className="flex items-center space-x-3 mb-2">
        <span className="w-32">Не согласны</span>
        <div className="flex-1 bg-gray-200 rounded-full h-6">
          <div
            className="bg-red-500 h-6 rounded-full text-xs text-white text-center leading-6"
            style={{ width: `${falsePercentage}%` }}
          >
            {falsePercentage.toFixed(1)}%
          </div>
        </div>
      </div>
    </div>
  );
};

export const TextStats = ({ question }: { question: QuestionStats }) => {
  const [showAll, setShowAll] = useState(false);

  // Подсчитываем частоту ответов
  const answerCounts = question.answers.reduce((acc, ans) => {
    acc[ans] = (acc[ans] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  // Сортируем по частоте и берем топ-5 (или первые 5, если уникальны)
  const sortedAnswers = Object.entries(answerCounts).sort((a, b) => b[1] - a[1]);
  const top5 = sortedAnswers.slice(0, 5);
  const remaining = sortedAnswers.slice(5);

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <ul className="list-disc pl-5">
        {top5.map(([answer, count]) => (
          <li key={answer}>
            {answer} ({count} раз)
          </li>
        ))}
      </ul>
      {remaining.length > 0 && (
        <button
          className="mt-2 text-blue-600 hover:underline"
          onClick={() => setShowAll(!showAll)}
        >
          {showAll ? 'Скрыть' : 'Показать остальные ответы'}
        </button>
      )}
      {showAll && (
        <ul className="list-disc pl-5 mt-2">
          {remaining.map(([answer, count]) => (
            <li key={answer}>
              {answer} ({count} раз)
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export const DateStats = ({ question }: { question: QuestionStats }) => {
  const dateCounts = question.answers.reduce((acc, ans) => {
    acc[ans] = (acc[ans] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const totalAnswers = question.answers.length;

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      {Object.entries(dateCounts).map(([date, count]) => {
        const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
        return (
          <div key={date} className="flex items-center space-x-3 mb-2">
            <span className="w-32">{date}</span>
            <div className="flex-1 bg-gray-200 rounded-full h-6">
              <div
                className="bg-blue-500 h-6 rounded-full text-xs text-white text-center leading-6"
                style={{ width: `${percentage}%` }}
              >
                {percentage.toFixed(1)}%
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
};


export const NumberStats = ({ question }: { question: QuestionStats }) => {
  const numbers = question.answers.map(Number);
  const totalAnswers = numbers.length;

  // Задаём интервалы (можно настроить под свои данные)
  const intervals = [0, 10, 20, 30, 40, 50];
  const frequency = intervals.map((start, i) => {
    const end = intervals[i + 1] || Infinity;
    return numbers.filter((num) => num >= start && num < end).length;
  });

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <div className="mt-2">
        {intervals.map((start, i) => {
          const end = intervals[i + 1] || 'и выше';
          const count = frequency[i];
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          return (
            <div key={start} className="flex items-center space-x-3 mb-2">
              <span className="w-32">{start} - {end}</span>
              <div className="flex-1 bg-gray-200 rounded-full h-6">
                <div
                  className="bg-green-500 h-6 rounded-full text-xs text-white text-center leading-6"
                  style={{ width: `${percentage}%` }}
                >
                  {percentage.toFixed(1)}%
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export const RatingStats = ({ question }: { question: QuestionStats }) => {
  const totalAnswers = question.answers.length;
  const sum = question.answers.reduce((acc, ans) => acc + Number(ans), 0);
  const average = totalAnswers > 0 ? (sum / totalAnswers).toFixed(1) : '0.0';

  // Подсчитываем частоту оценок (например, 1 звезда - 2 раза, 5 звёзд - 1 раз)
  const frequency = question.answers.reduce((acc, ans) => {
    const numAns = Number(ans);
    acc[numAns] = (acc[numAns] || 0) + 1;
    return acc;
  }, {} as Record<number, number>);

  const maxRating = 5; // Максимальный рейтинг (например, 5 звёзд)

  return (
    <div className="mb-6">
      <h3 className="font-semibold text-lg mb-2">{question.label}</h3>
      <p>Средний рейтинг: {average}</p>
      <div className="mt-2">
        {Array.from({ length: maxRating }, (_, i) => i + 1).map((star) => {
          const count = frequency[star] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          return (
            <div key={star} className="flex items-center space-x-3 mb-2">
              <span className="w-16">{star} звезда</span>
              <div className="flex-1 bg-gray-200 rounded-full h-6">
                <div
                  className="bg-yellow-500 h-6 rounded-full text-xs text-white text-center leading-6"
                  style={{ width: `${percentage}%` }}
                >
                  {percentage.toFixed(1)}%
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};