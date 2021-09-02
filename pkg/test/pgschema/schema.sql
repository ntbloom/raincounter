BEGIN;
DROP TABLE IF EXISTS rain;
CREATE TABLE rain
(
    id               SERIAL PRIMARY KEY,
    gw_timestamp     TIMESTAMPTZ NOT NULL,
    server_timestamp TIMESTAMPTZ NOT NULL,
    amount           FLOAT     NOT NULL
);

DROP TABLE IF EXISTS temperature;
CREATE TABLE temperature
(
    id               SERIAL PRIMARY KEY,
    gw_timestamp     TIMESTAMPTZ NOT NULL,
    server_timestamp TIMESTAMPTZ NOT NULL,
    value            INTEGER   NOT NULL
);

DROP TABLE IF EXISTS mappings;
CREATE TABLE mappings
(
    id       INTEGER PRIMARY KEY,
    longname TEXT
);

INSERT INTO mappings (id, longname)
VALUES (2, 'soft reset'),
       (3, 'hard reset'),
       (4, 'pause'),
       (5, 'unpause'),
       (6, NULL),
       (7, NULL)
;

DROP TABLE IF EXISTS event_log;
CREATE TABLE event_log
(
    id        INTEGER PRIMARY KEY,
    tag       INTEGER   NOT NULL,
    value     INTEGER   NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (tag) REFERENCES mappings (id)
);
COMMIT;