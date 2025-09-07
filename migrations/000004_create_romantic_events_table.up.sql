CREATE TABLE romantic_events(
    id SERIAL PRIMARY KEY,
    event_date TIMESTAMP,
    title VARCHAR,
    status VARCHAR,
    description TEXT,
    public_token TEXT,
    published_at TIMESTAMP,
    simp_target_id BIGINT REFERENCES simp_targets(id),
    user_id BIGINT REFERENCES users(id)
);