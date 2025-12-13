USE testdb;

-- MySQL - EU Region Data

-- Clients table
CREATE TABLE clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    country VARCHAR(50) NOT NULL,
    region VARCHAR(20) DEFAULT 'EU',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    client_id INT,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    region VARCHAR(20) DEFAULT 'EU',
    FOREIGN KEY (client_id) REFERENCES clients(id)
);

-- Insert EU region clients
INSERT INTO clients (name, email, country) VALUES
    ('João Silva', 'joao.silva@example.com', 'Portugal'),
    ('Maria Santos', 'maria.santos@example.com', 'Portugal'),
    ('Hans Mueller', 'hans.m@example.com', 'Germany'),
    ('Sophie Dubois', 'sophie.d@example.com', 'France'),
    ('Marco Rossi', 'marco.r@example.com', 'Italy'),
    ('Anna Kowalski', 'anna.k@example.com', 'Poland');

-- Insert EU region orders
INSERT INTO orders (client_id, order_date, total_amount, status) VALUES
    (1, '2024-11-16 11:00:00', 899.00, 'completed'),
    (1, '2024-11-29 15:30:00', 150.00, 'completed'),
    (2, '2024-11-19 10:45:00', 1500.00, 'completed'),
    (3, '2024-11-21 14:15:00', 650.00, 'shipped'),
    (3, '2024-12-03 09:30:00', 220.00, 'pending'),
    (4, '2024-11-26 12:00:00', 1100.00, 'completed'),
    (5, '2024-11-27 16:20:00', 85.00, 'completed'),
    (5, '2024-12-01 11:45:00', 450.00, 'shipped'),
    (6, '2024-11-23 13:50:00', 320.00, 'completed');
