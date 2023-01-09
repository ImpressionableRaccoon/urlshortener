CREATE TABLE links
(
    id      varchar(255) NOT NULL UNIQUE,
    url     varchar(255) NOT NULL UNIQUE,
    user_id uuid         NOT NULL
);