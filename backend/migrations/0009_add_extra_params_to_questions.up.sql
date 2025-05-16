BEGIN;

-- 1) Добавляем в survey_questions
ALTER TABLE survey_questions
  ADD COLUMN IF NOT EXISTS extra_params JSONB NOT NULL DEFAULT '{}'::jsonb;

-- 2) Добавляем в survey_questions_temp
ALTER TABLE survey_questions_temp
  ADD COLUMN IF NOT EXISTS extra_params JSONB NOT NULL DEFAULT '{}'::jsonb;

COMMIT;
