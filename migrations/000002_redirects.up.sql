CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE redirects
(
    redirect_id serial8                     NOT NULL,
    session_id  varchar(36)                 NOT NULL,
    link_id     varchar(30)                 NOT NULL REFERENCES links,
    time        timestamp without time zone NOT NULL
);
