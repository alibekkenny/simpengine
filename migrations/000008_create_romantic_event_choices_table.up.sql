CREATE TABLE romantic_event_choices(
    id SERIAL PRIMARY KEY,
    event_id BIGINT REFERENCES romantic_events(id) NOT NULL,
    step_id BIGINT REFERENCES event_steps(id) NOT NULL,
    options_ids INT[] NOT NULL
);