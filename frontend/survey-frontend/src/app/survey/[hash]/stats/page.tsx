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
import { FaArrowLeft } from 'react-icons/fa'; // Using FaArrowLeft for consistency
import Link from 'next/link';
import Spinner from '@/components/common/Spinner';
import { useTranslation } from 'next-i18next';

export default function SurveyStatsPage() {
  const { t } = useTranslation('translation', { keyPrefix: 'survey.stats' });
  const params = useParams();
  const hash = params.hash as string;

  const { data, isLoading, error } = useGetSurveyStats(hash);

  if (isLoading) return <Spinner />;
  if (error) return <div>{t('error', { message: error.message })}</div>;
  if (!data) return <div>{t('noData')}</div>;

  // Расчет процента завершения опросов
  const completionPercentage = data.started_interviews > 0
    ? ((data.completed_interviews / data.started_interviews) * 100).toFixed(1)
    : '0.0';

  return (
<div className="max-w-4xl mx-auto p-4 relative"> {/* Add relative positioning */}
 {/* Flex container with arrow and block */}
  <div className="flex flex-col">
    <div className="mb-4">
      <Link
        href={`/survey/${hash}`}
        className="flex items-center text-blue-600 hover:text-blue-800"
      >
        <FaArrowLeft className="mr-2" />
        {t('backToSurvey')}
      </Link>
    </div>
    <div>
      {/* Блок общей статистики with added margin-top */}
      <Block className="mt-4">
        <h2 className="text-xl font-semibold mb-2">{t('generalStatsTitle')}</h2>
        <div className="space-y-1">
          <p>{t('startedInterviews', { count: data.started_interviews })}</p>
          <p>{t('completedInterviews', { count: data.completed_interviews })}</p>
          <p>{t('completionPercentage', { percentage: completionPercentage })}</p>
          <p>{t('averageCompletionTime', { time: data.average_completion_time.toFixed(2) })}</p>
        </div>
        <div className="mt-4">
          <ExportToCsvButton data={data} />
        </div>
      </Block>
    </div>
  </div>



      {/* Блок результатов по вопросам */}
      <div>
        {data?.questions?.map((question) => {
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
              return <EmailStats key={question.id} question={question} />;
            case 'date':
              return <DateStats key={question.id} question={question} />;
            case 'number':
              return <NumberStats key={question.id} question={question} />;
            default:
              return (
                <div key={question.id} className="mb-4 p-4 border rounded">
                  <p className="font-medium">{question.label}</p>
                  <p>{t('type', { type: question.type })}</p>
                  <p>{t('answers', { answers: question.answers?.join(', ') || 'N/A' })}</p>
                </div>
              );
          }
        })}
      </div>
    </div>
  );
}