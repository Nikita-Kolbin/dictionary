UPDATE words
SET last_correct_answer = NOW()
WHERE last_correct_answer IS NULL;

ALTER TABLE IF EXISTS words
ALTER COLUMN last_correct_answer SET DEFAULT NOW();

ALTER TABLE IF EXISTS words
ALTER COLUMN last_correct_answer SET NOT NULL;