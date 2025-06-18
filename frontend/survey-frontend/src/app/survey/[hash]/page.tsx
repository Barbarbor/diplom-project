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
import { SurveyActions } from './surveyActions'; // Импортируем новый компонент
import { AccessModal } from './accessModal'; // Импортируем новый компонент
import { PreviewModal } from './previewModal';
import { Block } from '@/components/common/Block';
import { SurveyDistribution } from './surveyDistribution'; // Импортируем новый компонент
import Spinner from '@/components/common/Spinner';


const SURVEY_STATE = {'DRAFT': 'Черновик', 'ACTIVE': 'Активный'}

export default function SurveyPageClient() {
  const params = useParams();
  const hash = params.hash as string;
  const { data, isLoading } = useGetSurvey();
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
    <div className="flex min-h-screen">
      <div className="flex-1 flex justify-center p-6">
        <div className="max-w-3xl w-full">
          <Block>
            <div className="grid grid-cols-2 gap-4">
              <div className="font-medium">Название опроса</div>
              <div>
                <EditableLabel initialLabel={survey?.title} onLabelChange={handleTitleChange} />
              </div>
              <div className="font-medium">Автор</div>
              <div>{survey.creator || 'Unknown'}</div>
              <div className="font-medium">Дата создания</div>
              <div>
                {survey.created_at
                  ? new Date(survey.created_at).toLocaleString('ru-RU', {
                      timeZone: 'UTC',
                    })
                  : 'Unknown'}
              </div>
              <div className="font-medium">Состояние</div>
              <div>{SURVEY_STATE[survey.state]}</div>
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
           {[...questions].sort((a, b) => (a.question_order || 0) - (b.question_order || 0)).map((q) => (
  <QuestionBlock key={q.id} question={q} />
))}
          </ul>
        </div>
      </div>

      {/* Модальные окна */}
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