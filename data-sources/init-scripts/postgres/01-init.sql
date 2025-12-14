-- PostgreSQL - Core System (Basic Customer & Product Data)

-- Customers table - Basic customer information
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    country VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table - Basic product catalog
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    category VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders 2023 table - Historical orders
CREATE TABLE orders_2023 (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    product_id INTEGER REFERENCES products(id),
    order_date TIMESTAMP NOT NULL,
    quantity INTEGER NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'completed',
    year INTEGER DEFAULT 2023
);

-- Insert customers
INSERT INTO customers (name, email, country) VALUES
    ('John Smith', 'john.smith@example.com', 'USA'),
    ('Sarah Johnson', 'sarah.j@example.com', 'USA'),
    ('Michael Davis', 'michael.d@example.com', 'Canada'),
    ('Emily Wilson', 'emily.w@example.com', 'USA'),
    ('David Brown', 'david.b@example.com', 'Canada'),
    ('Maria Garcia', 'maria.g@example.com', 'Mexico'),
    ('James Anderson', 'james.a@example.com', 'USA'),
    ('Lisa Chen', 'lisa.c@example.com', 'Canada');

-- Insert products
INSERT INTO products (name, price, category) VALUES
    ('Laptop Pro 15"', 1299.99, 'Electronics'),
    ('Wireless Mouse', 29.99, 'Electronics'),
    ('Office Desk', 349.99, 'Furniture'),
    ('Gaming Chair', 299.99, 'Furniture'),
    ('USB-C Cable', 19.99, 'Electronics'),
    ('Monitor 27"', 399.99, 'Electronics'),
    ('Keyboard Mechanical', 89.99, 'Electronics'),
    ('Desk Lamp', 45.99, 'Furniture');

-- Insert 2023 orders
INSERT INTO orders_2023 (customer_id, product_id, order_date, quantity, total_amount, status) VALUES
    (1, 1, '2023-11-15 10:30:00', 1, 1299.99, 'completed'),
    (1, 2, '2023-11-15 10:30:00', 2, 59.98, 'completed'),
    (2, 3, '2023-10-18 09:15:00', 1, 349.99, 'completed'),
    (2, 8, '2023-10-18 09:15:00', 1, 45.99, 'completed'),
    (3, 4, '2023-09-20 11:00:00', 1, 299.99, 'completed'),
    (4, 6, '2023-12-25 13:30:00', 1, 399.99, 'completed'),
    (4, 7, '2023-12-25 13:30:00', 1, 89.99, 'completed'),
    (5, 5, '2023-08-22 15:20:00', 3, 59.97, 'completed'),
    (6, 1, '2023-07-10 14:00:00', 1, 1299.99, 'completed'),
    (7, 2, '2023-11-05 16:45:00', 5, 149.95, 'completed');
