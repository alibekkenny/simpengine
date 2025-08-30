CREATE TABLE event_steps(
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    description TEXT,
    step_order INT,
    event_id BIGINT REFERENCES romantic_events(id) ON DELETE CASCADE
);