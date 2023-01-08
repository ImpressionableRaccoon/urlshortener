CREATE TABLE links
(
    id      varchar(255) NOT NULL UNIQUE,
    url     varchar(255) NOT NULL,
    user_id uuid         NOT NULL
);