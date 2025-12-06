# data-sync

## Overview

Metadata sync API service for Trino - automatically discovers and caches catalogs and schemas for fast access via REST API.

Includes Trino distributed query engine setup with 1 coordinator and 3 workers for cross-database federated queries.


## Getting Started

### 1. Start all services (Trino + data-sync API)

```bash
docker-compose up -d --build
```

This starts:
- Trino coordinator (port 8080)
- 3 Trino workers
- data-sync API service (port 8081)



### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/catalogs` | List all cached catalogs |
| GET | `/catalogs/{name}` | Get specific catalog details |
| GET | `/catalogs/{name}/schemas` | List schemas for a catalog |
| POST | `/sync` | Manually trigger metadata sync |
| POST | `/query` | Execute Trino query directly |

### Examples

#### List all catalogs
```bash
curl http://localhost:8081/catalogs
```

Response:
```json
[
  {"Name": "mongodb", "Metadata": {}},
  {"Name": "mysql", "Metadata": {}},
  {"Name": "postgresql", "Metadata": {}},
  {"Name": "system", "Metadata": {}}
]
```

#### Get schemas for a catalog
```bash
curl http://localhost:8081/catalogs/system/schemas
```

Response:
```json
[
  {"Name": "information_schema", "CatalogName": "system", "Metadata": {}},
  {"Name": "jdbc", "CatalogName": "system", "Metadata": {}},
  {"Name": "metadata", "CatalogName": "system", "Metadata": {}},
  {"Name": "runtime", "CatalogName": "system", "Metadata": {}}
]
```

#### Trigger manual sync
```bash
curl -X POST http://localhost:8081/sync
```
