BEGIN;
DELETE FROM rain;
DELETE FROM temperature;
DELETE FROM status_log;
DELETE FROM event_log;
COMMIT;