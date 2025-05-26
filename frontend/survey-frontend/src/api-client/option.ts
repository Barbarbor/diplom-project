// src/api-client/option.ts

import request, { ApiResponse } from "@/lib/api";
import { OptionResponse, UpdateOptionLabelRequest, UpdateOptionOrderRequest } from "@/types/option";

/** Создать новую опцию */
export const createOption = async (
  hash: string,
  questionId: number
): Promise<ApiResponse<OptionResponse>> => {
  return request<OptionResponse>({
    method: "POST",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/option`
  });
};

/** Изменить лейбл опции */
export const updateOptionLabel = async (
  hash: string,
  questionId: number,
  optionId: number,
  data: UpdateOptionLabelRequest
): Promise<ApiResponse> => {
  return request({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/option/${optionId}`,
    data,
  });
};

/** Изменить порядок опции */
export const updateOptionOrder = async (
  hash: string,
  questionId: number,
  optionId: number,
  data: UpdateOptionOrderRequest
): Promise<ApiResponse> => {
  return request({
    method: "PATCH",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/option/${optionId}/order`,
    data
  });
};

/** Удалить опцию */
export const deleteOption = async (
  hash: string,
  questionId: number,
  optionId: number
): Promise<ApiResponse> => {
  return request({
    method: "DELETE",
    prefix: "/api",
    url: `/surveys/${hash}/question/${questionId}/option/${optionId}`,
  });
};
