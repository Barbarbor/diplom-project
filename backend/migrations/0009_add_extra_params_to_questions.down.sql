BEGIN;

ALTER TABLE survey_questions DROP COLUMN IF EXISTS extra_params;
ALTER TABLE survey_questions_temp DROP COLUMN IF EXISTS extra_params;

COMMIT;
