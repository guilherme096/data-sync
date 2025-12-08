import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { useQuery, useMutation, QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { api, type Catalog, type Schema, type Table as TableModel, type Column as ColumnModel } from "@/lib/api"
import { useState, useMemo } from "react"
import { Zap, Search } from "lucide-react"

const queryClient = new QueryClient();

function TableItem({ catalogName, schemaName, tableName }: { catalogName: string; schemaName: string; tableName: string }) {
  const { data: columns, isLoading } = useQuery<ColumnModel[], Error>({
    queryKey: ['columns', catalogName, schemaName, tableName],
    queryFn: () => api.discoverColumns(catalogName, schemaName, tableName),
  });

  return (
    <div className="border-b">
      <div className="py-2 px-3 font-medium">{tableName}</div>
      <div className="pb-2 px-3 pl-9 bg-muted/10">
        {isLoading ? (
          <div className="text-sm text-muted-foreground py-1">Loading...</div>
        ) : columns && columns.length > 0 ? (
          <div className="space-y-0.5">
            {columns.map((column) => (
              <div key={column.Name} className="flex items-baseline gap-3 text-sm py-0.5">
                <span className="min-w-[180px]">{column.Name}</span>
                <code className="text-xs bg-background px-1.5 py-0.5 rounded font-mono text-muted-foreground">
                  {column.DataType}
                </code>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-sm text-muted-foreground py-1">No columns</div>
        )}
      </div>
    </div>
  );
}

function SchemaSection({ catalogName, schemaName, tableSearch }: { catalogName: string; schemaName: string; tableSearch: string }) {
  const { data: tables, isLoading } = useQuery<TableModel[], Error>({
    queryKey: ['tables', catalogName, schemaName],
    queryFn: () => api.discoverTables(catalogName, schemaName),
  });

  const filteredTables = useMemo(() => {
    if (!tables) return [];
    if (!tableSearch) return tables;
    return tables.filter(table =>
      table.Name.toLowerCase().includes(tableSearch.toLowerCase())
    );
  }, [tables, tableSearch]);

  if (isLoading) return <div className="text-sm text-muted-foreground py-2">Loading tables...</div>;
  if (!filteredTables || filteredTables.length === 0) return null;

  return (
    <div className="mb-6">
      <div className="text-sm font-semibold text-muted-foreground mb-2 px-3">
        {catalogName} › {schemaName}
      </div>
      <div>
        {filteredTables.map((table) => (
          <TableItem
            key={table.Name}
            catalogName={catalogName}
            schemaName={schemaName}
            tableName={table.Name}
          />
        ))}
      </div>
    </div>
  );
}

function CatalogSection({ catalogName, tableSearch }: { catalogName: string; tableSearch: string }) {
  const { data: schemas, isLoading } = useQuery<Schema[], Error>({
    queryKey: ['schemas', catalogName],
    queryFn: () => api.listSchemas(catalogName),
  });

  if (isLoading) return <div className="text-sm text-muted-foreground py-2">Loading schemas...</div>;
  if (!schemas || schemas.length === 0) return null;

  return (
    <>
      {schemas.map((schema) => (
        <SchemaSection
          key={schema.Name}
          catalogName={catalogName}
          schemaName={schema.Name}
          tableSearch={tableSearch}
        />
      ))}
    </>
  );
}

function InventoryPageContent() {
  const [selectedCatalogName, setSelectedCatalogName] = useState<string | null>(null);
  const [selectedSchemaName, setSelectedSchemaName] = useState<string | null>(null);
  const [tableSearch, setTableSearch] = useState("");

  const { data: catalogs, isLoading: isLoadingCatalogs } = useQuery<Catalog[], Error>({
    queryKey: ['catalogs'],
    queryFn: api.listCatalogs,
  });

  const { data: schemas, isLoading: isLoadingSchemas } = useQuery<Schema[], Error>({
    queryKey: ['schemas', selectedCatalogName],
    queryFn: () => api.listSchemas(selectedCatalogName!),
    enabled: !!selectedCatalogName && selectedCatalogName !== '__ALL__',
  });

  const { data: tables, isLoading: isLoadingTables } = useQuery<TableModel[], Error>({
    queryKey: ['tables', selectedCatalogName, selectedSchemaName],
    queryFn: () => api.discoverTables(selectedCatalogName!, selectedSchemaName!),
    enabled: !!selectedCatalogName && selectedCatalogName !== '__ALL__' && !!selectedSchemaName && selectedSchemaName !== '__ALL__',
  });

  const syncMutation = useMutation({
    mutationFn: api.syncMetadata,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['catalogs'] });
      if (selectedCatalogName) {
        queryClient.invalidateQueries({ queryKey: ['schemas', selectedCatalogName] });
      }
    },
  });

  const filteredTables = useMemo(() => {
    if (!tables) return [];
    if (!tableSearch) return tables;
    return tables.filter(table =>
      table.Name.toLowerCase().includes(tableSearch.toLowerCase())
    );
  }, [tables, tableSearch]);

  const handleCatalogChange = (value: string) => {
    setSelectedCatalogName(value);
    setSelectedSchemaName(null);
  };

  const handleSchemaChange = (value: string) => {
    setSelectedSchemaName(value);
  };

  const showAllCatalogs = selectedCatalogName === '__ALL__';
  const showAllSchemas = selectedSchemaName === '__ALL__';
  const showSingleSchema = selectedCatalogName && selectedCatalogName !== '__ALL__' && selectedSchemaName && selectedSchemaName !== '__ALL__';

  return (
    <div className="h-full w-full flex flex-col">
      {/* Top Bar */}
      <div className="border-b bg-background p-4">
        <div className="flex items-center gap-4">
          <Select value={selectedCatalogName || ""} onValueChange={handleCatalogChange}>
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Select catalog..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="__ALL__">All Catalogs</SelectItem>
              {isLoadingCatalogs ? (
                <SelectItem value="_loading" disabled>Loading...</SelectItem>
              ) : (
                catalogs?.map((catalog) => (
                  <SelectItem key={catalog.Name} value={catalog.Name}>
                    {catalog.Name}
                  </SelectItem>
                ))
              )}
            </SelectContent>
          </Select>

          <Select
            value={selectedSchemaName || ""}
            onValueChange={handleSchemaChange}
            disabled={!selectedCatalogName || selectedCatalogName === '__ALL__'}
          >
            <SelectTrigger className="w-[200px]">
              <SelectValue placeholder="Select schema..." />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="__ALL__">All Schemas</SelectItem>
              {isLoadingSchemas ? (
                <SelectItem value="_loading" disabled>Loading...</SelectItem>
              ) : (
                schemas?.map((schema) => (
                  <SelectItem key={schema.Name} value={schema.Name}>
                    {schema.Name}
                  </SelectItem>
                ))
              )}
            </SelectContent>
          </Select>

          <div className="flex-1" />

          <Button
            onClick={() => syncMutation.mutate()}
            disabled={syncMutation.isPending}
            variant="outline"
            size="sm"
          >
            {syncMutation.isPending ? (
              <>
                <Zap className="w-4 h-4 mr-2 animate-spin" />
                Syncing...
              </>
            ) : (
              <>
                <Zap className="w-4 h-4 mr-2" />
                Sync
              </>
            )}
          </Button>
        </div>

        {(selectedSchemaName || selectedCatalogName === '__ALL__') && (
          <div className="mt-3 relative">
            <Search className="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search tables..."
              value={tableSearch}
              onChange={(e) => setTableSearch(e.target.value)}
              className="pl-9"
            />
          </div>
        )}
      </div>

      {/* Tables List */}
      <div className="flex-1 overflow-y-auto p-4">
        {!selectedCatalogName ? (
          <div className="flex items-center justify-center h-full text-muted-foreground">
            <p>Select a catalog to view tables</p>
          </div>
        ) : showAllCatalogs ? (
          <div className="max-w-4xl">
            {catalogs?.map((catalog) => (
              <CatalogSection
                key={catalog.Name}
                catalogName={catalog.Name}
                tableSearch={tableSearch}
              />
            ))}
          </div>
        ) : !selectedSchemaName ? (
          <div className="flex items-center justify-center h-full text-muted-foreground">
            <p>Select a schema to view tables</p>
          </div>
        ) : showAllSchemas ? (
          <div className="max-w-4xl">
            {schemas?.map((schema) => (
              <SchemaSection
                key={schema.Name}
                catalogName={selectedCatalogName}
                schemaName={schema.Name}
                tableSearch={tableSearch}
              />
            ))}
          </div>
        ) : isLoadingTables ? (
          <div className="text-muted-foreground">Loading tables...</div>
        ) : filteredTables.length === 0 ? (
          <div className="text-muted-foreground">
            {tableSearch ? "No tables match your search" : "No tables found"}
          </div>
        ) : (
          <div className="max-w-4xl">
            <div className="text-sm font-semibold text-muted-foreground mb-2 px-3">
              {selectedCatalogName} › {selectedSchemaName}
            </div>
            {filteredTables.map((table) => (
              <TableItem
                key={table.Name}
                catalogName={selectedCatalogName!}
                schemaName={selectedSchemaName!}
                tableName={table.Name}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export function InventoryPage() {
  return (
    <QueryClientProvider client={queryClient}>
      <InventoryPageContent />
    </QueryClientProvider>
  );
}
