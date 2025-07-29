CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR,
    email VARCHAR,
    password_hash VARCHAR,
    role VARCHAR,
    created_at TIMESTAMP
);