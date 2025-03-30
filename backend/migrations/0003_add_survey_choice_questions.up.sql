ALTER TABLE survey_questions
  ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW(),
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW(),
  ADD COLUMN IF NOT EXISTS question_order INTEGER NOT NULL DEFAULT 1;

-- Удаляем столбец options, если он существует
ALTER TABLE survey_questions
  DROP COLUMN IF EXISTS options;

-- Создаем новую таблицу survey_options для опций вопросов
CREATE TABLE IF NOT EXISTS survey_options (
  id SERIAL PRIMARY KEY,
  question_id INTEGER NOT NULL REFERENCES survey_questions(id) ON DELETE CASCADE,
  label TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

-- Создаем индекс для ускорения поиска опций по question_id
CREATE INDEX IF NOT EXISTS idx_survey_options_question_id
  ON survey_options(question_id);