CREATE SCHEMA customers;

CREATE TABLE customers.customers (
  id         UUID PRIMARY KEY,
  email      TEXT(254) NOT NULL UNIQUE,
  password   TEXT(255) NOT NULL,
  name       TEXT(255) NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP()
);

CREATE TABLE customers.sessions (
  access_token VARCHAR(512) PRIMARY KEY,
  customer_id  UUID NOT NULL,
  ip           VARCHAR(64) NOT NULL,
  user_agent   VARCHAR(4096) NOT NULL,
  created_at   TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP(),

  FOREIGN KEY (customer_id) REFERENCES customers.customers (id) ON DELETE CASCADE
);

CREATE TABLE customers.otp (
  id         UUID PRIMARY KEY,
  email      TEXT(254) NOT NULL,
  code       VARCHAR(6) NOT NULL,
  purpose    VARCHAR(255) NOT NULL,
  ip         VARCHAR(64) NOT NULL,
  expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP() NOT NULL
);