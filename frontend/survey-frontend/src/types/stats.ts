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
}

export interface SurveyStats {
  started_interviews: number;
  completed_interviews: number;
  questions: QuestionStats[];
}