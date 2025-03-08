export enum SurveyState {
  Draft = "DRAFT",
  Active = "ACTIVE",
}

export enum QuestionType {
  SingleChoice = "single_choice",
  MultiChoice = "multi_choice",
}

export enum SurveyStatus {
  InProgress = "in_progress",
  Completed = "completed",
}

export enum SurveyAction {
  Create = "create",
  Update = "update",
  Delete = "delete",
  Pass = "pass",
}

// Основной тип опроса
export interface Survey {
  id: number;
  title: string;
  author_id: number;
  created_at: string; // ISO строка даты, например, "2023-03-10T12:34:56Z"
  updated_at: string;
  hash: string;
  state: SurveyState;
}

// Тип для опций вопроса
export interface Option {
  id: number;
  label: string;
}

// Тип для вопросов опроса
export interface SurveyQuestion {
  id: number;
  survey_id: number;
  label: string;
  type: QuestionType;
  options: Option[];
}

// Тип для прохождений опросов (интервью)
export interface SurveyInterview {
  id: number;
  user_id: number;
  survey_id: number;
  status: SurveyStatus;
  start_time: string; // ISO строка даты
  end_time?: string; // Может быть undefined, если опрос ещё не завершён
}

// Тип для ответов на вопросы опроса
export interface SurveyAnswer {
  id: number;
  interview_id: number;
  question_id: number;
  options: number[]; // Массив ID выбранных опций
}

// Тип для статистики опроса
export interface SurveyStat {
  id: number;
  survey_id: number;
  views_count: number;
  completion_rate: number; // например, процент завершивших опрос
}

// Тип для ролей в опросах
export interface SurveyRole {
  id: number;
  survey_id: number;
  user_id: number;
  roles: string[]; // Массив ролей
}

// Тип для логов действий с опросами
export interface SurveyActionLog {
  id: number;
  action: SurveyAction;
  survey_id?: number; // Может быть null, если действие не связано напрямую с опросом
  user_id?: number;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  body?: any; // Дополнительная информация в формате JSON
  action_time: string; // ISO строка даты
}

export interface GetSurveyResponse {
  survey: {
    title: string;
    creator: string;
    created_at: Date;
    updated_at: Date;
    state: SurveyState;
  };
}
