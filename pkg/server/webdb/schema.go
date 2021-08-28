package webdb

const (
	//nolint
	webDBSchema = `
BEGIN TRANSACTION;
DROP TABLE IF EXISTS rain;
CREATE TABLE rain (
	id INTEGER PRIMARY KEY,
	gw_timestamp TEXT,
	server_timestamp TEXT,
	amount REAL
);

DROP TABLE IF EXISTS temperature;
CREATE TABLE temperature (
	if INTEGER PRIMARY KEY,
	timestamp TEXT NOT NULL,
	value INTEGER NOT NULL
);

DROP TABLE IF EXISTS mappings;
CREATE TABLE mappings (
	id INTEGER PRIMARY KEY,
	longname TEXT
);

INSERT INTO mappings (id, longname)
VALUES
	(2, "soft reset event"),
	(3, "hard reset event"),
	(4, "pause"),
	(5, "unpause"),
	(6, NULL),
	(7, NULL)
;

DROP TABLE IF EXISTS log;
CREATE TABLE log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tag INTEGER NOT NULL,
	value INTEGER NOT NULL,
	timestamp TEXT NOT NULL, --created by go
	FOREIGN KEY (tag) REFERENCES mappings(id)
);

COMMIT;
`
)
