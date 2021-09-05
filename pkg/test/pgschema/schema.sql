BEGIN;
DROP TABLE IF EXISTS rain CASCADE;
CREATE TABLE rain
(
    id               SERIAL PRIMARY KEY,
    gw_timestamp     TIMESTAMPTZ NOT NULL,
    server_timestamp TIMESTAMPTZ NOT NULL,
    amount           FLOAT       NOT NULL
);

DROP TABLE IF EXISTS temperature CASCADE;
CREATE TABLE temperature
(
    id               SERIAL PRIMARY KEY,
    gw_timestamp     TIMESTAMPTZ NOT NULL,
    server_timestamp TIMESTAMPTZ NOT NULL,
    value            INTEGER     NOT NULL
);

DROP TABLE IF EXISTS mappings CASCADE;
CREATE TABLE mappings
(
    id       INTEGER PRIMARY KEY,
    longname TEXT
);

DROP TABLE IF EXISTS status_codes CASCADE;
CREATE TABLE status_codes
(
    id    INTEGER PRIMARY KEY,
    asset TEXT
);

INSERT INTO status_codes (id, asset)
VALUES (1, 'sensor'),
       (2, 'gateway')
;

DROP TABLE IF EXISTS status_log CASCADE;
CREATE TABLE status_log
(
    id               SERIAL PRIMARY KEY,
    gw_timestamp     TIMESTAMPTZ NOT NULL,
    server_timestamp TIMESTAMPTZ NOT NULL,
    asset            INTEGER     NOT NULL,
    FOREIGN KEY (asset) REFERENCES status_codes (id)

);

INSERT INTO mappings (id, longname)
VALUES (2, 'soft reset'),
       (3, 'hard reset'),
       (4, 'pause'),
       (5, 'unpause'),
       (6, NULL),
       (7, NULL)
;

DROP TABLE IF EXISTS event_log CASCADE;
CREATE TABLE event_log
(
    id        INTEGER PRIMARY KEY,
    tag       INTEGER     NOT NULL,
    value     INTEGER     NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (tag) REFERENCES mappings (id)
);
COMMIT;