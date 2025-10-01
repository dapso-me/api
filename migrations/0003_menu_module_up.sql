CREATE SCHEMA module_menu;

CREATE TABLE module_menu.categories (
  id         UUID PRIMARY KEY,
  project_id UUID NOT NULL,
  position   INTEGER NOT NULL,

  UNIQUE (project_id, position)
);

CREATE TABLE module_menu.category_translations (
  category_id UUID NOT NULL,
  lang_code   VARCHAR(5) NOT NULL,
  name        VARCHAR(255) NOT NULL,

  UNIQUE (category_id, lang_code),

  FOREIGN KEY (category_id) REFERENCES module_menu.categories (id) ON DELETE CASCADE
);

CREATE TABLE module_menu.dishes (
  id             UUID PRIMARY KEY,
  category_id    UUID NOT NULL,
  position       INTEGER  NOT NULL,
  photo_url      VARCHAR(512) NOT NULL,
  photo_mini_url VARCHAR(512) NOT NULL,
  price          INTEGER NOT NULL,
  is_available   BOOLEAN NOT NULL,

  FOREIGN KEY (category_id) REFERENCES module_menu.categories (id)
);

CREATE TABLE module_menu.dish_translations (
  dish_id     UUID NOT NULL,
  lang_code   VARCHAR(5) NOT NULL,
  name        VARCHAR(255) NOT NULL,
  description VARCHAR(1024) NOT NULL,

  UNIQUE (dish_id, lang_code),

  FOREIGN KEY (dish_id) REFERENCES module_menu.dishes (id) ON DELETE CASCADE
);

CREATE TABLE module_menu.option_groups (
  id           UUID PRIMARY KEY,
  dish_id      UUID NOT NULL,
  position     INTEGER NOT NULL,
  min_quantity INTEGER NOT NULL,
  max_quantity INTEGER NOT NULL,

  FOREIGN KEY (dish_id) REFERENCES module_menu.dishes (id) ON DELETE CASCADE
);

CREATE TABLE module_menu.option_group_translations (
  option_group_id UUID NOT NULL,
  lang_code       VARCHAR(5) NOT NULL,
  name            VARCHAR(255) NOT NULL,

  UNIQUE (option_group_id, lang_code),

  FOREIGN KEY (option_group_id) REFERENCES module_menu.option_groups (id) ON DELETE CASCADE
);

CREATE TABLE module_menu.option_items (
  id              UUID PRIMARY KEY,
  option_group_id UUID NOT NULL,
  position        INTEGER NOT NULL,
  min_quantity    INTEGER NOT NULL,
  max_quantity    INTEGER NOT NULL,
  price_modifier  INTEGER NOT NULL,

  FOREIGN KEY (option_group_id) REFERENCES module_menu.option_groups (id) ON DELETE CASCADE
);

CREATE TABLE module_menu.option_item_translations (
  item_id   UUID NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  name      VARCHAR(255) NOT NULL,

  FOREIGN KEY (item_id) REFERENCES module_menu.option_items (id) ON DELETE CASCADE
);