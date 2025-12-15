<div align="center">
  <img src="docs/data-sync.png" alt="data-sync logo" width="200"/>
</div>

# data-sync

A distributed data integration platform that enables federated querying across multiple databases (PostgreSQL, MySQL, MongoDB) using Trino's distributed query engine. Features automatic metadata discovery, a REST API for catalog management, and a web interface for interactive query execution.

## Quickstart

Start all services (data sources, Trino cluster, backend API, and frontend):
```bash
./start.sh
```

Stop all services:
```bash
./stop.sh
```

Once running, access:
- **Frontend**: http://localhost:5173
- **API**: http://localhost:8081
- **Trino**: http://localhost:8080
