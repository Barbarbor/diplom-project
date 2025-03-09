import serverRequest from "@/lib/serverApi";
import request from "@/lib/api";
import { GetSurveyResponse } from "@/types/survey";

// Функция для создания опроса (POST /api/surveys)
export const createSurvey = async () => {
  const response = await request<{hash:string}>({
    method: "POST",
    prefix: "/api",
    url: "/surveys",
    cache: { disabled: true },
  });
  return response;
};

// Функция для получения опроса по hash (GET /api/surveys/:hash)
export const getSurvey = async (hash: string) => {
  const response = await serverRequest<GetSurveyResponse>({
    method: "GET",
    prefix: "/api",
    url: `/surveys/${hash}`,
    cache: { disabled: true },
  });
  console.log('resp', response);
  return response;
};