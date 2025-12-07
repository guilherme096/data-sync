# Frontend Architecture Plan: Data Sync & Schema Matcher

## Core Vision
A system to unify distributed data sources into "Imaginary Global Tables" that can be queried seamlessly via SQL or Natural Language. The UI prioritizes the consumption of data (Querying) while providing sophisticated tools for mapping raw data to logical schemas.

## App Structure

### 1. Home / Query Interface (`/`)
**Purpose:** The entry point for answering questions. Acts as a "Data Search Engine".
*   **Agent Mode:** Natural language input ("Show me high-value users").
*   **SQL Mode:** Advanced raw query editor.
*   **Context:** Toggle between querying "Global Tables" (Mapped) vs "Raw Catalogs" (Direct).
*   **Zero State:** Prompts user to configure Schema/Inventory if no global tables exist.

### 2. Schema Studio (`/schema`)
**Purpose:** The workspace for defining the "Global Tables" and logic.
*   **Global Table Builder:** Define virtual tables (e.g., `Unified_Customers`).
*   **Mapping Interface:** Drag-and-drop or selection interface to map Raw Columns -> Global Columns.
*   **Relation Definer:** UI to establish foreign key-like relationships across disparate sources.

### 3. Data Inventory (`/inventory`)
**Purpose:** The "Raw View" of connected physical sources.
*   **Catalog Browser:** Hierarchical view of connected sources (Postgres, Mongo, etc.).
*   **Sync Status:** Manual trigger to refresh metadata from the backend.
*   **Metadata Viewer:** Inspect raw schemas and types.

## Technical Stack
*   **Framework:** React + Vite
*   **UI Library:** Shadcn UI (Tailwind CSS)
*   **State Management:** React Query (Server state), Zustand or React Context (Client UI state).
*   **Routing:** React Router DOM
*   **Editor:** Monaco Editor or CodeMirror (for SQL).

## Directory Structure
```text
src/
├── features/
│   ├── inventory/          # Source visualization
│   ├── schema-studio/      # Mapping logic
│   └── query/              # Home page logic (Chat + SQL)
├── layouts/
│   └── MainLayout.tsx      # Sidebar + Shell
├── lib/
│   └── api.ts              # Typed backend fetchers
├── pages/
│   ├── QueryPage.tsx       # Home
│   ├── SchemaPage.tsx
│   └── InventoryPage.tsx
└── components/ui/          # Shadcn primitives
```
