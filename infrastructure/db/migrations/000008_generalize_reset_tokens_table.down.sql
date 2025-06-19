ALTER TABLE action_tokens DROP COLUMN payload;
ALTER TABLE action_tokens DROP COLUMN type;
ALTER TABLE action_tokens RENAME TO password_reset_tokens;
DROP TYPE token_action_type;