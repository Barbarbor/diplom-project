ALTER TABLE survey_questions_temp
ADD CONSTRAINT fk_survey_questions_temp_survey_id
FOREIGN KEY (survey_id)
REFERENCES surveys (id)
ON DELETE CASCADE;