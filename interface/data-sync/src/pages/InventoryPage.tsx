import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import { Separator } from "@/components/ui/separator"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { useQuery, useMutation, QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { api, type Catalog, type Schema } from "@/lib/api"
import { useState } from "react"
import { Zap, Database } from "lucide-react"

// Create a client
const queryClient = new QueryClient();

function InventoryPageContent() {
  const [selectedCatalogName, setSelectedCatalogName] = useState<string | null>(null);

  // Fetch catalogs
  const { data: catalogs, isLoading: isLoadingCatalogs, isError: isErrorCatalogs, error: catalogsError } = useQuery<Catalog[], Error>({
    queryKey: ['catalogs'],
    queryFn: api.listCatalogs,
  });

  // Fetch selected catalog details
  const { data: selectedCatalog, isLoading: isLoadingSelectedCatalog, isError: isErrorSelectedCatalog, error: selectedCatalogError } = useQuery<Catalog, Error>({
    queryKey: ['catalog', selectedCatalogName],
    queryFn: () => api.getCatalog(selectedCatalogName!),
    enabled: !!selectedCatalogName, // Only run if a catalog is selected
  });

  // Fetch schemas for selected catalog
  const { data: schemas, isLoading: isLoadingSchemas, isError: isErrorSchemas, error: schemasError } = useQuery<Schema[], Error>({
    queryKey: ['schemas', selectedCatalogName],
    queryFn: () => api.listSchemas(selectedCatalogName!),
    enabled: !!selectedCatalogName, // Only run if a catalog is selected
  });

  // Sync mutation
  const syncMutation = useMutation({
    mutationFn: api.syncMetadata,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['catalogs'] });
      // Invalidate schemas and selected catalog details to refresh if a catalog is selected
      if (selectedCatalogName) {
        queryClient.invalidateQueries({ queryKey: ['catalog', selectedCatalogName] });
        queryClient.invalidateQueries({ queryKey: ['schemas', selectedCatalogName] });
      }
      alert('Metadata sync completed successfully!'); // TODO: Replace with a proper toast notification
    },
    onError: (error) => {
      alert(`Metadata sync failed: ${error.message}`); // TODO: Replace with a proper toast notification
    },
  });

  const handleSync = () => {
    syncMutation.mutate();
  };

  return (
    <div className="h-full w-full bg-muted/20">
      <ResizablePanelGroup direction="horizontal" className="h-full w-full rounded-lg border">
        {/* Left Panel: Catalog List */}
        <ResizablePanel defaultSize={25} minSize={20} maxSize={40}>
          <div className="h-full p-4 flex flex-col gap-4">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold tracking-tight flex items-center gap-2">
                <Database className="w-5 h-5 text-primary" />
                Catalogs
              </h2>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={handleSync} 
                disabled={syncMutation.isPending}
                className="flex items-center gap-1"
              >
                {syncMutation.isPending ? (
                  <>
                    <Zap className="w-3 h-3 animate-spin" />
                    Syncing...
                  </>
                ) : (
                  <>
                    <Zap className="w-3 h-3" />
                    Sync Now
                  </>
                )}
              </Button>
            </div>
            <Separator />
            <div className="flex-1 overflow-y-auto space-y-2">
              {isLoadingCatalogs && <div className="text-sm text-muted-foreground">Loading catalogs...</div>}
              {isErrorCatalogs && <div className="text-sm text-destructive">Error loading catalogs: {catalogsError?.message}</div>}
              {catalogs?.length === 0 && !isLoadingCatalogs && <div className="text-sm text-muted-foreground">No catalogs found.</div>}
              {catalogs?.map((catalog) => (
                <Card 
                  key={catalog.Name} 
                  className={`cursor-pointer hover:bg-muted transition-colors ${selectedCatalogName === catalog.Name ? 'bg-muted border-primary' : ''}`}
                  onClick={() => setSelectedCatalogName(catalog.Name)}
                >
                  <CardContent className="p-3 flex items-center justify-between">
                    <span className="font-medium">{catalog.Name}</span>
                    {/* Could add a status indicator here later */}
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </ResizablePanel>

        <ResizableHandle withHandle />

        {/* Right Panel: Catalog Details and Schemas */}
        <ResizablePanel defaultSize={75}>
          <div className="h-full p-4 flex flex-col gap-4">
            <h2 className="text-lg font-semibold tracking-tight">
              {selectedCatalogName ? `Catalog: ${selectedCatalogName}` : 'Select a Catalog'}
            </h2>
            <Separator />

            {!selectedCatalogName ? (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                <p>Select a catalog from the left panel to view its details.</p>
              </div>
            ) : (
              <div className="flex-1 overflow-y-auto space-y-6">
                {/* Catalog Metadata */}
                <Card>
                  <CardHeader>
                    <CardTitle className="text-md">Metadata</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {isLoadingSelectedCatalog && <div className="text-sm text-muted-foreground">Loading metadata...</div>}
                    {isErrorSelectedCatalog && <div className="text-sm text-destructive">Error loading metadata: {selectedCatalogError?.message}</div>}
                    {selectedCatalog && Object.keys(selectedCatalog.Metadata).length > 0 ? (
                      <dl className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                        {Object.entries(selectedCatalog.Metadata).map(([key, value]) => (
                          <div key={key} className="col-span-1">
                            <dt className="font-medium text-muted-foreground">{key}:</dt>
                            <dd className="break-words">{value}</dd>
                          </div>
                        ))}
                      </dl>
                    ) : (
                      !isLoadingSelectedCatalog && !isErrorSelectedCatalog && <div className="text-sm text-muted-foreground">No metadata available.</div>
                    )}
                  </CardContent>
                </Card>

                {/* Schemas */}
                <Card>
                  <CardHeader>
                    <CardTitle className="text-md">Schemas</CardTitle>
                  </CardHeader>
                  <CardContent>
                    {isLoadingSchemas && <div className="text-sm text-muted-foreground">Loading schemas...</div>}
                    {isErrorSchemas && <div className="text-sm text-destructive">Error loading schemas: {schemasError?.message}</div>}
                    {schemas && schemas.length > 0 ? (
                      <Table>
                        <TableHeader>
                          <TableRow>
                            <TableHead>Name</TableHead>
                            <TableHead>Metadata</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {schemas.map((schema) => (
                            <TableRow key={schema.Name}>
                              <TableCell className="font-medium">{schema.Name}</TableCell>
                              <TableCell>
                                {Object.keys(schema.Metadata).length > 0 ? (
                                  <ul className="list-disc list-inside text-xs text-muted-foreground">
                                    {Object.entries(schema.Metadata).map(([key, value]) => (
                                      <li key={`${schema.Name}-${key}`}>{key}: {value}</li>
                                    ))}
                                  </ul>
                                ) : (
                                  <span className="text-muted-foreground text-xs">No metadata</span>
                                )}
                              </TableCell>
                            </TableRow>
                          ))}
                        </TableBody>
                      </Table>
                    ) : (
                      !isLoadingSchemas && !isErrorSchemas && <div className="text-sm text-muted-foreground">No schemas found for this catalog.</div>
                    )}
                  </CardContent>
                </Card>
              </div>
            )}
          </div>
        </ResizablePanel>
      </ResizablePanelGroup>
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
