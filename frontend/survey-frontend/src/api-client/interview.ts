import request, { ApiResponse } from "@/lib/api";
import { SurveyWithAnswers } from "@/types/interview";

export const startInterview = async (
  hash: string,
  interviewId: string,
  isDemo?: string
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "POST",
    prefix: "/api",
    url: `/interview/${hash}/start?interviewId=${interviewId}${isDemo === 'true' ? `&isDemo=${true}` : ''}`,
  });
};

export const getSurveyWithAnswers = async (
  hash: string,
  interviewId: string
): Promise<ApiResponse<SurveyWithAnswers>> => {
  return await request<SurveyWithAnswers>({
    method: "GET",
    prefix: "/api",
    url: `/interview/${hash}/survey?interviewId=${interviewId}`,
  });
};

export interface UpdateAnswerRequest {
  answer: string;
}

export const updateQuestionAnswer = async (
  interviewId: string,
  hash: string,
  questionId: number,
  data: UpdateAnswerRequest
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/interview/${hash}/${questionId}/answer?interviewId=${interviewId}`,
    data,
  });
};

export const finishInterview = async (
  hash: string,
  interviewId: string
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "POST",
    prefix: "/api",
    url: `/interview/${hash}/finish?interviewId=${interviewId}`,
  });
};
