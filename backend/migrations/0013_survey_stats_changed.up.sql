ALTER TABLE survey_stats
    DROP COLUMN views_count,               -- Удаляем поле views_count
    DROP COLUMN completion_rate,          -- Удаляем поле completion_rate
    ADD COLUMN started_interviews INT4 DEFAULT 0 NOT NULL,   -- Добавляем поле started_interviews
    ADD COLUMN completed_interviews INT4 DEFAULT 0 NOT NULL; -- Добавляем поле completed_interviews