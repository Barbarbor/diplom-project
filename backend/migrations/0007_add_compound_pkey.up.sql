BEGIN;

-- Для таблицы survey_questions
ALTER TABLE survey_questions DROP CONSTRAINT IF EXISTS survey_questions_pkey CASCADE;
ALTER TABLE survey_questions ADD CONSTRAINT survey_questions_pkey PRIMARY KEY (survey_id, id);
CREATE INDEX IF NOT EXISTS idx_survey_questions_survey_id ON survey_questions(survey_id);

-- Для таблицы survey_questions_temp
ALTER TABLE survey_questions_temp DROP CONSTRAINT IF EXISTS survey_questions_temp_pkey CASCADE;
ALTER TABLE survey_questions_temp ADD CONSTRAINT survey_questions_temp_pkey PRIMARY KEY (survey_id, id);
CREATE INDEX IF NOT EXISTS idx_survey_questions_temp_survey_id ON survey_questions_temp(survey_id);

-- Для таблицы survey_options
ALTER TABLE survey_options DROP CONSTRAINT IF EXISTS survey_options_pkey CASCADE;
ALTER TABLE survey_options ADD CONSTRAINT survey_options_pkey PRIMARY KEY (question_id, id);
CREATE INDEX IF NOT EXISTS idx_survey_options_question_id ON survey_options(question_id);

-- Для таблицы survey_options_temp
ALTER TABLE survey_options_temp DROP CONSTRAINT IF EXISTS survey_options_temp_pkey CASCADE;
ALTER TABLE survey_options_temp ADD CONSTRAINT survey_options_temp_pkey PRIMARY KEY (question_id, id);
CREATE INDEX IF NOT EXISTS idx_survey_options_temp_question_id ON survey_options_temp(question_id);

COMMIT;