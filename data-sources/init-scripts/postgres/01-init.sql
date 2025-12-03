CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    country VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER REFERENCES customers(id),
    product_name VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO customers (name, email, country) VALUES
    ('João Silva', 'joao@example.com', 'Portugal'),
    ('Maria Santos', 'maria@example.com', 'Portugal'),
    ('Pedro Costa', 'pedro@example.com', 'Brasil'),
    ('Ana Rodrigues', 'ana@example.com', 'Portugal'),
    ('Carlos Ferreira', 'carlos@example.com', 'Brasil');

INSERT INTO orders (customer_id, product_name, quantity, price) VALUES
    (1, 'Laptop', 1, 1200.00),
    (1, 'Mouse', 2, 25.50),
    (2, 'Keyboard', 1, 89.99),
    (3, 'Monitor', 2, 350.00),
    (4, 'Headphones', 1, 120.00),
    (5, 'USB Cable', 3, 15.00);
