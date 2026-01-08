<div align="center">
  <img src="docs/data-sync.png" alt="data-sync logo" width="200"/>
</div>

# data-sync

A distributed data integration platform that enables federated querying across multiple databases (PostgreSQL, MySQL, MongoDB) using Trino's distributed query engine. Features automatic metadata discovery, a REST API for catalog management, and a web interface for interactive query execution.

## Environment Configuration

Create a `.env` file in the root directory with the following variables:

```env
GEMINI_API_KEY=your_gemini_api_key_here
```

**GEMINI_API_KEY**: Your Google Gemini API key for AI-powered features (required).

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
