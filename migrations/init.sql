-- Создание таблицы пользователей (users)
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    balance BIGINT NOT NULL DEFAULT 1000,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы мерча (merch)
CREATE TABLE IF NOT EXISTS merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    price BIGINT NOT NULL
);

-- Вставка начальных товаров (мерча) в магазин
INSERT INTO merch (name, price) VALUES
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

-- Создание таблицы покупок (purchases) для хранения покупок пользователей
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    merch_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (merch_id) REFERENCES merch (id)
);

-- Создание таблицы для логирования транзакций с монетами (transaction_log)
-- Это таблица, которая будет содержать как покупки мерча, так и переводы монет между пользователями
CREATE TABLE IF NOT EXISTS transaction_log (
    id SERIAL PRIMARY KEY,
    source_id INT,
    destination_id INT,
    amount BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,  -- Тип транзакции: 'purchase' / 'transfer'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
