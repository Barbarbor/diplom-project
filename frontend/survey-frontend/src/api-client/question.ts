import request, { ApiResponse } from "@/lib/api";
import {
  QuestionType,
  QuestionResponse,
  UpdateQuestionLabelRequest,
  UpdateQuestionOrderRequest,
  UpdateQuestionExtraParamsRequest,
} from "@/types/question";

/**
 * Создать новый вопрос в опросе (POST /api/surveys/:hash/question).
 * @returns {question: SurveyQuestion}
 */
export const createQuestion = async (
  hash: string,
  data: { type: QuestionType }
): Promise<ApiResponse<QuestionResponse>> => {
  return await request<QuestionResponse>({
    method: "POST",
    prefix: "/api",
    url: `/surveys/${hash}/question${data?.type ? `?type=${data.type}`: ''}`
  });
};

/**
 * Обновить только метку вопроса (PATCH /api/surveys/:hash/question/:questionId).
 */
export const updateQuestionLabel = async (
  hash: string,
  questionId: number,
  data: UpdateQuestionLabelRequest
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}`,
    data,
  });
};

/**
 * Обновить тип вопроса (PATCH /api/surveys/:hash/question/:questionId/type).
 */
export const updateQuestionType = async (
  hash: string,
  questionId: number,
  data: { newType: QuestionType }
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/type`,
    data,
  });
};

/**
 * Обновить порядок вопроса (PATCH /api/surveys/:hash/question/:questionId/order).
 */
export const updateQuestionOrder = async (
  hash: string,
  questionId: number,
  data: UpdateQuestionOrderRequest
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/order`,
    data,
  });
};

/**
 * Обновить extra_params вопроса (PATCH /api/surveys/:hash/question/:questionId/extra_params).
 */
export const updateQuestionExtraParams = async (
  hash: string,
  questionId: number,
  data: UpdateQuestionExtraParamsRequest
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/extra_params`,
    data,
  });
};

/**
 * Восстановить вопрос (PUT /api/surveys/:hash/question/:questionId/restore).
 */
export const restoreQuestion = async (
  hash: string,
  questionId: number
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PUT",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/restore`,
  });
};

/**
 * Удалить вопрос (DELETE /api/surveys/:hash/question/:questionId).
 */
export const deleteQuestion = async (
  hash: string,
  questionId: number
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "DELETE",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}`,
  });
};

export { QuestionType };
