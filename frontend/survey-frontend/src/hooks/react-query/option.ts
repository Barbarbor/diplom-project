import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  createOption,
  updateOptionLabel,
  updateOptionOrder,
  deleteOption,
} from '@/api-client/option';
import { GetSurveyResponse } from '@/types/survey';
import { OptionResponse, UpdateOptionLabelRequest, UpdateOptionOrderRequest } from '@/types/option';

// Ключ для кэширования данных опроса
const SURVEY_QUERY_KEY = (hash: string) => ['survey', hash];

// Хук для создания новой опции
export const useCreateOption = () => {
  const queryClient = useQueryClient();
  return useMutation<
    OptionResponse,
    Error,
    { hash: string; questionId: number }
  >({
    mutationFn: async ({ hash, questionId }) => {
      const response = await createOption(hash, questionId);
      if (!response.success || !response.data) {
        throw new Error(response.error || 'Failed to create option');
      }
      return response.data;
    },
    onSuccess: (data, { hash, questionId }) => {
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        const updatedQuestions = old.survey.questions.map((question) => {
          if (question.id === questionId) {
            return {
              ...question,
              options: [...(question.options || []), data.option],
            };
          }
          return question;
        });
        queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash)});
        return {
          ...old,
          survey: {
            ...old.survey,
            questions: updatedQuestions,
          },
        };
      });
    },
  });
};

// Хук для обновления лейбла опции
export const useUpdateOptionLabel = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; optionId: number; data: UpdateOptionLabelRequest }
  >({
    mutationFn: async ({ hash, questionId, optionId, data }) => {
      const response = await updateOptionLabel(hash, questionId, optionId, data);
      if (!response.success) {
        throw new Error(response.error || 'Failed to update option label');
      }
    },
    onMutate: async ({ hash, questionId, optionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        const updatedQuestions = old.survey.questions.map((question) => {
          if (question.id === questionId) {
            return {
              ...question,
              options: question?.options?.map((option) =>
                option.id === optionId ? { ...option, label: data.label } : option
              ),
            };
          }
          return question;
        });
        return {
          ...old,
          survey: {
            ...old.survey,
            questions: updatedQuestions,
          },
        };
      });
      
      return { previousSurvey };
    },
    onSuccess: (_, { hash }) => {
      queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash) }); // Инвалидация
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Хук для обновления порядка опции
export const useUpdateOptionOrder = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; optionId: number; data: UpdateOptionOrderRequest }
  >({
    mutationFn: async ({ hash, questionId, optionId, data }) => {
      const response = await updateOptionOrder(hash, questionId, optionId, data);
      if (!response.success) {
        throw new Error(response.error || 'Failed to update option order');
      }
    },
    onMutate: async ({ hash, questionId, optionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        const updatedQuestions = old.survey.questions.map((question) => {
          if (question.id === questionId) {
            return {
              ...question,
              options: question?.options?.map((option) =>
                option.id === optionId ? { ...option, order: data.new_order } : option
              ),
            };
          }
          return question;
        });
        return {
          ...old,
          survey: {
            ...old.survey,
            questions: updatedQuestions,
          },
        };
      });
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Хук для удаления опции
export const useDeleteOption = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; optionId: number }
  >({
    mutationFn: async ({ hash, questionId, optionId }) => {
      const response = await deleteOption(hash, questionId, optionId);
      if (!response.success) {
        throw new Error(response.error || 'Failed to delete option');
      }
    },
    onMutate: async ({ hash, questionId, optionId }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        const updatedQuestions = old.survey.questions.map((question) => {
          if (question.id === questionId) {
            return {
              ...question,
              options: question?.options?.filter((option) => option.id !== optionId),
            };
          }
          return question;
        });
        return {
          ...old,
          survey: {
            ...old.survey,
            questions: updatedQuestions,
          },
        };
      });

      return { previousSurvey };
    },
    onSuccess: (_, { hash }) => {
      queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash) }); // Инвалидация
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};