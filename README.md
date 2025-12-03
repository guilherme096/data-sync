# data-sync

## Overview

Trino distributed query engine setup with 1 coordinator and 3 workers, connecting to multiple databases (PostgreSQL, MySQL, and MongoDB) for cross-database federated queries.


## Getting Started

### Start all services

```bash
docker-compose up -d
```

### Stop all services

```bash
docker-compose down
```

### Connect to Trino

```bash
docker exec -it trino-coordinator trino
```

### Query Examples

#### List available catalogs
```sql
SHOW CATALOGS;
```

#### Query on PostgreSQL
```sql
SELECT * FROM postgresql.public.customers LIMIT 5;
```

#### Query on MySQL
```sql
SELECT * FROM mysql.testdb.products LIMIT 5;
```

#### Query on MongoDB
```sql
SELECT * FROM mongodb.testdb.reviews LIMIT 5;
```

#### Cross-Database Query (Combining all 3 databases)
```sql
SELECT
    c.name as customer_name,
    c.country as customer_country,
    o.product_name as ordered_product,
    o.quantity,
    o.price as order_price,
    p.name as product_full_name,
    p.category,
    p.stock,
    r.rating,
    r.comment
FROM postgresql.public.orders o
JOIN postgresql.public.customers c ON o.customer_id = c.id
CROSS JOIN mysql.testdb.products p
LEFT JOIN mongodb.testdb.reviews r ON p.id = r.product_id
WHERE (
    LOWER(p.name) LIKE '%' || LOWER(o.product_name) || '%'
    OR LOWER(o.product_name) LIKE '%' || LOWER(SPLIT_PART(p.name, ' ', 1)) || '%'
)
ORDER BY c.name, r.rating DESC NULLS LAST
LIMIT 10;
```

### Add New Catalog Dynamically

```sql
-- Example: add new PostgreSQL connector
CREATE CATALOG newdb USING postgresql
WITH (
    "connection-url" = 'jdbc:postgresql://host:5432/database',
    "connection-user" = 'user',
    "connection-password" = 'password'
);
```


