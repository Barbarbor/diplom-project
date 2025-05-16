// src/types/question.ts

/** Тип вопроса */
export enum QuestionType {
    SingleChoice = "single_choice",
    MultiChoice  = "multi_choice",
    Consent      = "consent",
    Email        = "email",
    Rating       = "rating",
    Date         = "date",
    ShortText    = "short_text",
    LongText     = "long_text",
    Number       = "number",
  }
  
  /** Временные и постоянные состояния вопроса */
  export enum QuestionState {
    Actual  = "ACTUAL",
    New     = "NEW",
    Changed = "CHANGED",
    Deleted = "DELETED",
  }
  
  /** Опция вопроса (отображается при чтении) */
  export interface Option {
    id: number;
    label: string;
    order: number;
  }
  
  /** Полная структура вопроса с опциями */
  export interface SurveyQuestion {
    id: number;
    survey_id: number;
    label: string;
    type: QuestionType;
    order: number;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    extra_params: Record<string, any>;
    options: Option[];
    created_at: string; // ISO
    updated_at: string; // ISO
  }
  
  /** Ответ сервера на создание/обновление одного вопроса */
  export interface QuestionResponse {
    question: SurveyQuestion;
  }
  
  /** Тело создания вопроса */
  export interface CreateQuestionRequest {
    type: QuestionType;
  }
  
  /** Тело обновления только метки */
  export interface UpdateQuestionLabelRequest {
    label: string;
  }
  
  /** Тело обновления порядка */
  export interface UpdateQuestionOrderRequest {
    newOrder: number;
  }
  
  /** Тело обновления extra_params */
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export type UpdateQuestionExtraParamsRequest = Record<string, any>;
  