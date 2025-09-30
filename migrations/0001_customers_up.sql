CREATE SCHEMA customers;

CREATE TABLE customers.customers (
  id         UUID PRIMARY KEY,
  login      VARCHAR(254) NOT NULL UNIQUE,
  password   VARCHAR(255) NOT NULL,
  name       VARCHAR(255) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers.sessions (
  access_token VARCHAR(512) PRIMARY KEY,
  customer_id  UUID NOT NULL,
  ip           VARCHAR(64) NOT NULL,
  user_agent   VARCHAR(4096) NOT NULL,
  created_at   TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

  FOREIGN KEY (customer_id) REFERENCES customers.customers (id) ON DELETE CASCADE
);