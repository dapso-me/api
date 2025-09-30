CREATE SCHEMA project;

CREATE TABLE project.projects (
  id          UUID PRIMARY KEY,
  customer_id UUID NOT NULL,
  slug        VARCHAR(255) UNIQUE NOT NULL,
  name        VARCHAR(255) NOT NULL,

  removed_at  TIMESTAMP WITHOUT TIME ZONE DEFAULT NULL,
  created_at  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE project.billing_periods (
  id         UUID PRIMARY KEY,
  name       VARCHAR(255) NOT NULL,
  days_count INTEGER NOT NULL,
  discount   INTEGER NOT NULL
);

CREATE TABLE project.modules (
  id    UUID PRIMARY KEY,
  slug  VARCHAR(255) UNIQUE NOT NULL,
  name  VARCHAR(255) NOT NULL
);

CREATE TABLE project.module_prices (
  module_id         UUID NOT NULL,
  billing_period_id UUID NOT NULL,
  price             INTEGER NOT NULL
);

CREATE TABLE project.subscriptions (
  id                      UUID PRIMARY KEY,
  project_id              UUID NOT NULL,
  billing_period_id       UUID NOT NULL,
  billing_period_days     NOT NULL,
  billing_period_discount UUID NOT NULL,
  total_amount            INTEGER NOT NULL,
  is_active               BOOLEAN NOT NULL,
  expires_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at              TIMESTAMP WITHOUT TIME ZONE NOT NULL,

  FOREIGN KEY (project_id) REFERENCES project.projects (id),
  FOREIGN KEY (billing_period) REFERENCES project.billing_periods (id)
);

CREATE TABLE project.subscription_modules (
  subscription_id UUID NOT NULL,
  module_id       UUID NOT NULL,
  module_price    INTEGER NOT NULL,
  
  UNIQUE(subscription_id, module_id),

  FOREIGN KEY (subscription_id) REFERENCES project.subscriptions (id),
  FOREIGN KEY (module_id) REFERENCES project.modules (id)
);