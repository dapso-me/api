CREATE SCHEMA plans;

CREATE TABLE plans.plans (
  id               UUID PRIMARY KEY,
  name             VARCHAR(255) NOT NULL,
  period_days      INTEGER NOT NULL,
  price_per_period INTEGER NOT NULL
);

CREATE TABLE plans.modules (
  id    UUID PRIMARY KEY,
  slug  VARCHAR(20) UNIQUE NOT NULL,
  name  VARCHAR(255) NOT NULL,
  price INTEGER NOT NULL
);

CREATE TABLE plans.plan_modules (
  plan_id    UUID NOT NULL,
  module_id  UUID NOT NULL,
  is_default BOOLEAN NOT NULL,

  UNIQUE(plan_id, module_id),
  FOREIGN KEY (plan_id) REFERENCES plans.plans (id) ON DELETE CASCADE,
  FOREIGN KEY (module_id) REFERENCES plans.modules (id) ON DELETE CASCADE
);

CREATE TABLE plans.subscriptions (
  id           UUID PRIMARY KEY,
  customer_id  UUID UNIQUE NOT NULL,
  plan_id      UUID NOT NULL,
  total_amount INTEGER NOT NULL,
  is_active    BOOLEAN NOT NULL,
  expires_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL,
  created_at   TIMESTAMP WITHOUT TIME ZONE NOT NULL,

  FOREIGN KEY (customer_id) REFERENCES customers.customers (id) ON DELETE CASCADE,
  FOREIGN KEY (plan_id) REFERENCES plans.plans (id) ON DELETE CASCADE
);

CREATE TABLE plans.subscription_modules (
  subscription_id UUID NOT NULL,
  module_id UUID  NOT NULL,
  is_addon        BOOLEAN NOT NULL,
  addon_price     INTEGER NOT NULL,

  UNIQUE (subscription_id, module_id),
  FOREIGN KEY (subscription_id) REFERENCES plans.subscriptions (id) ON DELETE CASCADE,
  FOREIGN KEY (module_id) REFERENCES plans.modules (id) ON DELETE CASCADE,
);