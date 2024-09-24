CREATE TABLE IF NOT EXISTS short_url (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES "user" (user_id) ON DELETE CASCADE,
    original_url TEXT NOT NULL CHECK(original_url <> ''),
    custom_code VARCHAR(32) NOT NULL CHECK(custom_code <> ''),
    short_code CHAR(8) NOT NULL CHECK(short_code <> ''),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
