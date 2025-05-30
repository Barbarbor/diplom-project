// src/types/survey.ts
import { SurveyQuestion } from "./question";
// Повторно используемые перечисления
export enum SurveyState {
  Draft = "DRAFT",
  Active = "ACTIVE",
}

export enum SurveyAction {
  Create = "create",
  Update = "update",
  Delete = "delete",
  Pass = "pass",
}

// Типы запроса/ответа

/** Ответ сервера при создании опроса */
export interface CreateSurveyResponse {
  hash: string;
}

/** Краткая информация об опросе в списке */
export interface SurveySummary {
  title: string;
  created_at: string; // ISO
  updated_at: string; // ISO
  hash: string;
  state: SurveyState;
  completed_interviews: number;
}

/** Ответ сервера при запросе списка опросов */
export interface GetSurveysResponse {
  surveys: SurveySummary[];
}

/** Полная информация об опросе */
export interface SurveyDetail {
  title: string;
  creator: string;
  created_at: string;
  updated_at: string;
  state: SurveyState;
}

/** Ответ сервера при запросе одного опроса */
export interface GetSurveyResponse {
  survey: SurveyDetail;
}

/** Тело PATCH-запроса для обновления title */
export interface UpdateSurveyRequest {
  title: string;
}

export interface MutationResponse {
  message: string;
}
export interface SurveyDetail {
  title: string;
  hash: string;
  creator: string;
  created_at: string;
  updated_at: string;
  state: SurveyState;
  questions: SurveyQuestion[]; // Added
}

export type { SurveyQuestion };
