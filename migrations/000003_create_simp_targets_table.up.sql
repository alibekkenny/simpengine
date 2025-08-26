CREATE TABLE simp_targets (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    description TEXT NOT NULL,
    user_id BIGINT REFERENCES users(id)
);