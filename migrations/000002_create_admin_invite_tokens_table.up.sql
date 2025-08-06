CREATE TABLE admin_invite_tokens (
    id SERIAL PRIMARY KEY,
    token TEXT,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NOT NULL,
    created_by INT REFERENCES users(id),
    used_by INT REFERENCES users(id)
);