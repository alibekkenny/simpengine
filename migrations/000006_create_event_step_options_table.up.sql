CREATE TABLE event_step_options(
    id SERIAL PRIMARY KEY,
    label VARCHAR, 
    description TEXT, 
    img_id UUID, 
    event_step_id BIGINT REFERENCES event_steps(id) ON DELETE CASCADE
);