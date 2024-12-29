CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL,
    chat_id BIGINT NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS
    users_username_uidx ON users (username);

CREATE TABLE IF NOT EXISTS words (
    word TEXT NOT NULL,
    translated_word TEXT NOT NULL,
    example TEXT NOT NULL,
    translated_example TEXT NOT NULL,
    username TEXT NOT NULL,
    correct_answer_count INTEGER NOT NULL DEFAULT 0,
    last_correct_answer TIMESTAMP,
    created TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (word, username)
);