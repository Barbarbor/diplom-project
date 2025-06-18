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
  EmailStats,
} from '@/components/stats-components/questions';
import ExportToCsvButton from '@/components/stats-components/csvButton';
import { Block } from '@/components/common/Block';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import Spinner from '@/components/common/Spinner';

export default function SurveyStatsPage() {
  const params = useParams();
  const hash = params.hash as string;

  const { data, isLoading, error } = useGetSurveyStats(hash);
  console.log('data', data);

  if (isLoading) return <Spinner />;
  if (error) return <div>Ошибка: {error.message}</div>;
  if (!data) return <div>Данные статистики недоступны</div>;

  // Расчет процента завершения опросов
  const completionPercentage = data.started_interviews > 0
    ? ((data.completed_interviews / data.started_interviews) * 100).toFixed(1)
    : '0.0';

  return (
    <div className="max-w-4xl mx-auto p-4">
<Link href={`/survey/${hash}`} className="flex items-center space-x-2 text-blue-600 hover:underline mb-4">
        <ArrowLeftIcon className="h-5 w-5" />
        <span>Вернуться к опросу</span>
      </Link>
      {/* Блок общей статистики */}
      <Block>
        <h2 className="text-xl font-semibold mb-2">Общая статистика</h2>
        <div className="space-y-1">
          <p>Начато интервью: {data.started_interviews}</p>
          <p>Завершено интервью: {data.completed_interviews}</p>
          <p>Процент завершения: {completionPercentage}%</p>
          <p>Среднее время прохождения анкеты: {data.average_completion_time.toFixed(2)}s</p>
        </div>
        <div className="mt-4">
          <ExportToCsvButton data={data} />
        </div>
      </Block>

      {/* Блок результатов по вопросам */}
      <div >
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
            case 'email':
              return <EmailStats key={question.id} question={question}/>;
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