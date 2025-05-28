BEGIN;

-- Add the is_demo column with a default value of FALSE
ALTER TABLE survey_interviews
ADD COLUMN is_demo BOOLEAN NOT NULL DEFAULT FALSE;

COMMIT;