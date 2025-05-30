'use client';

import { useParams } from 'next/navigation';
import { useGetSurveyStats } from '@/hooks/react-query/survey';
import {
  SingleChoiceStats,
  ConsentStats,
  RatingStats,
  TextStats,
  DateStats,
  NumberStats,
  MultiChoiceStats,
} from '@/components/stats-components/questions';

export default function SurveyStatsPage() {
  const params = useParams();
  const hash = params.hash as string;

  const { data, isLoading, error } = useGetSurveyStats(hash);
  console.log('data', data);

  if (isLoading) return <div>Загрузка...</div>;
  if (error) return <div>Ошибка: {error.message}</div>;
  if (!data) return <div>Данные статистики недоступны</div>;

  const exportToCsv = () => {
    const csvData: string[][] = [];

    // Общая статистика
    csvData.push(['Общая статистика']);
    csvData.push(['Начато интервью', data.started_interviews.toString()]);
    csvData.push(['Завершено интервью', data.completed_interviews.toString()]);
    csvData.push([]);

    // Статистика по вопросам
    data.questions.forEach((question) => {
      csvData.push([`Вопрос: ${question.label}`, `Тип: ${question.type}`]);

      if (question.type === 'single_choice') {
        const totalAnswers = question.answers.length;
        const answerCounts = question.answers.reduce((acc, ans) => {
          const numAns = Number(ans);
          acc[numAns] = (acc[numAns] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([option.label, `${count} раз`, `${percentage.toFixed(1)}%`]);
        });
      } else if (question.type === 'multi_choice') {
        const totalInterviews = question.answers.length;
        const answerCounts = question.answers.reduce((acc, ans) => {
          const parsedAns = JSON.parse(ans) as number[];
          parsedAns.forEach((id) => {
            acc[id] = (acc[id] || 0) + 1;
          });
          return acc;
        }, {} as Record<number, number>);

        question?.options?.forEach((option) => {
          const count = answerCounts[option.id] || 0;
          const percentage = totalInterviews > 0 ? (count / totalInterviews) * 100 : 0;
          csvData.push([option.label, `${count} раз`, `${percentage.toFixed(1)}%`]);
        });
      } else if (question.type === 'consent') {
        const totalAnswers = question.answers.length;
        const trueCount = question.answers.filter((ans) => ans === 'true').length;
        const falseCount = totalAnswers - trueCount;
        const truePercentage = totalAnswers > 0 ? (trueCount / totalAnswers) * 100 : 0;
        const falsePercentage = totalAnswers > 0 ? (falseCount / totalAnswers) * 100 : 0;
        csvData.push(['Согласны', `${trueCount} раз`, `${truePercentage.toFixed(1)}%`]);
        csvData.push(['Не согласны', `${falseCount} раз`, `${falsePercentage.toFixed(1)}%`]);
      } else if (question.type === 'rating') {
        const totalAnswers = question.answers.length;
        const sum = question.answers.reduce((acc, ans) => acc + Number(ans), 0);
        const average = totalAnswers > 0 ? (sum / totalAnswers).toFixed(1) : '0.0';
        csvData.push(['Средний рейтинг', average]);

        const frequency = question.answers.reduce((acc, ans) => {
          const numAns = Number(ans);
          acc[numAns] = (acc[numAns] || 0) + 1;
          return acc;
        }, {} as Record<number, number>);

        for (let star = 1; star <= 5; star++) {
          const count = frequency[star] || 0;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([`${star} звезда`, `${count} раз`, `${percentage.toFixed(1)}%`]);
        }
      } else if (question.type === 'short_text' || question.type === 'long_text') {
        const answerCounts = question.answers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        const sortedAnswers = Object.entries(answerCounts).sort((a, b) => b[1] - a[1]);
        const top5 = sortedAnswers.slice(0, 5);
        top5.forEach(([answer, count]) => {
          csvData.push([answer, `${count} раз`]);
        });
      } else if (question.type === 'date') {
        const dateCounts = question.answers.reduce((acc, ans) => {
          acc[ans] = (acc[ans] || 0) + 1;
          return acc;
        }, {} as Record<string, number>);

        Object.entries(dateCounts).forEach(([date, count]) => {
          csvData.push([date, `${count} раз`]);
        });
      } else if (question.type === 'number') {
        const numbers = question.answers.map(Number);
        const totalAnswers = numbers.length;
        const intervals = [0, 10, 20, 30, 40, 50];
        intervals.forEach((start, i) => {
          const end = intervals[i + 1] || 'и выше';
          const count = numbers.filter((num) => num >= start && num < (end === 'и выше' ? Infinity : end)).length;
          const percentage = totalAnswers > 0 ? (count / totalAnswers) * 100 : 0;
          csvData.push([`${start} - ${end}`, `${count} раз`, `${percentage.toFixed(1)}%`]);
        });
      }

      csvData.push([]); // Пустая строка между вопросами
    });

    // Генерация CSV строки
    const csv = csvData.map(row => row.join(';')).join('\n');

    // Добавляем BOM для UTF-8
    const BOM = '\uFEFF'; // UTF-8 BOM
    const csvWithBom = BOM + csv;

    // Создание и скачивание файла
    const blob = new Blob([csvWithBom], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = `survey_stats_${new Date().toISOString()}.csv`;
    link.click();
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Статистика опроса</h1>
      <div className="mb-4">
        <p>Начато интервью: {data.started_interviews}</p>
        <p>Завершено интервью: {data.completed_interviews}</p>
      </div>
      <button
        className="mb-4 px-4 py-2 bg-blue-600 text-white rounded"
        onClick={exportToCsv}
      >
        Экспорт в CSV
      </button>
      <div>
        <h2 className="text-xl font-semibold mb-2">Вопросы</h2>
        {data.questions.map((question) => {
          switch (question.type) {
            case 'single_choice':
              return <SingleChoiceStats key={question.id} question={question} />;
            case 'multi_choice':
              return <MultiChoiceStats key={question.id} question={question} />;
            case 'consent':
              return <ConsentStats key={question.id} question={question} />;
            case 'rating':
              return <RatingStats key={question.id} question={question} />;
            case 'short_text':
            case 'long_text':
              return <TextStats key={question.id} question={question} />;
            case 'date':
              return <DateStats key={question.id} question={question} />;
            case 'number':
              return <NumberStats key={question.id} question={question} />;
            default:
              return (
                <div key={question.id} className="mb-4 p-4 border rounded">
                  <p className="font-medium">{question.label}</p>
                  <p>Тип: {question.type}</p>
                  <p>Ответы: {question.answers.join(', ')}</p>
                </div>
              );
          }
        })}
      </div>
    </div>
  );
}