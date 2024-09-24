CREATE TABLE IF NOT EXISTS "user" (
    user_id UUID PRIMARY KEY,
    username VARCHAR(32) NOT NULL CHECK(username <> ''),
    email VARCHAR(64) UNIQUE NOT NULL CHECK(email <> ''),
    hashed_password CHAR(60) NOT NULL CHECK(octet_length(hashed_password) <> 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
