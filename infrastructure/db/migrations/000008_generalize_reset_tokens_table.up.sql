CREATE TYPE token_action_type AS ENUM (
    'PASSWORD_RESET',
    'EMAIL_CONFIRMATION'
    );

ALTER TABLE password_reset_tokens
    RENAME TO action_tokens;


ALTER TABLE action_tokens
    ADD COLUMN type token_action_type NOT NULL;

ALTER TABLE action_tokens
    ADD COLUMN payload TEXT;