-- Создание таблицы пользователей (users)
CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) PRIMARY KEY,
    hashed_password TEXT NOT NULL,
    balance BIGINT NOT NULL
);

-- Создание таблицы мерча (merch)
CREATE TABLE IF NOT EXISTS merch (
    name VARCHAR(255) PRIMARY KEY,
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
    username varchar(255) NOT NULL,
    merch_item varchar(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (username) REFERENCES users (username),
    FOREIGN KEY (merch_item) REFERENCES merch (name)
);

-- Создание таблицы для логирования транзакций с монетами (transaction_log)
-- Это таблица, которая будет содержать как покупки мерча, так и переводы монет между пользователями
CREATE TABLE IF NOT EXISTS transaction_log (
    id SERIAL PRIMARY KEY,
    sender VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    amount BIGINT NOT NULL,
    type VARCHAR(50) NOT NULL,  -- Тип транзакции: 'purchase' / 'transfer'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
