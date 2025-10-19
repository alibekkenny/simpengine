ALTER TABLE event_step_options
DROP COLUMN img_id;

ALTER TABLE event_step_options
ADD COLUMN img_id BIGINT REFERENCES media(id);