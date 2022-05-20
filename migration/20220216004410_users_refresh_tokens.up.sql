CREATE TABLE IF NOT EXISTS users_refresh_tokens
(
    id            bigserial
        CONSTRAINT users_refresh_tokens_pk
            PRIMARY KEY,
    user_id       bigint      NOT NULL,
    refresh_token varchar(36) NOT NULL,
    expiry_date   timestamp   NOT NULL
);