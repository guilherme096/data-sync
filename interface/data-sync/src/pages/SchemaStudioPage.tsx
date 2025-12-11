import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useQuery, useMutation, useQueryClient, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { api, type GlobalTable, type GlobalColumn, type TableMapping, type ColumnMapping, type Catalog, type Schema, type Table as LocalTable, type Column as LocalColumn } from '@/lib/api'
import { Plus, Table, Columns3, Trash2, Link, X } from 'lucide-react'
import { Separator } from '@/components/ui/separator'

const queryClient = new QueryClient();

function SchemaStudioPageContent() {
  const queryClientInstance = useQueryClient()

  // Table creation state
  const [isCreatingTable, setIsCreatingTable] = useState(false)
  const [newTableName, setNewTableName] = useState('')
  const [newTableDescription, setNewTableDescription] = useState('')

  // Selected table
  const [selectedTable, setSelectedTable] = useState<string | null>(null)

  // Column creation state
  const [isAddingColumn, setIsAddingColumn] = useState(false)
  const [newColumnName, setNewColumnName] = useState('')
  const [newColumnType, setNewColumnType] = useState('varchar')
  const [newColumnDesc, setNewColumnDesc] = useState('')

  // Table mapping state
  const [isAddingTableMapping, setIsAddingTableMapping] = useState(false)
  const [mappingCatalog, setMappingCatalog] = useState('')
  const [mappingSchema, setMappingSchema] = useState('')
  const [mappingTable, setMappingTable] = useState('')

  // Column mapping state
  const [isAddingColumnMapping, setIsAddingColumnMapping] = useState(false)
  const [selectedGlobalColumn, setSelectedGlobalColumn] = useState<string | null>(null)
  const [colMappingCatalog, setColMappingCatalog] = useState('')
  const [colMappingSchema, setColMappingSchema] = useState('')
  const [colMappingTable, setColMappingTable] = useState('')
  const [colMappingColumn, setColMappingColumn] = useState('')

  // Queries - Global data
  const { data: globalTables, isLoading } = useQuery<GlobalTable[], Error>({
    queryKey: ['globalTables'],
    queryFn: api.listGlobalTables,
  })

  const { data: globalColumns } = useQuery<GlobalColumn[], Error>({
    queryKey: ['globalColumns', selectedTable],
    queryFn: () => api.listGlobalColumns(selectedTable!),
    enabled: !!selectedTable,
  })

  const { data: tableMappings } = useQuery<TableMapping[], Error>({
    queryKey: ['tableMappings', selectedTable],
    queryFn: () => api.listTableMappings(selectedTable!),
    enabled: !!selectedTable,
  })

  // Queries - Local data for selection
  const { data: catalogs } = useQuery<Catalog[], Error>({
    queryKey: ['catalogs'],
    queryFn: api.listCatalogs,
  })

  const { data: schemas } = useQuery<Schema[], Error>({
    queryKey: ['schemas', mappingCatalog || colMappingCatalog],
    queryFn: () => api.listSchemas(mappingCatalog || colMappingCatalog),
    enabled: !!(mappingCatalog || colMappingCatalog),
  })

  const { data: localTables } = useQuery<LocalTable[], Error>({
    queryKey: ['localTables', colMappingCatalog, colMappingSchema],
    queryFn: () => api.discoverTables(colMappingCatalog, colMappingSchema),
    enabled: !!(colMappingCatalog && colMappingSchema),
  })

  const { data: localColumns } = useQuery<LocalColumn[], Error>({
    queryKey: ['localColumns', colMappingCatalog, colMappingSchema, colMappingTable],
    queryFn: () => api.discoverColumns(colMappingCatalog, colMappingSchema, colMappingTable),
    enabled: !!(colMappingCatalog && colMappingSchema && colMappingTable),
  })

  // Mutations - Tables
  const createTableMutation = useMutation({
    mutationFn: (table: { Name: string; Description: string }) => api.createGlobalTable(table),
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
      setIsCreatingTable(false)
      setNewTableName('')
      setNewTableDescription('')
    },
  })

  const deleteTableMutation = useMutation({
    mutationFn: (name: string) => api.deleteGlobalTable(name),
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
      setSelectedTable(null)
    },
  })

  // Mutations - Columns
  const createColumnMutation = useMutation({
    mutationFn: ({ tableName, column }: { tableName: string; column: { Name: string; DataType: string; Description: string } }) =>
      api.createGlobalColumn(tableName, column),
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['globalColumns', selectedTable] })
      setIsAddingColumn(false)
      setNewColumnName('')
      setNewColumnType('varchar')
      setNewColumnDesc('')
    },
  })

  // Mutations - Table Mappings
  const createTableMappingMutation = useMutation({
    mutationFn: ({ tableName, mapping }: { tableName: string; mapping: { CatalogName: string; SchemaName: string; TableName: string } }) =>
      api.createTableMapping(tableName, mapping),
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['tableMappings', selectedTable] })
      setIsAddingTableMapping(false)
      setMappingCatalog('')
      setMappingSchema('')
      setMappingTable('')
    },
  })

  // Mutations - Column Mappings
  const createColumnMappingMutation = useMutation({
    mutationFn: ({ tableName, columnName, mapping }: { tableName: string; columnName: string; mapping: { CatalogName: string; SchemaName: string; TableName: string; ColumnName: string } }) =>
      api.createColumnMapping(tableName, columnName, mapping),
    onSuccess: () => {
      setIsAddingColumnMapping(false)
      setSelectedGlobalColumn(null)
      setColMappingCatalog('')
      setColMappingSchema('')
      setColMappingTable('')
      setColMappingColumn('')
    },
  })

  const handleCreateTable = () => {
    if (newTableName.trim()) {
      createTableMutation.mutate({
        Name: newTableName,
        Description: newTableDescription,
      })
    }
  }

  const handleCreateColumn = () => {
    if (selectedTable && newColumnName.trim()) {
      createColumnMutation.mutate({
        tableName: selectedTable,
        column: {
          Name: newColumnName,
          DataType: newColumnType,
          Description: newColumnDesc,
        },
      })
    }
  }

  const handleCreateTableMapping = () => {
    if (selectedTable && mappingCatalog && mappingSchema && mappingTable) {
      createTableMappingMutation.mutate({
        tableName: selectedTable,
        mapping: {
          CatalogName: mappingCatalog,
          SchemaName: mappingSchema,
          TableName: mappingTable,
        },
      })
    }
  }

  const handleCreateColumnMapping = () => {
    if (selectedTable && selectedGlobalColumn && colMappingCatalog && colMappingSchema && colMappingTable && colMappingColumn) {
      createColumnMappingMutation.mutate({
        tableName: selectedTable,
        columnName: selectedGlobalColumn,
        mapping: {
          CatalogName: colMappingCatalog,
          SchemaName: colMappingSchema,
          TableName: colMappingTable,
          ColumnName: colMappingColumn,
        },
      })
    }
  }

  return (
    <div className="h-full w-full flex flex-col p-6 gap-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Schema Studio</h1>
        <p className="text-muted-foreground mt-1">
          Create global tables and map them to your data sources
        </p>
      </div>

      {/* Content */}
      <div className="flex-1 grid grid-cols-2 gap-6 overflow-hidden">
        {/* Left Panel - Global Tables List */}
        <div className="flex flex-col gap-4 overflow-hidden">
          <Card className="flex-shrink-0">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle className="flex items-center gap-2">
                    <Table className="w-5 h-5" />
                    Global Tables
                  </CardTitle>
                  <CardDescription>
                    Unified views of your distributed data
                  </CardDescription>
                </div>
                <Button
                  size="sm"
                  onClick={() => setIsCreatingTable(true)}
                  disabled={isCreatingTable}
                >
                  <Plus className="w-4 h-4 mr-2" />
                  New
                </Button>
              </div>
            </CardHeader>
            <CardContent className="max-h-[400px] overflow-y-auto">
              {isCreatingTable && (
                <div className="mb-4 p-4 border rounded-lg space-y-3 bg-muted/20">
                  <Input
                    placeholder="Table name (e.g., global_customers)"
                    value={newTableName}
                    onChange={(e) => setNewTableName(e.target.value)}
                  />
                  <Input
                    placeholder="Description"
                    value={newTableDescription}
                    onChange={(e) => setNewTableDescription(e.target.value)}
                  />
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      onClick={handleCreateTable}
                      disabled={!newTableName.trim() || createTableMutation.isPending}
                    >
                      Create
                    </Button>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setIsCreatingTable(false)
                        setNewTableName('')
                        setNewTableDescription('')
                      }}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              )}

              {isLoading ? (
                <div className="text-sm text-muted-foreground">Loading...</div>
              ) : !globalTables || globalTables.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Table className="w-12 h-12 mx-auto mb-2 opacity-20" />
                  <p className="text-sm">No global tables yet</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {globalTables.map((table) => (
                    <div
                      key={table.Name}
                      className={`p-3 border rounded-lg cursor-pointer transition-colors ${
                        selectedTable === table.Name
                          ? 'bg-primary/10 border-primary'
                          : 'hover:bg-muted/50'
                      }`}
                      onClick={() => {
                        setSelectedTable(table.Name)
                        setIsAddingColumn(false)
                        setIsAddingTableMapping(false)
                        setIsAddingColumnMapping(false)
                      }}
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="font-medium">{table.Name}</div>
                          {table.Description && (
                            <div className="text-xs text-muted-foreground mt-1">
                              {table.Description}
                            </div>
                          )}
                        </div>
                        <Button
                          size="sm"
                          variant="ghost"
                          className="h-6 w-6 p-0 text-destructive hover:text-destructive hover:bg-destructive/10"
                          onClick={(e) => {
                            e.stopPropagation()
                            if (confirm(`Delete table "${table.Name}"?`)) {
                              deleteTableMutation.mutate(table.Name)
                            }
                          }}
                        >
                          <Trash2 className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>

        {/* Right Panel - Details */}
        <div className="flex flex-col gap-4 overflow-y-auto">
          {!selectedTable ? (
            <Card>
              <CardContent className="pt-6">
                <div className="text-center py-12 text-muted-foreground">
                  <Columns3 className="w-12 h-12 mx-auto mb-2 opacity-20" />
                  <p className="text-sm">Select a table from the left</p>
                </div>
              </CardContent>
            </Card>
          ) : (
            <>
              {/* Columns Card */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle className="flex items-center gap-2 text-base">
                      <Columns3 className="w-4 h-4" />
                      Columns
                    </CardTitle>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setIsAddingColumn(!isAddingColumn)}
                    >
                      <Plus className="w-3 h-3 mr-2" />
                      Add Column
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  {isAddingColumn && (
                    <div className="mb-4 p-3 border rounded-lg space-y-2 bg-muted/20">
                      <Input
                        placeholder="Column name (e.g., id)"
                        value={newColumnName}
                        onChange={(e) => setNewColumnName(e.target.value)}
                      />
                      <Select value={newColumnType} onValueChange={setNewColumnType}>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                        <SelectContent>
                          <SelectItem value="integer">integer</SelectItem>
                          <SelectItem value="bigint">bigint</SelectItem>
                          <SelectItem value="varchar">varchar</SelectItem>
                          <SelectItem value="text">text</SelectItem>
                          <SelectItem value="decimal">decimal</SelectItem>
                          <SelectItem value="boolean">boolean</SelectItem>
                          <SelectItem value="date">date</SelectItem>
                          <SelectItem value="timestamp">timestamp</SelectItem>
                        </SelectContent>
                      </Select>
                      <Input
                        placeholder="Description"
                        value={newColumnDesc}
                        onChange={(e) => setNewColumnDesc(e.target.value)}
                      />
                      <div className="flex gap-2">
                        <Button size="sm" onClick={handleCreateColumn} disabled={!newColumnName.trim()}>
                          Add
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => setIsAddingColumn(false)}>
                          Cancel
                        </Button>
                      </div>
                    </div>
                  )}

                  {!globalColumns || globalColumns.length === 0 ? (
                    <div className="text-sm text-muted-foreground py-4">No columns yet</div>
                  ) : (
                    <div className="space-y-2">
                      {globalColumns.map((column) => (
                        <div key={column.Name} className="p-2 border rounded flex items-center justify-between">
                          <div className="flex items-baseline gap-2">
                            <span className="font-medium text-sm">{column.Name}</span>
                            <code className="text-xs bg-muted px-1.5 py-0.5 rounded">{column.DataType}</code>
                          </div>
                          <Button
                            size="sm"
                            variant="ghost"
                            className="h-6 text-xs"
                            onClick={() => {
                              setSelectedGlobalColumn(column.Name)
                              setIsAddingColumnMapping(true)
                              setIsAddingTableMapping(false)
                            }}
                          >
                            <Link className="w-3 h-3 mr-1" />
                            Map
                          </Button>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Table Mappings Card */}
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle className="text-base">Table Mappings</CardTitle>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => {
                        setIsAddingTableMapping(!isAddingTableMapping)
                        setIsAddingColumnMapping(false)
                      }}
                    >
                      <Plus className="w-3 h-3 mr-2" />
                      Add Mapping
                    </Button>
                  </div>
                </CardHeader>
                <CardContent>
                  {isAddingTableMapping && (
                    <div className="mb-4 p-3 border rounded-lg space-y-2 bg-muted/20">
                      <Select value={mappingCatalog} onValueChange={setMappingCatalog}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select catalog" />
                        </SelectTrigger>
                        <SelectContent>
                          {catalogs?.map((cat) => (
                            <SelectItem key={cat.Name} value={cat.Name}>{cat.Name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select value={mappingSchema} onValueChange={setMappingSchema} disabled={!mappingCatalog}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select schema" />
                        </SelectTrigger>
                        <SelectContent>
                          {schemas?.map((sch) => (
                            <SelectItem key={sch.Name} value={sch.Name}>{sch.Name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Input
                        placeholder="Table name"
                        value={mappingTable}
                        onChange={(e) => setMappingTable(e.target.value)}
                      />
                      <div className="flex gap-2">
                        <Button size="sm" onClick={handleCreateTableMapping} disabled={!mappingCatalog || !mappingSchema || !mappingTable}>
                          Add
                        </Button>
                        <Button size="sm" variant="outline" onClick={() => setIsAddingTableMapping(false)}>
                          Cancel
                        </Button>
                      </div>
                    </div>
                  )}

                  {!tableMappings || tableMappings.length === 0 ? (
                    <div className="text-sm text-muted-foreground py-4">No mappings yet</div>
                  ) : (
                    <div className="space-y-1">
                      {tableMappings.map((mapping, idx) => (
                        <div key={idx} className="text-sm p-2 bg-muted/30 rounded">
                          {mapping.CatalogName}.{mapping.SchemaName}.{mapping.TableName}
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Column Mapping Card */}
              {isAddingColumnMapping && selectedGlobalColumn && (
                <Card>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-base">Map Column: {selectedGlobalColumn}</CardTitle>
                      <Button
                        size="sm"
                        variant="ghost"
                        onClick={() => {
                          setIsAddingColumnMapping(false)
                          setSelectedGlobalColumn(null)
                        }}
                      >
                        <X className="w-4 h-4" />
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-2">
                      <Select value={colMappingCatalog} onValueChange={setColMappingCatalog}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select catalog" />
                        </SelectTrigger>
                        <SelectContent>
                          {catalogs?.map((cat) => (
                            <SelectItem key={cat.Name} value={cat.Name}>{cat.Name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select value={colMappingSchema} onValueChange={setColMappingSchema} disabled={!colMappingCatalog}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select schema" />
                        </SelectTrigger>
                        <SelectContent>
                          {schemas?.map((sch) => (
                            <SelectItem key={sch.Name} value={sch.Name}>{sch.Name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select value={colMappingTable} onValueChange={setColMappingTable} disabled={!colMappingSchema}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select table" />
                        </SelectTrigger>
                        <SelectContent>
                          {localTables?.map((tbl) => (
                            <SelectItem key={tbl.Name} value={tbl.Name}>{tbl.Name}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Select value={colMappingColumn} onValueChange={setColMappingColumn} disabled={!colMappingTable}>
                        <SelectTrigger>
                          <SelectValue placeholder="Select column" />
                        </SelectTrigger>
                        <SelectContent>
                          {localColumns?.map((col) => (
                            <SelectItem key={col.Name} value={col.Name}>
                              {col.Name} ({col.DataType})
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <Button size="sm" onClick={handleCreateColumnMapping} disabled={!colMappingColumn} className="w-full">
                        Create Mapping
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}

export function SchemaStudioPage() {
  return (
    <QueryClientProvider client={queryClient}>
      <SchemaStudioPageContent />
    </QueryClientProvider>
  )
}
