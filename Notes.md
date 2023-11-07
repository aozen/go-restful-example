# Notes

## Postgres

(WIP)
psql postgres

CREATE ROLE "username" WITH LOGIN PASSWORD 'userpassword';
ALTER ROLE "username" CREATEDB;

psql postgres -U username

create database pq_test;
grant all privileges on database pq_test to username;

\c pq_test

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
(WIP)