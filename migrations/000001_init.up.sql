CREATE TABLE links
(

    link_id varchar(30)  NOT NULL UNIQUE,
    url     varchar(250) NOT NULL,

    PRIMARY KEY (link_id)
);

CREATE INDEX ON links (link_id) INCLUDE (url);
CREATE INDEX ON links USING HASH (url);
