import {SurveyQuestion } from "./question";



export interface SurveyQuestionWithAnswer extends Omit<SurveyQuestion, 'created_at'| 'updated_at' | 'question_original_id' | 'question_state' > {
    answer?: string;
};

export type SurveyWithAnswers = SurveyQuestionWithAnswer[];