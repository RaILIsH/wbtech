-- Создание таблицы orders
CREATE TABLE IF NOT EXISTS orders (
    order_uid          VARCHAR(255) PRIMARY KEY,
    track_number       VARCHAR(255),
    entry              VARCHAR(255),
    locale             VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id        VARCHAR(255),
    delivery_service   VARCHAR(255),
    shardkey           VARCHAR(255),
    sm_id              INTEGER,
    date_created       TIMESTAMP,
    oof_shard          VARCHAR(255)
);

-- Создание таблицы delivery
CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    name      VARCHAR(255) NOT NULL,
    phone     VARCHAR(255) NOT NULL,
    zip       VARCHAR(255) NOT NULL,
    city      VARCHAR(255) NOT NULL,
    address   VARCHAR(255) NOT NULL,
    region    VARCHAR(255) NOT NULL,
    email     VARCHAR(255) NOT NULL,
    PRIMARY KEY (order_uid)
);

-- Создание таблицы payment
CREATE TABLE IF NOT EXISTS payment (
    order_uid     VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction   VARCHAR(255) NOT NULL,
    request_id    VARCHAR(255),
    currency      VARCHAR(3) NOT NULL,
    provider      VARCHAR(255) NOT NULL,
    amount        INTEGER NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          VARCHAR(255) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total   INTEGER NOT NULL,
    custom_fee    INTEGER DEFAULT 0,
    PRIMARY KEY (order_uid)
);

-- Создание таблицы items
CREATE TABLE IF NOT EXISTS items (
    id           SERIAL PRIMARY KEY,
    order_uid    VARCHAR(255) REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id      INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price        INTEGER NOT NULL,
    rid          VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    sale         INTEGER NOT NULL,
    size         VARCHAR(255) NOT NULL,
    total_price  INTEGER NOT NULL,
    nm_id        INTEGER NOT NULL,
    brand        VARCHAR(255) NOT NULL,
    status       INTEGER NOT NULL
);

-- Индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_orders_order_uid ON orders(order_uid);
CREATE INDEX IF NOT EXISTS idx_items_order_uid ON items(order_uid);