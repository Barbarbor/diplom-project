import { useGetSurvey } from "./react-query/survey";

export const useSurveyQuestion = (questionId: number) => {
  const { data: surveyData, isLoading, error } = useGetSurvey();

  const question = surveyData?.survey.questions.find((q) => q.id === questionId) || null;

  return {
    question,
    isLoading,
    error,
    refetch: () => {}, // Placeholder; useGetSurvey's refetch if needed
  };
};