DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'survey_state') THEN
        CREATE TYPE survey_state AS ENUM ('DRAFT', 'ACTIVE');
    END IF;
END$$;

-- 2. Добавляем колонку state в таблицу surveys
ALTER TABLE surveys
ADD COLUMN state survey_state NOT NULL DEFAULT 'DRAFT';

-- 3. Обновляем существующие записи, если нужно (в данном случае они останутся DRAFT)
UPDATE surveys SET state = 'DRAFT' WHERE state IS NULL;