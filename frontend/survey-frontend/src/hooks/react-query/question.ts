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
import { SURVEY_QUERY_KEY } from "./survey";



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
            question_order: data.question.question_order,
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


export const useUpdateStateBadge = () => {
  const queryClient = useQueryClient();

  return useMutation<
    void,
    Error,
    { hash: string; questionId: number; newState: string }
  >({
    mutationFn: async ({ hash, questionId, newState }) => {
      // Проверяем текущее состояние из кэша
      const currentSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      const currentQuestion = currentSurvey?.survey.questions.find(
        (q) => q.id === questionId
      );
      console.log('curr q', currentQuestion, 'newstate', newState)
      // Логика: не обновляем, если новое состояние совпадает с текущим или равно "NEW"
      if (
        currentQuestion?.question_state === newState ||
         currentQuestion?.question_state === "NEW"
      ) {
        return; // Ничего не делаем
      }
    },
    onMutate: async ({ hash, questionId, newState }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );

      // Оптимистичное обновление кэша
      if (previousSurvey) {
        queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, question_state: newState } : q
              ),
            },
          };
        });
      }

      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   if (context?.previousSurvey) {
    //     queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context.previousSurvey);
    //   }
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
    onSuccess: (_, { hash }) => {
    queryClient.refetchQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};
export const useUpdateQuestionType = () => {
  const queryClient = useQueryClient();

  return useMutation<
    SurveyQuestion, // Возвращаем обновленный вопрос
    Error,
    { hash: string; questionId: number; data: { newType: QuestionType } }
  >({
    mutationFn: async ({ hash, questionId, data }) => {
      const response = await updateQuestionType(hash, questionId, data);
      if (!response.success || !response.data) {
        throw new Error(response.error || "Failed to update question type");
      }
      return response.data; // Возвращаем обновленный вопрос
    },
    onMutate: async ({ hash }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });

    },
    onSuccess: (question, { hash, questionId }) => {
      // @ts-expect-error Property 'data' does not exist on type 'SurveyQuestion'.
      const updatedQuestion = question.data;
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? updatedQuestion : q
              ),
            },
          };
        }
      );
      queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash) })
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
      // Update the data object to use new_order instead of newOrder
      const response = await updateQuestionOrder(hash, questionId, {
        new_order: data.new_order,
      });

      if (!response.success) {
        throw new Error(response.error || "Failed to update question order");
      }
      return; // Успешное завершение
    },
    onMutate: async ({ hash, questionId, data }) => {
      // Cancel any ongoing queries for the survey
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });

      // Store the previous survey state for rollback
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );

      // Optimistic update (optional, can be removed if not needed)
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === questionId ? { ...q, question_order: data.new_order } : q
              ),
            },
          };
        }
      );
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   console.log('Mutation error:', err.message);
    //   // Rollback to the previous state if the mutation fails
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
    onSuccess: (_, { hash }) => {
      // Refetch the survey query to get updated data from the server
      queryClient.refetchQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
    },
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
    onSuccess: (_, { hash }) => {
      queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash) })
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};
export const useRestoreQuestion = () => {
  const queryClient = useQueryClient();

  return useMutation<
    SurveyQuestion, // Теперь возвращаем объект вопроса
    Error,
    { hash: string; questionId: number }
  >({
    mutationFn: async ({ hash, questionId }) => {
      const response = await restoreQuestion(hash, questionId);
      if (!response.success || !response.data) {
        throw new Error(response.error || "Failed to restore question");
      }
      return response.data;
    },
    onMutate: async ({ hash}) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
      return { previousSurvey };
    },
    onSuccess: (restoredQuestion, { hash }) => {
      // @ts-expect-error Property 'question' does not exist on type 'SurveyQuestion'.
      const question = restoredQuestion.question;
      // Обновляем кэш данными из ответа
      queryClient.setQueryData(
        SURVEY_QUERY_KEY(hash),
        (old: GetSurveyResponse | undefined) => {
          if (!old) return old;
          return {
            ...old,
            survey: {
              ...old.survey,
              questions: old.survey.questions.map((q) =>
                q.id === question.id ? question : q
              ),
            },
          };
        }
        
      );
       queryClient.refetchQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
    },
    // onError: (err, { hash }, context) => {
    //   if (context?.previousSurvey) {
    //     queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context.previousSurvey);
    //   }
    })
};

// Hook для удаления вопроса
export const useDeleteQuestion = () => {
  const queryClient = useQueryClient();
   const updateStateBadge = useUpdateStateBadge();
  return useMutation<void, Error, { hash: string; questionId: number }>({
    mutationFn: async ({ hash, questionId }) => {
      const response = await deleteQuestion(hash, questionId);
      if (!response.success) {
        throw new Error(response.error || "Failed to delete question");
      }
    },
    onMutate: async ({ hash }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(
        SURVEY_QUERY_KEY(hash)
      );
     
      return { previousSurvey };
      
    },
     onSuccess: (_, { hash, questionId }) => {
      updateStateBadge.mutate({ hash, questionId, newState: "DELETED" });
      queryClient.invalidateQueries({ queryKey: SURVEY_QUERY_KEY(hash) })
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
  });
};
