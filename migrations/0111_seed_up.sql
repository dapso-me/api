-- customer
-- password is "test"
INSERT INTO customers.customers (id, username, password, name, role, created_at)
VALUES ('f3bbd575-4135-4ac7-9fe4-1bafeedceee0', 'test', '$2a$10$7egYXtMzvwjARpcG8HsUtO4ae.7cYv3WrVTu.xsY1YDRZS3L4quku', 'Test Name', 'customer', NOW());  

INSERT INTO customers.customers (id, username, password, name, role, created_at)
VALUES ('ec931089-2dfb-4f63-974e-7cff75809d43', 'owner', '$2a$10$7egYXtMzvwjARpcG8HsUtO4ae.7cYv3WrVTu.xsY1YDRZS3L4quku', 'Owner', 'owner', NOW());  

-- project
INSERT INTO project.projects (id, customer_id, slug, name, created_at)
VALUES (
  'e4cbd11c-e7c8-4ee0-9f8a-a3281711922f',
  'f3bbd575-4135-4ac7-9fe4-1bafeedceee0',
  'burger-king',
  'Burger King',
  NOW()
);