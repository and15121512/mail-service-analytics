#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE ROLE analytics WITH LOGIN PASSWORD 'secret';
 	CREATE DATABASE analytics OWNER analytics;
	\connect analytics;
	CREATE TABLE event(
		event_id VARCHAR(50) PRIMARY KEY
		, task_id VARCHAR(50)
		, time TIMESTAMP
		, type INTEGER
		, status INTEGER
	);
	INSERT INTO event VALUES (
		'111'
		, '222'
		, NOW()
		, 0
		, 0
	);
	GRANT ALL PRIVILEGES ON TABLE event TO analytics;
EOSQL

echo "END SETUP"