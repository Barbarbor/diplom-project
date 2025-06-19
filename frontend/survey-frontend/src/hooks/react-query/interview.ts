import {useQuery, useMutation, useQueryClient} from '@tanstack/react-query';
import {
    startInterview,
    getSurveyWithAnswers,
    updateQuestionAnswer,
    UpdateAnswerRequest,
    finishInterview
} from '@/api-client/interview';
import {SurveyWithAnswers} from '@/types/interview';
import {CustomError} from '@/lib/error';

export const surveyWithAnswersQueryKey = (hash: string, interviewId: string) => [
    'surveyWithAnswers',
    hash,
    interviewId
];

export const useStartInterview = () => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, {hash: string; interviewId: string; isDemo?: string}>({
        mutationFn: async ({hash, interviewId, isDemo}) => {
            const response = await startInterview(hash, interviewId, isDemo);
            if (!response.success || !response.data) {
                throw new CustomError(
                    response.error || 'Failed to fetch survey stats',
                    response.status
                );
            }
            return response.data;
        },
        onSuccess: (_, {hash, interviewId}) => {
            // Инвалидируем данные опроса после начала интервью
            queryClient.invalidateQueries({queryKey: surveyWithAnswersQueryKey(hash, interviewId)});
        },
        retry: false,
        throwOnError: true
    });
};

export const useGetSurveyWithAnswers = (hash: string, interviewId: string) => {
    return useQuery<SurveyWithAnswers, Error>({
        queryKey: surveyWithAnswersQueryKey(hash, interviewId),
        queryFn: async () => {
            const response = await getSurveyWithAnswers(hash, interviewId);
            if (!response.success || !response.data) {
                throw new CustomError(
                    response.error || 'Failed to fetch survey stats',
                    response.status
                );
            }
            return response.data;
        },
        enabled: !!hash && !!interviewId, // Запрос выполняется только если hash и interviewId не пустые,
        retry: false,
        refetchOnMount: false,
        throwOnError: true
    });
};

export const useUpdateQuestionAnswer = () => {
    const queryClient = useQueryClient();
    return useMutation<
        void,
        Error,
        {hash: string; interviewId: string; questionId: number; data: UpdateAnswerRequest}
    >({
        mutationFn: ({hash, interviewId, questionId, data}) =>
            updateQuestionAnswer(interviewId, hash, questionId, data).then(res => res.data),
        onSuccess: (_, {hash, interviewId, questionId, data}) => {
            // Обновляем данные в кэше
            queryClient.setQueryData(
                surveyWithAnswersQueryKey(hash, interviewId),
                (oldData: SurveyWithAnswers | undefined) => {
                    if (!oldData) return oldData;
                    return {
                        ...oldData,
                        questions: oldData.questions.map(question =>
                            question.id === questionId
                                ? {...question, answer: data.answer}
                                : question
                        )
                    };
                }
            );
        }
    });
};

export const useFinishInterview = () => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, {hash: string; interviewId: string}>({
        mutationFn: ({hash, interviewId}) =>
            finishInterview(hash, interviewId).then(res => res.data),
        onSuccess: (_, {hash, interviewId}) => {
            // Инвалидируем данные опроса после завершения интервью
        }
    });
};
