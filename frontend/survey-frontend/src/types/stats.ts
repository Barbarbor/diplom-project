/* eslint-disable @typescript-eslint/no-explicit-any */
import { QuestionType } from "./question";

interface OptionStats {
  id: number;
  label: string;
}

export interface QuestionStats {
  id: number;
  label: string;
  type: QuestionType;
  options?: OptionStats[];
  answers: string[];
  extra_params?: Record<string, any>;
}

export interface SurveyStats {
  started_interviews: number;
  completed_interviews: number;
  average_completion_time: number;
  questions: QuestionStats[];
}