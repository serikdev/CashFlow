#!/bin/bash

set -e

# Load environment variables from .env file
if [ -f .env ]; then
  export $(grep -v '^#' .env | xargs -d '\n')
else
  echo ".env file not found"
  exit 1
fi

# Check PostgreSQL connection
echo "Checking PostgreSQL connection..."
pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > /dev/null

if [ $? -ne 0 ]; then
  echo "Failed to connect to the database"
  exit 1
fi

# Read the command argument
CMD=$1

if [ -z "$CMD" ]; then
  echo "Please provide a command: up | down | status | redo | create <name>"
  exit 1
fi

# Migration config
MIGRATIONS_DIR="./migrations"
DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

echo "Executing goose command: $CMD"

# Execute the corresponding goose command
case $CMD in
  up)
    goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" up
    ;;
  down)
    goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" down
    ;;
  status)
    goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" status
    ;;
  redo)
    goose -dir "$MIGRATIONS_DIR" postgres "$DB_URL" redo
    ;;
  create)
    NAME=$2
    if [ -z "$NAME" ]; then
      echo "Please provide a migration name: ./migrate.sh create add_users_table"
      exit 1
    fi
    goose -dir "$MIGRATIONS_DIR" create "$NAME" sql
    ;;
  *)
    echo "Unknown command: $CMD"
    exit 1
    ;;
esac

echo "Goose command '$CMD' executed successfully!"
