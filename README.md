# data-sync

## Overview

Trino distributed query engine setup with 1 coordinator and 3 workers for cross-database federated queries.


## Getting Started

### 1. Start Trino cluster

This creates the `trino-network` that data sources will join:

```bash
docker-compose up -d
```

### 2. (Optional) Start test data sources

**Important:** Trino must be started first to create the shared network.

```bash
cd data-sources
docker-compose up -d
cd ..
```

### 3. Connect to Trino

```bash
docker exec -it trino-coordinator trino
```


### Dynamic Catalog Management

You can also add catalogs dynamically via SQL (requires `catalog.management=dynamic` in config):

```sql
CREATE CATALOG newdb USING postgresql
WITH (
    "connection-url" = 'jdbc:postgresql://host:5432/database',
    "connection-user" = 'user',
    "connection-password" = 'password'
);
```

## Query Examples

### List available catalogs
```sql
SHOW CATALOGS;
```

### Query specific data source
```sql
-- PostgreSQL
SELECT * FROM postgresql.public.customers LIMIT 5;

-- MySQL
SELECT * FROM mysql.testdb.products LIMIT 5;

-- MongoDB
SELECT * FROM mongodb.testdb.reviews LIMIT 5;
```

### Cross-Database Query
```sql
SELECT
    c.name as customer_name,
    c.country as customer_country,
    o.product_name as ordered_product,
    p.name as product_full_name,
    p.category,
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
