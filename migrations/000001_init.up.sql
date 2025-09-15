CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password_hash TEXT NOT NULL,
                       role VARCHAR(50) NOT NULL,
                       date_create TIMESTAMPTZ DEFAULT NOW()
);
