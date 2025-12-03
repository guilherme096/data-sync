USE testdb;

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE suppliers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    contact_email VARCHAR(100),
    country VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (name, category, price, stock) VALUES
    ('Laptop Dell XPS', 'Electronics', 1500.00, 25),
    ('Mouse Logitech', 'Electronics', 30.00, 150),
    ('Keyboard Mechanical', 'Electronics', 120.00, 80),
    ('Monitor 27"', 'Electronics', 400.00, 45),
    ('USB-C Cable', 'Accessories', 20.00, 200),
    ('Webcam HD', 'Electronics', 85.00, 60);

INSERT INTO suppliers (name, contact_email, country) VALUES
    ('Tech Supply Ltd', 'contact@techsupply.com', 'USA'),
    ('Euro Components', 'info@eurocomp.eu', 'Germany'),
    ('Asia Electronics', 'sales@asiaelec.com', 'China'),
    ('Local Parts PT', 'vendas@localparts.pt', 'Portugal');
