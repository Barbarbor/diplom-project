-- 0003_add_survey_tables_with_enums.sql

-- Создание ENUM для действий с опросами
DO $$ BEGIN
    CREATE TYPE action_enum AS ENUM ('create', 'update', 'delete', 'pass');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Создание ENUM для статусов прохождения опроса
DO $$ BEGIN
    CREATE TYPE status_enum AS ENUM ('in_progress', 'completed');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Создание ENUM для типов вопросов
DO $$ BEGIN
    CREATE TYPE question_type_enum AS ENUM ('single_choice', 'multi_choice');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

-- Таблица "Пользователи"
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Таблица "Профили пользователей"
CREATE TABLE IF NOT EXISTS user_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id),
    first_name TEXT,
    last_name TEXT,
    birth_date DATE,
    phone_number TEXT,
    lang TEXT DEFAULT 'en'
);

-- Таблица "Роли пользователей"
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    roles TEXT[] NOT NULL
);

-- Таблица "Опросы"
CREATE TABLE IF NOT EXISTS surveys (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    author_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    hash TEXT UNIQUE NOT NULL
);

-- Таблица "Вопросы опросов"
CREATE TABLE IF NOT EXISTS survey_questions (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES surveys(id),
    label TEXT NOT NULL,
    type question_type_enum NOT NULL,
    options JSON NOT NULL
);

-- Таблица "Прохождения опросов"
CREATE TABLE IF NOT EXISTS survey_interviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    survey_id INTEGER REFERENCES surveys(id),
    status status_enum NOT NULL DEFAULT 'in_progress',
    start_time TIMESTAMP DEFAULT NOW(),
    end_time TIMESTAMP
);

-- Таблица "Ответы на вопросы"
CREATE TABLE IF NOT EXISTS survey_answers (
    id SERIAL PRIMARY KEY,
    interview_id INTEGER NOT NULL REFERENCES survey_interviews(id),
    question_id INTEGER NOT NULL,
    options JSON NOT NULL
);

-- Таблица "Статистика опросов"
CREATE TABLE IF NOT EXISTS survey_stats (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES surveys(id),
    views_count INTEGER DEFAULT 0,
    completion_rate FLOAT DEFAULT 0
);

-- Таблица "Роли в опросах"
CREATE TABLE IF NOT EXISTS survey_roles (
    id SERIAL PRIMARY KEY,
    survey_id INTEGER NOT NULL REFERENCES surveys(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    roles TEXT[] NOT NULL
);

-- Таблица "Действия с опросами"
CREATE TABLE IF NOT EXISTS survey_actions (
    id SERIAL PRIMARY KEY,
    action action_enum NOT NULL,
    survey_id INTEGER REFERENCES surveys(id),
    user_id INTEGER REFERENCES users(id),
    body JSON,
    action_time TIMESTAMP DEFAULT NOW()
);

-- Создание индексов для ускорения запросов
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_surveys_hash ON surveys(hash);
CREATE INDEX IF NOT EXISTS idx_survey_questions_survey_id ON survey_questions(survey_id);
CREATE INDEX IF NOT EXISTS idx_survey_interviews_status ON survey_interviews(status);
CREATE INDEX IF NOT EXISTS idx_survey_answers_question_id ON survey_answers(question_id);
CREATE INDEX IF NOT EXISTS idx_survey_roles_user_id ON survey_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_survey_actions_action ON survey_actions(action);
