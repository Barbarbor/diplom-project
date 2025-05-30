'use client';

import React, { useState } from 'react';
import { useGetSurvey, useRestoreSurvey } from '@/hooks/react-query/survey';
import { useCreateQuestion } from '@/hooks/react-query/question';
import { useUpdateSurvey } from '@/hooks/react-query/survey';
import { usePublishSurvey } from '@/hooks/react-query/survey';
import { QuestionType } from '@/api-client/question';
import { SurveyDetail } from '@/types/survey';
import { useParams } from 'next/navigation';
import QuestionBlock from '@/components/QuestionBlock';
import EditableLabel from '@/components/questions/EditableLabel';
import RestoreIcon from '@/components/common/RestoreIcon';

export default function SurveyPageClient() {
  const params = useParams();
  const hash = params.hash as string;
  const { data, isLoading } = useGetSurvey();
  const createQ = useCreateQuestion();
  const updateSurvey = useUpdateSurvey();
  const publishSurvey = usePublishSurvey();
  const restoreSurvey = useRestoreSurvey();

  const survey = data?.survey ?? ({} as SurveyDetail);
  const questions = data?.survey.questions ?? [];

  // State to control the visibility of the preview modal
  const [isPreviewOpen, setIsPreviewOpen] = useState(false);

  if (isLoading) {
    return <>Loading</>;
  }

  const handleTitleChange = (newTitle: string) => {
    updateSurvey.mutate({ hash, data: { title: newTitle } });
  };

  const handlePublish = () => {
    publishSurvey.mutate(hash);
  };

  const onAddSingle = () => {
    createQ.mutate({ hash, data: { type: QuestionType.SingleChoice } });
  };

  const handleRestoreSurvey = () => {
    restoreSurvey.mutate(hash);
  };

  // Function to open the preview modal
  const handlePreview = () => {
    setIsPreviewOpen(true);
  };

  // Function to close the preview modal
  const handleClosePreview = () => {
    setIsPreviewOpen(false);
  };

  return (
    <div className="flex min-h-screen">
      <div className="flex-1 flex justify-center p-6">
        <div className="max-w-3xl w-full">
          <div className="mb-6 p-4 bg-white border border-gray-300 rounded shadow-sm">
            <h2 className="text-xl font-semibold mb-4">Survey Details</h2>
            <div className="grid grid-cols-2 gap-4">
              <div className="font-medium">Name:</div>
              <div>
                <EditableLabel initialLabel={survey?.title} onLabelChange={handleTitleChange} />
              </div>
              <div className="font-medium">Author:</div>
              <div>{survey.creator || 'Unknown'}</div>
              <div className="font-medium">Created:</div>
              <div>
                {survey.created_at
                  ? new Date(survey.created_at).toLocaleString('ru-RU', {
                      timeZone: 'UTC',
                    })
                  : 'Unknown'}
              </div>
            </div>
            <button
              className="mt-4 px-4 py-2 bg-green-600 text-white rounded"
              onClick={handlePublish}
            >
              Publish Survey
            </button>
            <RestoreIcon onRestore={handleRestoreSurvey} entityType="survey" />
          </div>

          <ul className="space-y-4">
            {questions.map((q) => (
              <QuestionBlock key={q.id} question={q} />
            ))}
          </ul>
        </div>
      </div>

      <div className="fixed top-20 right-6 w-48 p-4 bg-white border border-gray-300 rounded shadow-md">
        <button
          className="w-full mb-2 px-4 py-2 bg-blue-600 text-white rounded"
          onClick={onAddSingle}
        >
          Добавить вопрос
        </button>
        <button
          className="w-full px-4 py-2 bg-gray-500 text-white rounded"
          onClick={handlePreview}
        >
          Предпросмотр
        </button>
      </div>

      {/* Preview Modal */}
      {isPreviewOpen && (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50">
          <div className="bg-white p-4 rounded shadow-md w-11/12 h-5/6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-semibold">Предпросмотр опроса</h2>
              <button
                className="px-4 py-2 bg-red-600 text-white rounded"
                onClick={handleClosePreview}
              >
                Закрыть
              </button>
            </div>
            <iframe
              src={`/poll/${hash}?isDemo=true`}
              className="w-full h-full border-0"
              title="Survey Preview"
            />
          </div>
        </div>
      )}
    </div>
  );
}