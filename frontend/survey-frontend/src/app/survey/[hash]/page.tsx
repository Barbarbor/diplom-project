'use client';
import { useGetSurvey, useUpdateSurvey, usePublishSurvey, useRestoreSurvey } from '@/hooks/react-query/survey';
import { useCreateQuestion, useUpdateQuestionLabel, useUpdateQuestionType, useUpdateQuestionOrder, useRestoreQuestion, useDeleteQuestion } from '@/hooks/react-query/question';
import { useParams } from 'next/navigation';
import { useState } from 'react';
import { QuestionType, SurveyQuestion } from '@/types/question';
import { SurveyState } from '@/types/survey';

export default function SurveyPage() {
  const params = useParams();
  const hash = params.hash as string;

  const { data, isLoading, error } = useGetSurvey(hash);


  const survey = data?.survey;

  const [title, setTitle] = useState(survey?.title);
  const updateSurveyMutation = useUpdateSurvey();

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setTitle(e.target.value);
  };

  const handleTitleSave = () => {
    updateSurveyMutation.mutate({ hash, data: { title } });
  };

  const publishSurveyMutation = usePublishSurvey();
  const restoreSurveyMutation = useRestoreSurvey();

  const handlePublish = () => {
    publishSurveyMutation.mutate(hash);
  };

  const handleRestore = () => {
    restoreSurveyMutation.mutate(hash);
  };

  const [newQuestionType, setNewQuestionType] = useState<QuestionType>(QuestionType.SingleChoice);
  const createQuestionMutation = useCreateQuestion();

  const handleCreateQuestion = () => {
    createQuestionMutation.mutate({ hash, data: { type: newQuestionType } });
  };

  const sortedQuestions = [...(survey?.questions ?? [])].sort((a, b) => a.order - b.order);

  return (
    <div>
      <h1>Survey: {survey?.title}</h1>
      <p>State: {survey?.state}</p>
      <p>Created by: {survey?.creator}</p>
      <p>Created at: {survey?.created_at}</p>
      <p>Updated at: {survey?.updated_at}</p>

      <div>
        <input value={title} onChange={handleTitleChange} placeholder="Survey Title" />
        <button onClick={handleTitleSave}>Save Title</button>
      </div>

      {survey?.state === SurveyState.Draft && (
        <button onClick={handlePublish}>Publish Survey</button>
      )}
      {/* Add condition for restore button if applicable */}

      <h2>Questions</h2>
      {sortedQuestions.map((question) => (
        <Question key={question.id} question={question} hash={hash} />
      ))}

      <div>
        <select
          value={newQuestionType}
          onChange={(e) => setNewQuestionType(e.target.value as QuestionType)}
        >
          {Object.values(QuestionType).map((type) => (
            <option key={type} value={type}>
              {type}
            </option>
          ))}
        </select>
        <button onClick={handleCreateQuestion}>Add Question</button>
      </div>
    </div>
  );
}

interface QuestionProps {
  question: SurveyQuestion;
  hash: string;
}

function Question({ question, hash }: QuestionProps) {
  const [label, setLabel] = useState(question.label);
  const [type, setType] = useState(question.type);
  const [order, setOrder] = useState(question.order.toString());

  const updateLabelMutation = useUpdateQuestionLabel();
  const updateTypeMutation = useUpdateQuestionType();
  const updateOrderMutation = useUpdateQuestionOrder();
  const restoreMutation = useRestoreQuestion();
  const deleteMutation = useDeleteQuestion();

  const handleLabelChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setLabel(e.target.value);
  };

  const handleLabelSave = () => {
    updateLabelMutation.mutate({ hash, questionId: question.id, data: { label } });
  };

  const handleTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setType(e.target.value as QuestionType);
  };

  const handleTypeSave = () => {
    updateTypeMutation.mutate({ hash, questionId: question.id, data: { newType: type } });
  };

  const handleOrderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setOrder(e.target.value);
  };

  const handleOrderSave = () => {
    const newOrder = parseInt(order, 10);
    if (!isNaN(newOrder)) {
      updateOrderMutation.mutate({ hash, questionId: question.id, data: { newOrder } });
    }
  };

  const handleRestore = () => {
    restoreMutation.mutate({ hash, questionId: question.id });
  };

  const handleDelete = () => {
    deleteMutation.mutate({ hash, questionId: question.id });
  };

  return (
    <div>
      <input value={label} onChange={handleLabelChange} placeholder="Question Label" />
      <button onClick={handleLabelSave}>Save Label</button>

      <select value={type} onChange={handleTypeChange}>
        {Object.values(QuestionType).map((t) => (
          <option key={t} value={t}>
            {t}
          </option>
        ))}
      </select>
      <button onClick={handleTypeSave}>Save Type</button>

      <input type="number" value={order} onChange={handleOrderChange} placeholder="Order" />
      <button onClick={handleOrderSave}>Save Order</button>

      <button onClick={handleRestore}>Restore</button>
      <button onClick={handleDelete}>Delete</button>
    </div>
  );
}