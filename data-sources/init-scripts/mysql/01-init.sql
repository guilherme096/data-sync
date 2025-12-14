USE testdb;

-- MySQL - CRM/Enhanced System (Extended Customer & Product Data)

-- Customer Profiles table - Extended customer information from CRM
-- Links to PostgreSQL customers via customer_id
CREATE TABLE customer_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,  -- References postgres.customers.id
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(50),
    postal_code VARCHAR(20),
    loyalty_tier VARCHAR(20) DEFAULT 'Bronze',
    total_lifetime_spent DECIMAL(12, 2) DEFAULT 0,
    preferred_payment_method VARCHAR(30),
    last_purchase_date TIMESTAMP,
    UNIQUE KEY unique_customer (customer_id)
);

-- Product Inventory table - Warehouse and supplier information
-- Links to PostgreSQL products via product_id
CREATE TABLE product_inventory (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,  -- References postgres.products.id
    stock_level INT NOT NULL DEFAULT 0,
    warehouse_location VARCHAR(50),
    supplier VARCHAR(100),
    reorder_point INT DEFAULT 10,
    last_restock_date TIMESTAMP,
    unit_cost DECIMAL(10, 2),
    UNIQUE KEY unique_product (product_id)
);

-- Orders 2024 table - Recent orders (same schema as postgres.orders_2023)
CREATE TABLE orders_2024 (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer_id INT NOT NULL,
    product_id INT NOT NULL,
    order_date TIMESTAMP NOT NULL,
    quantity INT NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    year INT DEFAULT 2024
);

-- Insert customer profiles (matching postgres customer IDs)
INSERT INTO customer_profiles (customer_id, phone, address, city, postal_code, loyalty_tier, total_lifetime_spent, preferred_payment_method, last_purchase_date) VALUES
    (1, '+1-555-0101', '123 Main St', 'New York', '10001', 'Gold', 2850.50, 'Credit Card', '2024-11-15'),
    (2, '+1-555-0102', '456 Oak Ave', 'Los Angeles', '90001', 'Silver', 1890.25, 'PayPal', '2024-10-20'),
    (3, '+1-555-0103', '789 Pine Rd', 'Toronto', 'M5H 2N2', 'Bronze', 650.00, 'Debit Card', '2024-09-15'),
    (4, '+1-555-0104', '321 Elm St', 'Chicago', '60601', 'Platinum', 4200.00, 'Credit Card', '2024-12-28'),
    (5, '+1-555-0105', '654 Maple Dr', 'Vancouver', 'V6B 1A1', 'Silver', 1250.75, 'Credit Card', '2024-11-30'),
    (6, '+1-555-0106', '987 Cedar Ln', 'Mexico City', '06600', 'Gold', 3100.00, 'Cash', '2024-12-10'),
    (7, '+1-555-0107', '147 Birch Ct', 'Seattle', '98101', 'Bronze', 890.50, 'PayPal', '2024-11-20'),
    (8, '+1-555-0108', '258 Spruce Way', 'Montreal', 'H3A 1A1', 'Silver', 1650.00, 'Debit Card', '2024-12-05');

-- Insert product inventory (matching postgres product IDs)
INSERT INTO product_inventory (product_id, stock_level, warehouse_location, supplier, reorder_point, last_restock_date, unit_cost) VALUES
    (1, 45, 'Warehouse A - NY', 'TechSupply Inc', 15, '2024-11-01', 899.99),
    (2, 230, 'Warehouse B - CA', 'Peripherals Co', 50, '2024-11-15', 12.50),
    (3, 18, 'Warehouse C - TX', 'Furniture Direct', 5, '2024-10-20', 210.00),
    (4, 12, 'Warehouse C - TX', 'Furniture Direct', 5, '2024-10-20', 180.00),
    (5, 450, 'Warehouse B - CA', 'Cable Masters', 100, '2024-12-01', 5.99),
    (6, 38, 'Warehouse A - NY', 'Display Tech Ltd', 10, '2024-11-10', 285.00),
    (7, 67, 'Warehouse B - CA', 'Peripherals Co', 20, '2024-11-20', 45.00),
    (8, 95, 'Warehouse C - TX', 'Furniture Direct', 25, '2024-10-15', 22.50);

-- Insert 2024 orders
INSERT INTO orders_2024 (customer_id, product_id, order_date, quantity, total_amount, status) VALUES
    (1, 6, '2024-11-16 11:00:00', 1, 399.99, 'completed'),
    (2, 7, '2024-10-19 15:30:00', 1, 89.99, 'completed'),
    (2, 5, '2024-10-19 15:30:00', 2, 39.98, 'completed'),
    (3, 2, '2024-09-21 14:15:00', 3, 89.97, 'completed'),
    (4, 1, '2024-12-26 09:30:00', 1, 1299.99, 'shipped'),
    (4, 7, '2024-12-26 09:30:00', 1, 89.99, 'shipped'),
    (5, 3, '2024-11-28 12:00:00', 1, 349.99, 'completed'),
    (5, 8, '2024-11-28 12:00:00', 1, 45.99, 'completed'),
    (6, 1, '2024-12-08 16:20:00', 2, 2599.98, 'completed'),
    (7, 4, '2024-11-18 13:50:00', 1, 299.99, 'completed'),
    (8, 6, '2024-12-03 10:30:00', 1, 399.99, 'shipped'),
    (8, 2, '2024-12-03 10:30:00', 2, 59.98, 'shipped');
