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

export type GlobalTable = {
  Name: string;
  Description: string;
};

export type GlobalColumn = {
  GlobalTableName: string;
  Name: string;
  DataType: string;
  Description: string;
};

export type TableMapping = {
  GlobalTableName: string;
  CatalogName: string;
  SchemaName: string;
  TableName: string;
};

export type ColumnMapping = {
  GlobalTableName: string;
  GlobalColumnName: string;
  CatalogName: string;
  SchemaName: string;
  TableName: string;
  ColumnName: string;
};

export type ColumnRelationship = {
  SourceGlobalTableName: string;
  SourceGlobalColumnName: string;
  TargetGlobalTableName: string;
  TargetGlobalColumnName: string;
  RelationshipName?: string;
  Description?: string;
};

export type TableSource = {
  type: 'physical' | 'relation';
  catalog?: string;
  schema?: string;
  table?: string;
  relationId?: string;
};

export type JoinColumn = {
  left: string;
  right: string;
};

export type TableRelation = {
  id: string;
  name: string;
  leftTable: TableSource;
  rightTable: TableSource;
  relationType: 'JOIN' | 'UNION';
  joinColumn?: JoinColumn;
  description?: string;
};

export type QueryResult = {
  Rows: Record<string, unknown>[] | null;
};

export type ToolResult = {
  toolName: string;
  data: any;
};

export type ChatResponse = {
  message: string;
  toolResults?: ToolResult[];
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
  },

  // Global Tables API
  listGlobalTables: async (): Promise<GlobalTable[]> => {
    const res = await fetch(`${API_BASE}/global/tables`);
    if (!res.ok) throw new Error('Failed to fetch global tables');
    return res.json();
  },

  getGlobalTable: async (name: string): Promise<GlobalTable> => {
    const res = await fetch(`${API_BASE}/global/tables/${name}`);
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || `Failed to fetch global table ${name}`);
    }
    return res.json();
  },

  createGlobalTable: async (table: { Name: string; Description: string }): Promise<GlobalTable> => {
    const res = await fetch(`${API_BASE}/global/tables`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(table),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create global table');
    }
    return res.json();
  },

  deleteGlobalTable: async (name: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/global/tables/${name}`, {
      method: 'DELETE',
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete global table');
    }
  },

  // Global Columns API
  listGlobalColumns: async (globalTableName: string): Promise<GlobalColumn[]> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns`);
    if (!res.ok) throw new Error(`Failed to fetch columns for ${globalTableName}`);
    return res.json();
  },

  createGlobalColumn: async (globalTableName: string, column: { Name: string; DataType: string; Description: string }): Promise<GlobalColumn> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(column),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create global column');
    }
    return res.json();
  },

  deleteGlobalColumn: async (globalTableName: string, columnName: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns/${columnName}`, {
      method: 'DELETE',
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete global column');
    }
  },

  // Table Mappings API
  listTableMappings: async (globalTableName: string): Promise<TableMapping[]> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/mappings/tables`);
    if (!res.ok) throw new Error(`Failed to fetch table mappings for ${globalTableName}`);
    return res.json();
  },

  createTableMapping: async (globalTableName: string, mapping: { CatalogName: string; SchemaName: string; TableName: string }): Promise<TableMapping> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/mappings/tables`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(mapping),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create table mapping');
    }
    return res.json();
  },

  deleteTableMapping: async (globalTableName: string, mapping: { CatalogName: string; SchemaName: string; TableName: string }): Promise<void> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/mappings/tables`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(mapping),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete table mapping');
    }
  },

  // Column Mappings API
  listColumnMappings: async (globalTableName: string, globalColumnName: string): Promise<ColumnMapping[]> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns/${globalColumnName}/mappings`);
    if (!res.ok) throw new Error(`Failed to fetch column mappings`);
    return res.json();
  },

  createColumnMapping: async (globalTableName: string, globalColumnName: string, mapping: { CatalogName: string; SchemaName: string; TableName: string; ColumnName: string }): Promise<ColumnMapping> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns/${globalColumnName}/mappings`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(mapping),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create column mapping');
    }
    return res.json();
  },

  deleteColumnMapping: async (globalTableName: string, globalColumnName: string, mapping: { CatalogName: string; SchemaName: string; TableName: string; ColumnName: string }): Promise<void> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/columns/${globalColumnName}/mappings`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(mapping),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete column mapping');
    }
  },

  // Column Relationships API
  listColumnRelationships: async (globalTableName: string): Promise<ColumnRelationship[]> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/relationships`);
    if (!res.ok) throw new Error(`Failed to fetch relationships for ${globalTableName}`);
    return res.json();
  },

  createColumnRelationship: async (
    globalTableName: string,
    relationship: {
      SourceGlobalTableName: string;
      SourceGlobalColumnName: string;
      TargetGlobalTableName: string;
      TargetGlobalColumnName: string;
      RelationshipName?: string;
      Description?: string;
    }
  ): Promise<ColumnRelationship> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/relationships`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(relationship),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create column relationship');
    }
    return res.json();
  },

  deleteColumnRelationship: async (
    globalTableName: string,
    relationship: {
      SourceGlobalTableName: string;
      SourceGlobalColumnName: string;
      TargetGlobalTableName: string;
      TargetGlobalColumnName: string;
    }
  ): Promise<void> => {
    const res = await fetch(`${API_BASE}/global/tables/${globalTableName}/relationships`, {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(relationship),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete column relationship');
    }
  },

  // Chat API
  sendChatMessage: async (message: string, conversationHistory: unknown[]): Promise<ChatResponse> => {
    // Map conversation history to backend format (role and content only)
    const history = (conversationHistory as Array<{ role: string; content: string }>).map(m => ({
      role: m.role,
      content: m.content,
    }));

    const res = await fetch(`${API_BASE}/chatbot/message`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ message, history }),
    });

    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to send chat message');
    }

    const data = await res.json();
    return data; // Return full ChatResponse with message and toolResults
  },

  // Table Relations API
  listTableRelations: async (): Promise<TableRelation[]> => {
    const res = await fetch(`${API_BASE}/relations`);
    if (!res.ok) throw new Error('Failed to fetch table relations');
    return res.json();
  },

  getTableRelation: async (id: string): Promise<TableRelation> => {
    const res = await fetch(`${API_BASE}/relations/${id}`);
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || `Failed to fetch relation ${id}`);
    }
    return res.json();
  },

  createTableRelation: async (relation: TableRelation): Promise<TableRelation> => {
    const res = await fetch(`${API_BASE}/relations`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(relation),
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to create table relation');
    }
    return res.json();
  },

  deleteTableRelation: async (id: string): Promise<void> => {
    const res = await fetch(`${API_BASE}/relations/${id}`, {
      method: 'DELETE',
    });
    if (!res.ok) {
      const errText = await res.text();
      throw new Error(errText || 'Failed to delete table relation');
    }
  },
};
