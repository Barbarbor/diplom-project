// src/types/option.ts

/** Временные и реальные состояния опции */
export enum OptionState {
  Actual  = "ACTUAL",
  New     = "NEW",
  Changed = "CHANGED",
  Deleted = "DELETED",
}

/** Базовая структура опции, возвращаемая API */
export interface Option {
  id: number;
  question_id: number;
  label: string;
  option_order: number;
  option_state: OptionState;
  created_at: string; // ISO
  updated_at: string; // ISO
}

/** Ответ API при работе с одной опцией */
export interface OptionResponse {
  option: Option;
}

/** Пакет всех опций для вопроса */
export interface OptionsListResponse {
  options: Option[];
}

/** Тело POST /option */
// eslint-disable-next-line @typescript-eslint/no-empty-object-type
export interface CreateOptionRequest {
  // no body: type inferred from question
}

/** Тело PATCH /order */
export interface UpdateOptionOrderRequest {
  new_order: number;
}

/** Тело PATCH / (label) */
export interface UpdateOptionLabelRequest {
  label: string;
}
