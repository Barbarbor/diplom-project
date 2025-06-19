'use client';

import React, { useState } from 'react';
import { useGetSurvey, useRestoreSurvey } from '@/hooks/react-query/survey';
import { useCreateQuestion } from '@/hooks/react-query/question';
import { useUpdateSurvey } from '@/hooks/react-query/survey';
import { QuestionType } from '@/api-client/question';
import { SurveyDetail } from '@/types/survey';
import { useParams } from 'next/navigation';
import QuestionBlock from '@/components/QuestionBlock';
import EditableLabel from '@/components/questions/EditableLabel';
import { SurveyActions } from './surveyActions';
import { AccessModal } from './accessModal';
import { PreviewModal } from './previewModal';
import { Block } from '@/components/common/Block';
import { SurveyDistribution } from './surveyDistribution';
import Spinner from '@/components/common/Spinner';
import { useTranslation } from 'next-i18next';
import { FaArrowRight } from 'react-icons/fa';
import Link from 'next/link';

export default function SurveyPageClient() {
  // Use the 'survey' prefix for translations
  const { t, i18n } = useTranslation('translation', { keyPrefix: 'survey' });
  const params = useParams();
  const hash = params.hash as string;
  const { data, isLoading, error } = useGetSurvey();

  if (error) {
    throw new Error('Failed to load survey');
  }

  const createQ = useCreateQuestion();
  const updateSurvey = useUpdateSurvey();
  const restoreSurvey = useRestoreSurvey();

  const survey = data?.survey ?? ({} as SurveyDetail);
  const questions = data?.survey.questions ?? [];

  // State to control the visibility of modals
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);
  const [isAccessModalOpen, setIsAccessModalOpen] = useState(false);

  if (isLoading) {
    return <Spinner />;
  }

  const handleTitleChange = (newTitle: string) => {
    updateSurvey.mutate({ hash, data: { title: newTitle } });
  };

  const onAddQuestion = (type: QuestionType) => {
    createQ.mutate({ hash, data: { type } });
  };

  const handleRestoreSurvey = () => {
    restoreSurvey.mutate(hash);
  };

  return (
    <div className="flex min-h-screen relative">
      {/* "Go to Statistics" link in the top-left corner */}
      <div className="absolute top-4 left-4 z-10">
        <Link
          href={`/survey/${hash}/stats`}
          className="flex items-center text-blue-600 hover:text-blue-800"
        >
          <FaArrowRight className="mr-2" />
          {t('goToStats')}
        </Link>
      </div>

      <div className="flex-1 flex justify-center p-6">
        <div className="max-w-3xl w-full">
          <Block>
            <div className="grid grid-cols-2 gap-4">
              <div className="font-medium">{t('surveyTitle')}</div>
              <div>
                <EditableLabel
                  initialLabel={survey?.title}
                  onLabelChange={handleTitleChange}
                />
              </div>
              <div className="font-medium">{t('creator')}</div>
              <div>{survey.creator || t('unknown')}</div>
              <div className="font-medium">{t('createdAt')}</div>
              <div>
                {survey.created_at
                  ? new Date(survey.created_at).toLocaleString(i18n.language, {
                      timeZone: 'UTC',
                    })
                  : t('unknown')}
              </div>
              <div className="font-medium">{t('state')}</div>
              <div>{t(`states.${survey.state}`)}</div>
            </div>
            <SurveyActions
              state={survey.state}
              hash={hash}
              questionCount={questions.length}
              onAddQuestion={onAddQuestion}
              onOpenAccessModal={() => setIsAccessModalOpen(true)}
              onOpenPreviewModal={() => setIsPreviewOpen(true)}
              onRestoreSurvey={handleRestoreSurvey}
            />
          </Block>

          <SurveyDistribution hash={hash} />

          <ul className="space-y-4">
            {[...questions]
              .sort((a, b) => (a.question_order || 0) - (b.question_order || 0))
              .map((q) => (
                <QuestionBlock key={q.id} question={q} />
              ))}
          </ul>
        </div>
      </div>

      {/* Modals */}
      <AccessModal
        isOpen={isAccessModalOpen}
        onClose={() => setIsAccessModalOpen(false)}
      />
      <PreviewModal
        isOpen={isPreviewOpen}
        onClose={() => setIsPreviewOpen(false)}
        hash={hash}
      />
    </div>
  );
}