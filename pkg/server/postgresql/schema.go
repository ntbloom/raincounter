package postgresql

// schema for the postgresql postgresql
const (
	sqlSchema = `
	BEGIN TRANSACTION;

/* rain gauges, will probably only be one but let's plan for more just in case */
DROP TABLE IF EXISTS gauge CASCADE;
CREATE TABLE gauge (
	id SERIAL PRIMARY KEY,
	short_name TEXT NOT NULL,
	amount_mm REAL,
	amount_in REAL,
	about TEXT NOT NULL
);

/* populate the data for our rain gauge */
INSERT INTO gauge VALUES (
	"The Blue House", 0.2794, 0.11, "the rain gauge that started it all"	
);

/* list of all events */
DROP TABLE IF EXISTS events CASCADE;
CREATE TABLE events (
	id SERIAL PRIMARY KEY,
	short_name TEXT NOT NULL,
	description TEXT NOT NULL
);

/* various sensor events */
DROP TABLE IF EXISTS sensor_events;
CREATE TABLE sensor_events (
	id SERIAL PRIMARY KEY,
	short_name TEXT NOT NULL,
);

/* populate the sensor events */
INSERT INTO sensor_events (short_name) VALUES (
	"softReset",
	"hardReset",
	"pause",
	"unpause"
);

/* LOG THE ACTUAL EVENTS */

/* log temperature events */
DROP TABLE IF EXISTS temperature_log CASCADE;
CREATE TABLE temperature_log (
	id SERIAL PRIMARY KEY,
	db_timestamp TIMESTAMP WITH TIME ZONE DEFAULT now(),
	gauge_id INTEGER,
	gauge_timestamp TIMESTAMP WITH TIME ZONE,
	tempC INTEGER,
	FOREIGN KEY(gauge_id) REFERENCES gauge(id)
);

/* log rain events */
DROP TABLE IF EXISTS rain_log CASCADE;
CREATE TABLE rain_log (
	id SERIAL PRIMARY KEY,
	db_timestamp TIMESTAMP WITH TIME ZONE DEFAULT now(),
	gauge_id INTEGER,
	gauge_timestamp TIMESTAMP WITH TIME ZONE,
	FOREIGN KEY(gauge_id) REFERENCES gauge(id)
);

/* log sensor events */
DROP TABLE IF EXISTS sensor_event_log;
CREATE TABLE sensor_event_log (
	id SERIAL PRIMARY KEY,
	db_timestamp TIMESTAMP WITH TIME ZONE DEFAULT now(),
	gauge_id INTEGER,
	gauge_timestamp TIMESTAMP WITH TIME ZONE,
	sensor_event_id INTEGER,
	FOREIGN KEY(sensor_event_id) REFERENCES sensor_events(id)
);

COMMIT;
	
	
	
	`
)
