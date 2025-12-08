export type Catalog = {
  Name: string;
  Metadata: Record<string, string>;
};

export type Schema = {
  Name: string;
  CatalogName: string;
  Metadata: Record<string, string>;
};

export type Table = {
  Name: string;
  SchemaName: string;
  CatalogName: string;
  Metadata: Record<string, string>;
};

export type Column = {
  Name: string;
  TableName: string;
  SchemaName: string;
  CatalogName: string;
  DataType: string;
  Metadata: Record<string, string>;
};

export type QueryResult = {
  Rows: Record<string, unknown>[] | null;
};

export type SyncResponse = {
  status: string;
  message: string;
};

const API_BASE = '/api';

export const api = {
  health: async () => {
    const res = await fetch(`${API_BASE}/health`);
    if (!res.ok) throw new Error('Health check failed');
    return res.json();
  },

  // Temporarily excluding executeQuery as per user's request
  // executeQuery: async (query: string, params: Record<string, unknown> = {}): Promise<QueryResult> => {
  //   const res = await fetch(`${API_BASE}/query`, {
  //     method: 'POST',
  //     headers: { 'Content-Type': 'application/json' },
  //     body: JSON.stringify({ query, params }),
  //   });
      
  //   if (!res.ok) {
  //       const errText = await res.text();
  //       throw new Error(errText || 'Query failed');
  //   }
      
  //   return res.json();
  // },

  listCatalogs: async (): Promise<Catalog[]> => {
    const res = await fetch(`${API_BASE}/catalogs`);
    if (!res.ok) throw new Error('Failed to fetch catalogs');
    return res.json();
  },

  getCatalog: async (name: string): Promise<Catalog> => {
    const res = await fetch(`${API_BASE}/catalogs/${name}`);
    if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || `Failed to fetch catalog ${name}`);
    }
    return res.json();
  },

  listSchemas: async (catalogName: string): Promise<Schema[]> => {
    const res = await fetch(`${API_BASE}/catalogs/${catalogName}/schemas`);
    if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || `Failed to fetch schemas for catalog ${catalogName}`);
    }
    return res.json();
  },

  syncMetadata: async (): Promise<SyncResponse> => {
    const res = await fetch(`${API_BASE}/sync`, { method: 'POST' });
    if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || 'Sync failed');
    }
    return res.json();
  },

  discoverTables: async (catalogName: string, schemaName: string): Promise<Table[]> => {
    const res = await fetch(`${API_BASE}/discover/catalogs/${catalogName}/schemas/${schemaName}/tables`);
    if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || `Failed to discover tables for ${catalogName}.${schemaName}`);
    }
    return res.json();
  },

  discoverColumns: async (catalogName: string, schemaName: string, tableName: string): Promise<Column[]> => {
    const res = await fetch(`${API_BASE}/discover/catalogs/${catalogName}/schemas/${schemaName}/tables/${tableName}/columns`);
    if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || `Failed to discover columns for ${catalogName}.${schemaName}.${tableName}`);
    }
    return res.json();
  }
};
