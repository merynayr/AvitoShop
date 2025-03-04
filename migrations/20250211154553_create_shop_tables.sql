-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  coins INTEGER NOT NULL DEFAULT 1000
);

CREATE TABLE IF NOT EXISTS transactions (
  id SERIAL PRIMARY KEY,
  from_user_id INTEGER NOT NULL,
  to_user_id INTEGER NOT NULL,
  amount INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT fk_transactions_from FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_transactions_to FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS inventory (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  item_name TEXT NOT NULL,
  quantity INTEGER NOT NULL,
  CONSTRAINT fk_inventory_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Добавляем таблицу для хранения цен на мерч
CREATE TABLE IF NOT EXISTS merch_prices (
  item_name TEXT PRIMARY KEY,
  price INTEGER NOT NULL
);

-- Вставляем цены на мерч
INSERT INTO merch_prices (item_name, price) VALUES
  ('t-shirt', 80),
  ('cup', 20),
  ('book', 50),
  ('pen', 10),
  ('powerbank', 200),
  ('hoody', 300),
  ('umbrella', 200),
  ('socks', 10),
  ('wallet', 50),
  ('pink-hoody', 500);

CREATE INDEX IF NOT EXISTS idx_transactions_from ON transactions (from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to ON transactions (to_user_id);
CREATE INDEX IF NOT EXISTS idx_inventory_user ON inventory (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_item ON inventory (user_id, item_name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS merch_prices;
-- +goose StatementEnd
