import { useQueryClient, useMutation } from "@tanstack/react-query";
import {
  createQuestion,
  updateQuestionLabel,
  updateQuestionType,
  updateQuestionOrder,
  updateQuestionExtraParams,
  restoreQuestion,
  deleteQuestion,
  QuestionType,
} from "@/api-client/question";
import {
  UpdateQuestionLabelRequest,
  UpdateQuestionOrderRequest,
  UpdateQuestionExtraParamsRequest,
} from "@/types/question";
import { GetSurveyResponse, SurveyQuestion } from "@/types/survey";
import { QuestionResponse } from "@/types/question";

// Query keys
const SURVEY_QUERY_KEY = (hash: string) => ["survey", hash];

// Hook для создания вопроса
export const useCreateQuestion = () => {
  const queryClient = useQueryClient();
  return useMutation<
    QuestionResponse,
    Error,
    { hash: string; data: { type: QuestionType } }
  >({
    mutationFn: async ({ hash, data }) => {
      const response = await createQuestion(hash, data);
      if (!response.success || !response.data) {
        throw new Error(response.error || "Failed to create question");
      }
      return response.data;
    },
    onMutate: async ({ hash }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      return { previousSurvey };
    },
    onSuccess: (data, { hash }) => {
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          const newQuestion: SurveyQuestion = {
            ...data.question,
            id: data.question.id,
           // survey_id: old.survey.id || 0, // Assume survey_id is available or adjust
            label: data.question.label || "",
            type: data.question.type,
            order: data.question.order,
            extra_params: data.question.extra_params || {},
            options: data.question.options || [],
            created_at: data.question.created_at,
            updated_at: data.question.updated_at,
          };
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: [...(old.survey.questions || []), newQuestion],
            },
          };
        }
      );
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для обновления метки вопроса
export const useUpdateQuestionLabel = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; data: UpdateQuestionLabelRequest }
  >({
    mutationFn: async ({ hash, questionId, data }) => {
      const response = await updateQuestionLabel(hash, questionId, data);
      if (!response.success) {
        throw new Error(response.error || "Failed to update question label");
      }
    },
    onMutate: async ({ hash, questionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, label: data.label } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для обновления типа вопроса
export const useUpdateQuestionType = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; data: { newType: QuestionType } }
  >({
    mutationFn: async ({ hash, questionId, data }) => {
      const response = await updateQuestionType(hash, questionId, data);
      if (!response.success) {
        throw new Error(response.error || "Failed to update question type");
      }
    },
    onMutate: async ({ hash, questionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, type: data.newType } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для обновления порядка вопроса
export const useUpdateQuestionOrder = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; data: UpdateQuestionOrderRequest }
  >({
    mutationFn: async ({ hash, questionId, data }) => {
      const response = await updateQuestionOrder(hash, questionId, data);
      if (!response.success) {
        throw new Error(response.error || "Failed to update question order");
      }
    },
    onMutate: async ({ hash, questionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, order: data.newOrder } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для обновления extra_params вопроса
export const useUpdateQuestionExtraParams = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; data: UpdateQuestionExtraParamsRequest }
  >({
    mutationFn: async ({ hash, questionId, data }) => {
      const response = await updateQuestionExtraParams(hash, questionId, data);
      if (!response.success) {
        throw new Error(
          response.error || "Failed to update question extra params"
        );
      }
    },
    onMutate: async ({ hash, questionId, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, extra_params: data } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для восстановления вопроса
export const useRestoreQuestion = () => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, { hash: string; questionId: number }>({
    mutationFn: async ({ hash, questionId }) => {
      const response = await restoreQuestion(hash, questionId);
      if (!response.success) {
        throw new Error(response.error || "Failed to restore question");
      }
    },
    onMutate: async ({ hash, questionId }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, state: "ACTUAL" } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};

// Hook для удаления вопроса
export const useDeleteQuestion = () => {
  const queryClient = useQueryClient();
  return useMutation<void, Error, { hash: string; questionId: number }>({
    mutationFn: async ({ hash, questionId }) => {
      const response = await deleteQuestion(hash, questionId);
      if (!response.success) {
        throw new Error(response.error || "Failed to delete question");
      }
    },
    onMutate: async ({ hash, questionId }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.filter(
                (q) => q.id !== questionId
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};
