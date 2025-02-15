-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id             BIGINT(20) UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,
    xid            VARCHAR(20)         NOT NULL,
    created_at     DATETIME(3)         NOT NULL,
    updated_at     DATETIME(3)         NOT NULL,
    deleted_at     DATETIME(3)         NULL,
    email          VARCHAR(191)        NOT NULL,
    password       LONGTEXT            NULL,
    recovery_token VARCHAR(256)        NULL,
    CONSTRAINT uni_users_email UNIQUE (email),
    CONSTRAINT uni_users_xid UNIQUE (xid),
    INDEX idx_users_xid (xid),
    INDEX idx_users_deleted_at (deleted_at),
    INDEX idx_users_email (email),
    INDEX idx_users_recovery_token (recovery_token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
