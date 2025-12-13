-- PostgreSQL - US Region Data

-- Clients table
CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    country VARCHAR(50) NOT NULL,
    region VARCHAR(20) DEFAULT 'US',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    client_id INTEGER REFERENCES clients(id),
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    region VARCHAR(20) DEFAULT 'US'
);

-- Insert US region clients
INSERT INTO clients (name, email, country) VALUES
    ('John Smith', 'john.smith@example.com', 'USA'),
    ('Sarah Johnson', 'sarah.j@example.com', 'USA'),
    ('Michael Davis', 'michael.d@example.com', 'Canada'),
    ('Emily Wilson', 'emily.w@example.com', 'USA'),
    ('David Brown', 'david.b@example.com', 'Canada');

-- Insert US region orders
INSERT INTO orders (client_id, order_date, total_amount, status) VALUES
    (1, '2024-11-15 10:30:00', 1245.50, 'completed'),
    (1, '2024-11-28 14:20:00', 89.99, 'completed'),
    (2, '2024-11-18 09:15:00', 450.00, 'completed'),
    (2, '2024-12-01 16:45:00', 1200.00, 'pending'),
    (3, '2024-11-20 11:00:00', 750.00, 'completed'),
    (4, '2024-11-25 13:30:00', 120.00, 'shipped'),
    (4, '2024-12-02 10:00:00', 399.99, 'pending'),
    (5, '2024-11-22 15:20:00', 45.00, 'completed');
