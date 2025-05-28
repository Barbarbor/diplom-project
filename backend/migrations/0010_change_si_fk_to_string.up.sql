-- 20250526163500_alter_survey_interviews_id.up.sql

BEGIN;

-- 1) Удаляем внешний ключ survey_answers_interview_id_fkey из survey_answers
ALTER TABLE survey_answers DROP CONSTRAINT survey_answers_interview_id_fkey;

-- 2) Удаляем ограничение первичного ключа на survey_interviews
ALTER TABLE survey_interviews DROP CONSTRAINT survey_interviews_pkey;

-- 3) Изменяем тип столбца id в survey_interviews на VARCHAR(36)
ALTER TABLE survey_interviews ALTER COLUMN id TYPE VARCHAR(36) USING (id::VARCHAR);

-- 4) Изменяем тип столбца interview_id в survey_answers на VARCHAR(36)
ALTER TABLE survey_answers ALTER COLUMN interview_id TYPE VARCHAR(36) USING (interview_id::VARCHAR);

-- 5) Устанавливаем id как первичный ключ в survey_interviews
ALTER TABLE survey_interviews ADD PRIMARY KEY (id);

-- 6) Восстанавливаем внешний ключ в survey_answers
ALTER TABLE survey_answers ADD CONSTRAINT survey_answers_interview_id_fkey
    FOREIGN KEY (interview_id) REFERENCES survey_interviews(id);

COMMIT;