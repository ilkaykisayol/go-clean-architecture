CREATE TABLE IF NOT EXISTS users
(
    id              bigserial
    CONSTRAINT users_pk
    PRIMARY KEY,
    username        varchar(30)       NOT NULL,
    email           varchar(100),
    password_hash   varchar(500)      NOT NULL,
    is_active       boolean DEFAULT FALSE NOT NULL,
    is_programmatic boolean DEFAULT FALSE NOT NULL
);