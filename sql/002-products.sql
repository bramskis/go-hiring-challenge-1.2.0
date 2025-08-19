CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32),
    price DECIMAL(10, 2) NOT NULL,
    -- NOT NULL if product must have a category, but nullable and ON DELETE SET NULL if product does not have to have a category
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
