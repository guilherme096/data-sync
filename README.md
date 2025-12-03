# data-sync

## Overview

Simple setup for Trino to connect to multiple databases (PostgreSQL and MySQL) and perform cross-database queries


## How to use 

### Connect to Trino

```bash
docker exec -it trino trino
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

#### Cross-Database Query (JOIN between PostgreSQL and MySQL)
```sql
SELECT
    c.name as customer_name,
    c.country,
    p.name as product_name,
    p.category
FROM postgresql.public.customers c
CROSS JOIN mysql.testdb.products p
WHERE c.country = 'Portugal' AND p.category = 'Electronics'
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


