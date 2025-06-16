// src/types/question.ts
import { Option } from "./option";
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
  

  
  /** Полная структура вопроса с опциями */
  export interface SurveyQuestion {
    id: number;
    question_original_id: number;
    survey_id: number;
    label: string;
    type: QuestionType;
    question_order: number;
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    extra_params: Record<string, any>;
    options?: Option[];
    question_state: QuestionState;
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
    new_order: number;
  }
  
  /** Тело обновления extra_params */
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  export type UpdateQuestionExtraParamsRequest = Record<string, any>;
  