CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    -- unique constraint enforces one-to-one relationship for user and token
    -- the downside to that is that users cannot log in from multiple devices
    user_id INT UNIQUE,

    -- the hash of the session token, as we will not store the token itself
    -- this ensures even if the database is leaked, attackers will not be able
    -- to impersonate users unless they manage to reverse the hash
    -- the hash must be unique to ensure no two users have the same token
    token_hash TEXT UNIQUE NOT NULL
);
