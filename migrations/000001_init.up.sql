CREATE TABLE IF NOT EXISTS users (
    username TEXT NOT NULL,
    chat_id BIGINT NOT NULL,
    notification_word_count SMALLINT NOT NULL DEFAULT 10,
    created TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS
    users_username_uidx ON users (username);

CREATE TABLE IF NOT EXISTS words (
    id BIGSERIAL PRIMARY KEY,
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

CREATE TABLE IF NOT EXISTS notification_times (
    time TIME NOT NULL,
    username TEXT NOT NULL,
    UNIQUE (time, username)
);
CREATE INDEX IF NOT EXISTS
    notification_times_idx ON notification_times (time);
