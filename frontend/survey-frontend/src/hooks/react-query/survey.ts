import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getSurvey,
  updateSurvey,
  publishSurvey,
  restoreSurvey,
  UpdateSurveyRequest,
  getSurveyStats,
  getSurveys,
} from '@/api-client/survey';
import { GetSurveyResponse, GetSurveysResponse } from '@/types/survey';
import { useSurveyHash } from '../survey';
import { SurveyStats } from '@/types/stats';

// Query keys
export const SURVEY_QUERY_KEY = (hash: string) => ['survey', hash];

export const useGetSurveys = () => {
  return useQuery<GetSurveysResponse, Error>({
    queryKey: ["surveys"],
    queryFn: async () => {
      const response = await getSurveys();
      if (!response.success || !response.data) {
        throw new Error(response.error || 'Failed to fetch survey');
      }
      return response.data;
    },

  });
};

// Hook для получения опроса
export const useGetSurvey = () => {
  const hash = useSurveyHash();
  return useQuery<GetSurveyResponse, Error>({
    queryKey: SURVEY_QUERY_KEY(hash),
    queryFn: async () => {
      const response = await getSurvey(hash);
      if (!response.success || !response.data) {
        throw new Error(response.error || 'Failed to fetch survey');
      }
      return response.data;
    },

  });
};

// Hook для обновления опроса
export const useUpdateSurvey = () => {
  const queryClient = useQueryClient();
  return useMutation<
    void,
    Error,
    { hash: string; data: UpdateSurveyRequest }
  >({
    mutationFn: async ({ hash, data }) => {
      const response = await updateSurvey(hash, data);
      if (!response.success) {
        throw new Error(response.error || 'Failed to update survey');
      }
    },
    onMutate: async ({ hash, data }) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        return {
          ...old,
          survey: { ...old.survey, title: data.title },
        };
      });
      return { previousSurvey };
    },
    // onError: (err, { hash }, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
    // onSuccess: (_, { hash, data }) => {
    //   // No need to update cache, as onMutate already set the title
    // },
  });
};

// Hook для публикации опроса
export const usePublishSurvey = () => {
  const queryClient = useQueryClient();
  const hash = useSurveyHash();
  return useMutation<void, Error, string>({
    mutationFn: async (hash) => {
      const response = await publishSurvey(hash);
      if (!response.success) {
        throw new Error(response.error || 'Failed to publish survey');
      }
    },
    onMutate: async (hash) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        return {
          ...old,
          survey: { ...old.survey, state: 'ACTIVE' },
        };
      });
      return { previousSurvey };
    },
    // onError: (err, hash, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      // No need to update cache
    },
  });
};

// Hook для восстановления опроса
export const useRestoreSurvey = () => {
  const queryClient = useQueryClient();
    const hash = useSurveyHash();
  return useMutation<void, Error, string>({
    mutationFn: async (hash) => {
      const response = await restoreSurvey(hash);
      if (!response.success) {
        throw new Error(response.error || 'Failed to restore survey');
      }
    },
    onMutate: async (hash) => {
      await queryClient.cancelQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
      const previousSurvey = queryClient.getQueryData<GetSurveyResponse>(SURVEY_QUERY_KEY(hash));
      queryClient.setQueryData(SURVEY_QUERY_KEY(hash), (old: GetSurveyResponse | undefined) => {
        if (!old) return old;
        return {
          ...old,
          survey: { ...old.survey, state: 'DRAFT' },
        };
      });
      return { previousSurvey };
    },
    // onError: (err, hash, context) => {
    //   queryClient.setQueryData(SURVEY_QUERY_KEY(hash), context?.previousSurvey);
    // },
    onSuccess: () => {
      // No need to update cache
      queryClient.refetchQueries({ queryKey: SURVEY_QUERY_KEY(hash) });
    },
  });
};

export function useGetSurveyStats(hash: string) {
  return useQuery<SurveyStats, Error>({
    queryKey: ['surveyStats', hash],
     queryFn: async () => {
      const response = await getSurveyStats(hash);
      if (!response.success || !response.data) {
        throw new Error(response.error || 'Failed to fetch survey stats');
      }
      return response.data;
    },
    enabled: !!hash, // Запрос выполняется только если hash существует
  });
}