CREATE TABLE users
(
    id              UUID PRIMARY KEY,
    username        VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    roles           TEXT[]       NOT NULL,
    active          BOOLEAN      NOT NULL DEFAULT TRUE,
    failed_attempts INT          NOT NULL DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    email_verified  bool         NOT NULL DEFAULT FALSE,

    CONSTRAINT chk_valid_roles CHECK (roles <@ ARRAY ['ADMIN', 'MANAGER', 'CASHIER', 'STOCK_CLERK']::TEXT[])
);

CREATE INDEX idx_users_roles ON users USING GIN (roles);
CREATE INDEX idx_users_username ON users (username);