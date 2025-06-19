ALTER TABLE survey_roles
ADD CONSTRAINT unique_survey_user UNIQUE (survey_id, user_id);