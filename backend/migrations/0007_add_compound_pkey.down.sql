BEGIN;

-- Восстанавливаем первичный ключ для survey_questions на поле id
ALTER TABLE survey_questions DROP CONSTRAINT IF EXISTS survey_questions_pkey;
ALTER TABLE survey_questions ADD CONSTRAINT survey_questions_pkey PRIMARY KEY (id);

DROP INDEX IF EXISTS idx_survey_questions_survey_id;

-- Восстанавливаем первичный ключ для survey_questions_temp на поле id
ALTER TABLE survey_questions_temp DROP CONSTRAINT IF EXISTS survey_questions_temp_pkey;
ALTER TABLE survey_questions_temp ADD CONSTRAINT survey_questions_temp_pkey PRIMARY KEY (id);

DROP INDEX IF EXISTS idx_survey_questions_temp_survey_id;

-- Восстанавливаем первичный ключ для survey_options на поле id
ALTER TABLE survey_options DROP CONSTRAINT IF EXISTS survey_options_pkey;
ALTER TABLE survey_options ADD CONSTRAINT survey_options_pkey PRIMARY KEY (id);

DROP INDEX IF EXISTS idx_survey_options_question_id;

-- Восстанавливаем первичный ключ для survey_options_temp на поле id
ALTER TABLE survey_options_temp DROP CONSTRAINT IF EXISTS survey_options_temp_pkey;
ALTER TABLE survey_options_temp ADD CONSTRAINT survey_options_temp_pkey PRIMARY KEY (id);

DROP INDEX IF EXISTS idx_survey_options_temp_question_id;

COMMIT;
