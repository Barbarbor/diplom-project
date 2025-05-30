// app/survey/[hash]/stats/page.tsx
'use client';

import { useParams } from 'next/navigation';
import { useGetSurveyStats } from '@/hooks/react-query/survey';
import {
SingleChoiceStats ,
  ConsentStats,
  RatingStats,
  TextStats,
  DateStats,
  NumberStats,
  MultiChoiceStats,
} from '@/components/stats-components/questions'; // Предполагаемый путь к компонентам

export default function SurveyStatsPage() {
  const params = useParams();
  const hash = params.hash as string;

  const { data, isLoading, error } = useGetSurveyStats(hash);
  console.log('data', data);

  if (isLoading) return <div>Загрузка...</div>;
  if (error) return <div>Ошибка: {error.message}</div>;
  if (!data) return <div>Данные статистики недоступны</div>;

  return (
    <div className="max-w-4xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Статистика опроса</h1>
      <div className="mb-4">
        <p>Начато интервью: {data.started_interviews}</p>
        <p>Завершено интервью: {data.completed_interviews}</p>
      </div>
      <div>
        <h2 className="text-xl font-semibold mb-2">Вопросы</h2>
        {data.questions.map((question) => {
          switch (question.type) {
            case 'single_choice':
              return <SingleChoiceStats key={question.id} question={question} />;
              case 'multi_choice':
                return <MultiChoiceStats key={question.id} question={question} /> ;
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