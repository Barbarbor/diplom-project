

-- 1. Создаем ENUM question_state, если он не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'question_state') THEN
        CREATE TYPE question_state AS ENUM ('ACTUAL', 'NEW', 'CHANGED', 'DELETED');
    END IF;
END$$;

-- 2. Создаем ENUM option_state, если он не существует
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'option_state') THEN
        CREATE TYPE option_state AS ENUM ('ACTUAL', 'NEW', 'CHANGED', 'DELETED');
    END IF;
END$$;

-- 3. Добавляем новое поле question_state в таблицу survey_questions_temp
ALTER TABLE survey_questions_temp
  ADD COLUMN IF NOT EXISTS question_state question_state NOT NULL DEFAULT 'ACTUAL';

-- 4. Добавляем новое поле option_state в таблицу survey_options_temp
ALTER TABLE survey_options_temp
  ADD COLUMN IF NOT EXISTS option_state option_state NOT NULL DEFAULT 'ACTUAL';
