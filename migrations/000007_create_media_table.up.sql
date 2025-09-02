CREATE TABLE media(
    id SERIAL PRIMARY KEY,
    object_name VARCHAR NOT NULL,
    original_name VARCHAR NOT NULL,
    mime_type VARCHAR,
    user_id BIGINT REFERENCES users(id),
    size BIGINT,
    created_at TIMESTAMP DEFAULT NOW()
);