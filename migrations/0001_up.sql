CREATE SCHEMA customers;

CREATE TABLE customers.customers (
  id         UUID PRIMARY KEY,
  email      TEXT(254) NOT NULL UNIQUE,
  password   TEXT(255) NOT NULL,
  name       TEXT(255) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP()
);

CREATE TABLE customers.sessions (
  access_token TEXT(512) PRIMARY KEY,
  customer_id  UUID NOT NULL,
  ip           TEXT(64) NOT NULL,
  user_agent   TEXT(4096) NOT NULL,
  created_at   TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP()
);