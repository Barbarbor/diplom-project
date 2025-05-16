import request, { ApiResponse } from "@/lib/api";
import {
  CreateSurveyResponse,
  GetSurveysResponse,
  GetSurveyResponse,
  UpdateSurveyRequest,
} from "@/types/survey";

/**
 * Создать новый опрос. Клиентский вызов.
 */
export const createSurvey = async (): Promise<ApiResponse<CreateSurveyResponse>> => {
  return await request<CreateSurveyResponse>({
    method: "POST",
    prefix: "/api",
    url: "/surveys",
  });
};

/**
 * Получить список опросов, созданных пользователем. Клиентский вызов.
 */
export const getSurveys = async (): Promise<ApiResponse<GetSurveysResponse>> => {
  return await request<GetSurveysResponse>({
    method: "GET",
    prefix: "/api",
    url: "/surveys",
  });
};

/**
 * Получить детали одного опроса по hash. Клиентский вызов.
 */
export const getSurvey = async (hash: string): Promise<ApiResponse<GetSurveyResponse>> => {
  return await request<GetSurveyResponse>({
    method: "GET",
    prefix: "/api",
    url: `/surveys/${hash}`,
  });
};

/**
 * Обновить заголовок опроса. Клиентский вызов.
 */
export const updateSurvey = async (
  hash: string,
  data: UpdateSurveyRequest
): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}`,
    data,
  });
};

/**
 * Опубликовать опрос. Клиентский вызов.
 */
export const publishSurvey = async (hash: string): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "POST",
    prefix: "/api",
    url: `/surveys/${hash}/publish`,
  });
};

/**
 * Восстановить опрос из временной области. Клиентский вызов.
 */
export const restoreSurvey = async (hash: string): Promise<ApiResponse<void>> => {
  return await request<void>({
    method: "PUT",
    prefix: "/api",
    url: `/surveys/${hash}/restore`,
  });
};
export type { UpdateSurveyRequest };

