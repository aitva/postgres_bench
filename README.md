# Postgres Bench

A repository to benchmark Postgres driver for Golang.

## Setting up Postgres

1. Connect to Postgres as an admin:  `sudo -l postgres psql`
2. Create an user and a database for the application:
  ```sql
  CREATE DATABASE postgres_bench;
  CREATE USER postgres_bench PASSWORD 'postgres_bench';
  ALTER DATABASE postgres_bench OWNER TO postgres_bench;
  ```
3. fill the database: `go run .`

## Running the benchmark

Use the Go benchmark program: `go test -bench .`
