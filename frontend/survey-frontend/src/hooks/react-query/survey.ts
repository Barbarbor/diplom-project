import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getSurvey,
  updateSurvey,
  publishSurvey,
  restoreSurvey,
  UpdateSurveyRequest,
} from '@/api-client/survey';
import { GetSurveyResponse } from '@/types/survey';

// Query keys
const SURVEY_QUERY_KEY = (hash: string) => ['survey', hash];

// Hook для получения опроса
export const useGetSurvey = (hash: string) => {
  return useQuery<GetSurveyResponse, Error>({
    queryKey: SURVEY_QUERY_KEY(hash),
    queryFn: async () => {
      const response = await getSurvey(hash);
      if (!response.success || !response.data) {
        throw new Error(response.error || 'Failed to fetch survey');
      }
      return response.data;
    },
    enabled: !!hash,
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
      // No need to update cache
    },
  });
};

// Hook для восстановления опроса
export const useRestoreSurvey = () => {
  const queryClient = useQueryClient();
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
    },
  });
};