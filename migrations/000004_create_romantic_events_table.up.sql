CREATE TABLE romantic_events(
    id SERIAL PRIMARY KEY,
    event_date TIMESTAMP,
    title VARCHAR,
    description TEXT,
    simp_target_id BIGINT REFERENCES simp_targets(id),
    user_id BIGINT REFERENCES users(id)
);