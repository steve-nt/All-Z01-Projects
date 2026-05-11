-- Optional explicit schema
CREATE TABLE IF NOT EXISTS sales (
    transaction_id     TEXT,
    customer_id        TEXT,
    product_id         TEXT,
    product_name       TEXT,
    quantity           NUMERIC,
    price_per_unit     NUMERIC,
    total_price        NUMERIC,
    transaction_date   DATE
);
