BEGIN;

-- 1. Создаем таблицу surveys_temp для временного хранения опросов.
-- Убираем поле state, так как оно остается в основной таблице.
CREATE TABLE IF NOT EXISTS surveys_temp (
    survey_original_id INTEGER PRIMARY KEY REFERENCES surveys(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 2. Создаем таблицу survey_questions_temp для временного хранения вопросов.
-- Здесь survey_id ссылается на surveys_temp.
CREATE TABLE IF NOT EXISTS survey_questions_temp (
    id SERIAL PRIMARY KEY,
    question_original_id INTEGER REFERENCES survey_questions(id) ON DELETE CASCADE,
    survey_id INTEGER NOT NULL REFERENCES surveys_temp(survey_original_id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    type TEXT NOT NULL,
    question_order INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 3. Создаем таблицу survey_options_temp для временного хранения опций вопросов.
-- Здесь question_id ссылается на survey_questions_temp.
CREATE TABLE IF NOT EXISTS survey_options_temp (
    id SERIAL PRIMARY KEY,
    option_original_id INTEGER REFERENCES survey_options(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES survey_questions_temp(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 4. Создаем индексы для ускорения поиска по внешним ключам

CREATE INDEX IF NOT EXISTS idx_surveys_temp_original_id 
    ON surveys_temp(survey_original_id);

CREATE INDEX IF NOT EXISTS idx_survey_questions_temp_original_id 
    ON survey_questions_temp(question_original_id);

CREATE INDEX IF NOT EXISTS idx_survey_questions_temp_survey_id 
    ON survey_questions_temp(survey_id);

CREATE INDEX IF NOT EXISTS idx_survey_options_temp_question_id 
    ON survey_options_temp(question_id);

CREATE INDEX IF NOT EXISTS idx_survey_options_temp_original_id 
    ON survey_options_temp(option_original_id);

COMMIT;