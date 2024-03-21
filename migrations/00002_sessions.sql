-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    -- unique constraint enforces one-to-one relationship for user and token
    -- the downside to that is that users cannot log in from multiple devices
    -- 'references users (id)' creates a constraint for foreign key on users' id
    -- 'on delete cascade' deletes sessions automatically when a USER is deleted
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,

    -- the hash of the session token, as we will not store the token itself
    -- this ensures even if the database is leaked, attackers will not be able
    -- to impersonate users unless they manage to reverse the hash
    -- the hash must be unique to ensure no two users have the same token
    token_hash TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd
