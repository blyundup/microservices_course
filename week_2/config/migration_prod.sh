#!/bin/bash
source prod.env

export MIGRATION_DSN="host=pg-prod port=5432 dbname=$PG_DB user=$PG_USER password=$PG_PASSWORD sslmode=disable"

goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v