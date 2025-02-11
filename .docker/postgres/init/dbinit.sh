#!/bin/bash
set -e
export PGPASSWORD=postgres;
psql -v ON_ERROR_STOP=1 --username "postgres" <<-EOSQL
  CREATE DATABASE recipebot;
  GRANT ALL PRIVILEGES ON DATABASE recipebot TO "postgres";
EOSQL