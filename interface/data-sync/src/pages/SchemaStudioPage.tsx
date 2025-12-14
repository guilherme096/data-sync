import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { useQuery, useMutation, useQueryClient, QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { api, type Catalog, type Schema, type Table as TableModel, type Column as ColumnModel, type TableRelation, type TableSource, type GlobalTable, type GlobalColumn, type AutoMatchResponse } from '@/lib/api'
import { Plus, Database, GitMerge, Search, Zap, X, Globe, Sparkles, Loader2 } from 'lucide-react'
import { useMemo } from 'react'

const queryClient = new QueryClient();

// ============================================
// DATA SOURCES TAB COMPONENTS (from InventoryPage)
// ============================================

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

// ============================================
// GLOBAL TABLES TAB COMPONENTS
// ============================================

function GlobalTableItem({ tableName }: { tableName: string }) {
  const { data: columns, isLoading } = useQuery<GlobalColumn[], Error>({
    queryKey: ['globalColumns', tableName],
    queryFn: () => api.listGlobalColumns(tableName),
  });

  return (
    <div className="border-b">
      <div className="py-2 px-3 font-medium flex items-center gap-2">
        <Globe className="w-4 h-4 text-primary" />
        {tableName}
      </div>
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
                {column.Description && (
                  <span className="text-xs text-muted-foreground italic">{column.Description}</span>
                )}
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

// ============================================
// STUDIO TAB TYPES AND COMPONENTS
// ============================================

type RelationType = 'JOIN' | 'UNION'

function RelationCard({
  relation,
  relations,
  onDelete,
  onClick
}: {
  relation: TableRelation
  relations: TableRelation[]
  onDelete: () => void
  onClick: () => void
}) {
  const getSourceDisplay = (source: TableSource) => {
    if (source.type === 'physical') {
      return `${source.catalog}.${source.schema}.${source.table}`
    } else {
      const rel = relations.find(r => r.id === source.relationId)
      return rel ? rel.name : 'Unknown Relation'
    }
  }

  const leftDisplay = getSourceDisplay(relation.leftTable)
  const rightDisplay = getSourceDisplay(relation.rightTable)

  return (
    <Card
      className="cursor-pointer hover:bg-muted/50 transition-colors"
      onClick={onClick}
    >
      <CardContent className="">
        {/* Relation Name */}
        <div className="mb-3 flex items-center justify-between">
          <h3 className="text-lg font-semibold">{relation.name}</h3>
          <Button
            size="sm"
            variant="ghost"
            className="text-destructive hover:text-destructive hover:bg-destructive/10"
            onClick={(e) => {
              e.stopPropagation()
              if (confirm(`Delete relation "${relation.name}"?`)) {
                onDelete()
              }
            }}
          >
            <X className="w-4 h-4" />
          </Button>
        </div>

        {/* Relation Details */}
        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 text-sm font-medium px-3 py-1.5 bg-blue-500/10 text-blue-600 dark:text-blue-400 rounded">
            <Database className="w-3.5 h-3.5" />
            {leftDisplay}
          </div>

          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">←</span>
            <div className="px-2 py-1 bg-muted rounded text-xs font-medium">
              {relation.relationType}
              {relation.joinColumn && (
                <span className="ml-1 text-muted-foreground">
                  ({relation.joinColumn.left})
                </span>
              )}
            </div>
            <span className="text-xs text-muted-foreground">→</span>
          </div>

          <div className="flex items-center gap-2 text-sm font-medium px-3 py-1.5 bg-green-500/10 text-green-600 dark:text-green-400 rounded">
            <Database className="w-3.5 h-3.5" />
            {rightDisplay}
          </div>
        </div>

        {relation.description && (
          <div className="mt-2 text-xs text-muted-foreground">
            {relation.description}
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function AddRelationDialog({
  onAdd,
  existingRelations
}: {
  onAdd: (relation: Omit<TableRelation, 'id'>) => void
  existingRelations: TableRelation[]
}) {
  const [open, setOpen] = useState(false)
  const [relationName, setRelationName] = useState('')
  const [relationType, setRelationType] = useState<RelationType>('JOIN')
  const [description, setDescription] = useState('')

  // Source type selection
  const [leftSourceType, setLeftSourceType] = useState<'physical' | 'relation'>('physical')
  const [rightSourceType, setRightSourceType] = useState<'physical' | 'relation'>('physical')

  // For physical table sources
  const [leftCatalog, setLeftCatalog] = useState('')
  const [leftSchema, setLeftSchema] = useState('')
  const [leftTable, setLeftTable] = useState('')
  const [leftColumn, setLeftColumn] = useState('')

  const [rightCatalog, setRightCatalog] = useState('')
  const [rightSchema, setRightSchema] = useState('')
  const [rightTable, setRightTable] = useState('')
  const [rightColumn, setRightColumn] = useState('')

  // For relation sources
  const [leftRelationId, setLeftRelationId] = useState('')
  const [rightRelationId, setRightRelationId] = useState('')

  const { data: catalogs } = useQuery<Catalog[], Error>({
    queryKey: ['catalogs'],
    queryFn: api.listCatalogs,
  })

  const { data: leftSchemas } = useQuery<Schema[], Error>({
    queryKey: ['schemas', leftCatalog],
    queryFn: () => api.listSchemas(leftCatalog),
    enabled: !!leftCatalog,
  })

  const { data: leftTables } = useQuery<TableModel[], Error>({
    queryKey: ['tables', leftCatalog, leftSchema],
    queryFn: () => api.discoverTables(leftCatalog, leftSchema),
    enabled: !!leftCatalog && !!leftSchema,
  })

  const { data: leftColumns } = useQuery<ColumnModel[], Error>({
    queryKey: ['columns', leftCatalog, leftSchema, leftTable],
    queryFn: () => api.discoverColumns(leftCatalog, leftSchema, leftTable),
    enabled: !!leftCatalog && !!leftSchema && !!leftTable,
  })

  const { data: rightSchemas } = useQuery<Schema[], Error>({
    queryKey: ['schemas', rightCatalog],
    queryFn: () => api.listSchemas(rightCatalog),
    enabled: !!rightCatalog,
  })

  const { data: rightTables } = useQuery<TableModel[], Error>({
    queryKey: ['tables', rightCatalog, rightSchema],
    queryFn: () => api.discoverTables(rightCatalog, rightSchema),
    enabled: !!rightCatalog && !!rightSchema,
  })

  const { data: rightColumns } = useQuery<ColumnModel[], Error>({
    queryKey: ['columns', rightCatalog, rightSchema, rightTable],
    queryFn: () => api.discoverColumns(rightCatalog, rightSchema, rightTable),
    enabled: !!rightCatalog && !!rightSchema && !!rightTable,
  })

  const handleAdd = () => {
    const relation: Omit<TableRelation, 'id'> = {
      name: relationName,
      leftTable: leftSourceType === 'physical'
        ? { type: 'physical', catalog: leftCatalog, schema: leftSchema, table: leftTable }
        : { type: 'relation', relationId: leftRelationId },
      rightTable: rightSourceType === 'physical'
        ? { type: 'physical', catalog: rightCatalog, schema: rightSchema, table: rightTable }
        : { type: 'relation', relationId: rightRelationId },
      relationType,
      description: description || undefined,
    }

    if (relationType === 'JOIN' && leftColumn && rightColumn) {
      relation.joinColumn = {
        left: leftColumn,
        right: rightColumn,
      }
    }

    onAdd(relation)
    setOpen(false)
    // Reset form
    setRelationName('')
    setLeftSourceType('physical')
    setRightSourceType('physical')
    setLeftCatalog('')
    setLeftSchema('')
    setLeftTable('')
    setLeftColumn('')
    setRightCatalog('')
    setRightSchema('')
    setRightTable('')
    setRightColumn('')
    setLeftRelationId('')
    setRightRelationId('')
    setDescription('')
  }

  // Validate that join columns have matching data types
  const leftColumnType = leftColumns?.find(col => col.Name === leftColumn)?.DataType
  const rightColumnType = rightColumns?.find(col => col.Name === rightColumn)?.DataType
  const doJoinColumnsMatch = relationType === 'UNION' || !leftColumn || !rightColumn || leftColumnType === rightColumnType

  const isNameValid = relationName.trim() && !existingRelations.some(r => r.name === relationName)
  const isLeftValid = leftSourceType === 'physical'
    ? (leftCatalog && leftSchema && leftTable)
    : leftRelationId
  const isRightValid = rightSourceType === 'physical'
    ? (rightCatalog && rightSchema && rightTable)
    : rightRelationId
  const areColumnsValid = relationType === 'UNION' || (leftColumn && rightColumn)

  const isValid = isNameValid && isLeftValid && isRightValid && areColumnsValid && doJoinColumnsMatch

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Card className="cursor-pointer hover:bg-muted/50 transition-colors border-dashed">
          <CardContent className="">
            <div className="flex items-center justify-center gap-2 text-muted-foreground">
              <Plus className="w-5 h-5" />
              <span className="text-sm font-medium">Add Relation</span>
            </div>
          </CardContent>
        </Card>
      </DialogTrigger>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Create Table Relation</DialogTitle>
          <DialogDescription>
            Create a named relation between tables or other relations
          </DialogDescription>
        </DialogHeader>

        {/* Relation Name */}
        <div className="space-y-2 mt-4">
          <label className="text-sm font-semibold">Relation Name *</label>
          <Input
            placeholder="e.g., combined_users, all_orders..."
            value={relationName}
            onChange={(e) => setRelationName(e.target.value)}
          />
          {relationName && !isNameValid && (
            <p className="text-xs text-destructive">Name already exists or is invalid</p>
          )}
        </div>

        <div className="grid grid-cols-2 gap-6 mt-4">
          {/* Left Source */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Left Source</h3>

            <div>
              <label className="text-xs text-muted-foreground">Source Type</label>
              <Select value={leftSourceType} onValueChange={(v) => setLeftSourceType(v as 'physical' | 'relation')}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="physical">Physical Table</SelectItem>
                  <SelectItem value="relation">Existing Relation</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {leftSourceType === 'physical' ? (
              <>
                <div>
                  <label className="text-xs text-muted-foreground">Catalog</label>
                  <Select value={leftCatalog} onValueChange={setLeftCatalog}>
                    <SelectTrigger><SelectValue placeholder="Select catalog" /></SelectTrigger>
                    <SelectContent>
                      {catalogs?.map(cat => (
                        <SelectItem key={cat.Name} value={cat.Name}>{cat.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <label className="text-xs text-muted-foreground">Schema</label>
                  <Select value={leftSchema} onValueChange={setLeftSchema} disabled={!leftCatalog}>
                    <SelectTrigger><SelectValue placeholder="Select schema" /></SelectTrigger>
                    <SelectContent>
                      {leftSchemas?.map(sch => (
                        <SelectItem key={sch.Name} value={sch.Name}>{sch.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <label className="text-xs text-muted-foreground">Table</label>
                  <Select value={leftTable} onValueChange={setLeftTable} disabled={!leftSchema}>
                    <SelectTrigger><SelectValue placeholder="Select table" /></SelectTrigger>
                    <SelectContent>
                      {leftTables?.map(tbl => (
                        <SelectItem key={tbl.Name} value={tbl.Name}>{tbl.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </>
            ) : (
              <div>
                <label className="text-xs text-muted-foreground">Relation</label>
                <Select value={leftRelationId} onValueChange={setLeftRelationId}>
                  <SelectTrigger><SelectValue placeholder="Select relation" /></SelectTrigger>
                  <SelectContent>
                    {existingRelations.map(rel => (
                      <SelectItem key={rel.id} value={rel.id}>{rel.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
          </div>

          {/* Right Source */}
          <div className="space-y-3">
            <h3 className="text-sm font-semibold">Right Source</h3>

            <div>
              <label className="text-xs text-muted-foreground">Source Type</label>
              <Select value={rightSourceType} onValueChange={(v) => setRightSourceType(v as 'physical' | 'relation')}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="physical">Physical Table</SelectItem>
                  <SelectItem value="relation">Existing Relation</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {rightSourceType === 'physical' ? (
              <>
                <div>
                  <label className="text-xs text-muted-foreground">Catalog</label>
                  <Select value={rightCatalog} onValueChange={setRightCatalog}>
                    <SelectTrigger><SelectValue placeholder="Select catalog" /></SelectTrigger>
                    <SelectContent>
                      {catalogs?.map(cat => (
                        <SelectItem key={cat.Name} value={cat.Name}>{cat.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <label className="text-xs text-muted-foreground">Schema</label>
                  <Select value={rightSchema} onValueChange={setRightSchema} disabled={!rightCatalog}>
                    <SelectTrigger><SelectValue placeholder="Select schema" /></SelectTrigger>
                    <SelectContent>
                      {rightSchemas?.map(sch => (
                        <SelectItem key={sch.Name} value={sch.Name}>{sch.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <label className="text-xs text-muted-foreground">Table</label>
                  <Select value={rightTable} onValueChange={setRightTable} disabled={!rightSchema}>
                    <SelectTrigger><SelectValue placeholder="Select table" /></SelectTrigger>
                    <SelectContent>
                      {rightTables?.map(tbl => (
                        <SelectItem key={tbl.Name} value={tbl.Name}>{tbl.Name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              </>
            ) : (
              <div>
                <label className="text-xs text-muted-foreground">Relation</label>
                <Select value={rightRelationId} onValueChange={setRightRelationId}>
                  <SelectTrigger><SelectValue placeholder="Select relation" /></SelectTrigger>
                  <SelectContent>
                    {existingRelations.map(rel => (
                      <SelectItem key={rel.id} value={rel.id}>{rel.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            )}
          </div>
        </div>

        {/* Relation Type */}
        <div className="space-y-2 mt-4">
          <label className="text-sm font-semibold">Relation Type</label>
          <Select value={relationType} onValueChange={(v) => setRelationType(v as RelationType)}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="JOIN">JOIN</SelectItem>
              <SelectItem value="UNION">UNION</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Join Columns - appears only when JOIN is selected */}
        {relationType === 'JOIN' && (
          <div className="space-y-3 mt-4">
            <div className="grid grid-cols-2 gap-6 p-4 border rounded-lg bg-muted/30">
              <div className="space-y-2">
                <label className="text-sm font-semibold">Left Join Column</label>
                <Select
                  value={leftColumn}
                  onValueChange={setLeftColumn}
                  disabled={leftSourceType === 'physical' ? !leftTable : !leftRelationId}
                >
                  <SelectTrigger><SelectValue placeholder="Select column" /></SelectTrigger>
                  <SelectContent>
                    {leftSourceType === 'physical' && leftColumns?.map(col => (
                      <SelectItem key={col.Name} value={col.Name}>
                        {col.Name} ({col.DataType})
                      </SelectItem>
                    ))}
                    {leftSourceType === 'physical' && (!leftColumns || leftColumns.length === 0) && (
                      <SelectItem value="_no_columns" disabled>No columns available</SelectItem>
                    )}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold">Right Join Column</label>
                <Select
                  value={rightColumn}
                  onValueChange={setRightColumn}
                  disabled={rightSourceType === 'physical' ? !rightTable : !rightRelationId}
                >
                  <SelectTrigger><SelectValue placeholder="Select column" /></SelectTrigger>
                  <SelectContent>
                    {rightSourceType === 'physical' && rightColumns?.map(col => (
                      <SelectItem key={col.Name} value={col.Name}>
                        {col.Name} ({col.DataType})
                      </SelectItem>
                    ))}
                    {rightSourceType === 'physical' && (!rightColumns || rightColumns.length === 0) && (
                      <SelectItem value="_no_columns" disabled>No columns available</SelectItem>
                    )}
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Data type mismatch warning */}
            {leftColumn && rightColumn && leftSourceType === 'physical' && rightSourceType === 'physical' && !doJoinColumnsMatch && (
              <div className="flex items-start gap-2 p-3 bg-destructive/10 border border-destructive/20 rounded-md">
                <span className="text-destructive text-sm font-medium">⚠</span>
                <div className="flex-1 text-sm">
                  <p className="font-semibold text-destructive">Data type mismatch</p>
                  <p className="text-destructive/90 mt-1">
                    Cannot join columns with different data types: <code className="bg-destructive/20 px-1 rounded">{leftColumnType}</code> and <code className="bg-destructive/20 px-1 rounded">{rightColumnType}</code>
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Description */}
        <div className="space-y-2">
          <label className="text-sm font-semibold">Description (optional)</label>
          <Input
            placeholder="Describe this relation..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </div>

        <div className="flex justify-end gap-2 mt-4">
          <Button variant="outline" onClick={() => setOpen(false)}>Cancel</Button>
          <Button onClick={handleAdd} disabled={!isValid}>Create Relation</Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}

function RelationDetailsSidebar({
  relation,
  open,
  onOpenChange
}: {
  relation: TableRelation | null
  open: boolean
  onOpenChange: (open: boolean) => void
}) {
  if (!relation) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle>Relation Details</DialogTitle>
          <DialogDescription>
            View detailed information about this table relation
          </DialogDescription>
        </DialogHeader>

        <div className="mt-4 space-y-6">
          <div>
            <h3 className="text-sm font-semibold mb-2">Relation Type</h3>
            <div className="px-3 py-2 bg-muted rounded text-sm font-mono">
              {relation.relationType}
            </div>
          </div>

          <div>
            <h3 className="text-sm font-semibold mb-2">Left Table</h3>
            <div className="space-y-1 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Catalog:</span>
                <span className="font-mono">{relation.leftTable.catalog}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Schema:</span>
                <span className="font-mono">{relation.leftTable.schema}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Table:</span>
                <span className="font-mono">{relation.leftTable.table}</span>
              </div>
              {relation.joinColumn && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Join Column:</span>
                  <span className="font-mono">{relation.joinColumn.left}</span>
                </div>
              )}
            </div>
          </div>

          <div>
            <h3 className="text-sm font-semibold mb-2">Right Table</h3>
            <div className="space-y-1 text-sm">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Catalog:</span>
                <span className="font-mono">{relation.rightTable.catalog}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Schema:</span>
                <span className="font-mono">{relation.rightTable.schema}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Table:</span>
                <span className="font-mono">{relation.rightTable.table}</span>
              </div>
              {relation.joinColumn && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Join Column:</span>
                  <span className="font-mono">{relation.joinColumn.right}</span>
                </div>
              )}
            </div>
          </div>

          {relation.description && (
            <div>
              <h3 className="text-sm font-semibold mb-2">Description</h3>
              <p className="text-sm text-muted-foreground">
                {relation.description}
              </p>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}

// ============================================
// MAIN COMPONENT
// ============================================

function SchemaStudioPageContent() {
  const queryClientInstance = useQueryClient()

  // Data Sources tab state
  const [selectedCatalogName, setSelectedCatalogName] = useState<string | null>(null)
  const [selectedSchemaName, setSelectedSchemaName] = useState<string | null>(null)
  const [tableSearch, setTableSearch] = useState('')

  // Studio tab state
  const [selectedRelation, setSelectedRelation] = useState<TableRelation | null>(null)
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [matchResults, setMatchResults] = useState<AutoMatchResponse | null>(null)

  // Load relations from API
  const { data: relations = [], isLoading: isLoadingRelations } = useQuery<TableRelation[], Error>({
    queryKey: ['tableRelations'],
    queryFn: api.listTableRelations,
  })

  // Load global tables
  const { data: globalTables = [], isLoading: isLoadingGlobalTables } = useQuery<GlobalTable[], Error>({
    queryKey: ['globalTables'],
    queryFn: api.listGlobalTables,
  })

  // Queries
  const { data: catalogs, isLoading: isLoadingCatalogs } = useQuery<Catalog[], Error>({
    queryKey: ['catalogs'],
    queryFn: api.listCatalogs,
  })

  const { data: schemas, isLoading: isLoadingSchemas } = useQuery<Schema[], Error>({
    queryKey: ['schemas', selectedCatalogName],
    queryFn: () => api.listSchemas(selectedCatalogName!),
    enabled: !!selectedCatalogName && selectedCatalogName !== '__ALL__',
  })

  const { data: tables, isLoading: isLoadingTables } = useQuery<TableModel[], Error>({
    queryKey: ['tables', selectedCatalogName, selectedSchemaName],
    queryFn: () => api.discoverTables(selectedCatalogName!, selectedSchemaName!),
    enabled: !!selectedCatalogName && selectedCatalogName !== '__ALL__' && !!selectedSchemaName && selectedSchemaName !== '__ALL__',
  })

  const syncMutation = useMutation({
    mutationFn: api.syncMetadata,
    onSuccess: async () => {
      // Invalidate and refetch catalogs
      await queryClientInstance.invalidateQueries({ queryKey: ['catalogs'] })
      await queryClientInstance.refetchQueries({ queryKey: ['catalogs'] })

      // Invalidate and refetch global tables
      await queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
      await queryClientInstance.refetchQueries({ queryKey: ['globalTables'] })

      // Invalidate and refetch schemas if a catalog is selected
      if (selectedCatalogName && selectedCatalogName !== '__ALL__') {
        await queryClientInstance.invalidateQueries({ queryKey: ['schemas', selectedCatalogName] })
        await queryClientInstance.refetchQueries({ queryKey: ['schemas', selectedCatalogName] })
      }

      // Invalidate and refetch tables if both catalog and schema are selected
      if (selectedCatalogName && selectedCatalogName !== '__ALL__' && selectedSchemaName && selectedSchemaName !== '__ALL__') {
        await queryClientInstance.invalidateQueries({ queryKey: ['tables', selectedCatalogName, selectedSchemaName] })
        await queryClientInstance.refetchQueries({ queryKey: ['tables', selectedCatalogName, selectedSchemaName] })
      }
    },
  })

  const filteredTables = useMemo(() => {
    if (!tables) return []
    if (!tableSearch) return tables
    return tables.filter(table =>
      table.Name.toLowerCase().includes(tableSearch.toLowerCase())
    )
  }, [tables, tableSearch])

  const handleCatalogChange = (value: string) => {
    setSelectedCatalogName(value)
    setSelectedSchemaName(null)
  }

  const handleSchemaChange = (value: string) => {
    setSelectedSchemaName(value)
  }

  const showAllCatalogs = selectedCatalogName === '__ALL__'
  const showAllSchemas = selectedSchemaName === '__ALL__'

  const createRelationMutation = useMutation({
    mutationFn: (relation: Omit<TableRelation, 'id'>) => {
      const newRelation: TableRelation = {
        ...relation,
        id: Date.now().toString(),
      }
      return api.createTableRelation(newRelation)
    },
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['tableRelations'] })
      queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
    },
  })

  const deleteRelationMutation = useMutation({
    mutationFn: api.deleteTableRelation,
    onSuccess: () => {
      queryClientInstance.invalidateQueries({ queryKey: ['tableRelations'] })
      queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
    },
  })

  const autoMatchMutation = useMutation({
    mutationFn: () => api.autoMatchRelations({ maxSuggestions: 10, autoCreate: true }),
    onSuccess: (data) => {
      setMatchResults(data)
      queryClientInstance.invalidateQueries({ queryKey: ['tableRelations'] })
      queryClientInstance.invalidateQueries({ queryKey: ['globalTables'] })
    },
  })

  const handleAddRelation = (relation: Omit<TableRelation, 'id'>) => {
    createRelationMutation.mutate(relation)
  }

  const handleDeleteRelation = (id: string) => {
    deleteRelationMutation.mutate(id)
  }

  const handleRelationClick = (relation: TableRelation) => {
    setSelectedRelation(relation)
    setSidebarOpen(true)
  }

  return (
    <div className="h-full w-full flex flex-col p-6 gap-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Schema Studio</h1>
        <p className="text-muted-foreground mt-1">
          Manage data sources and create table relations
        </p>
      </div>

      {/* Tabs */}
      <Tabs defaultValue="data-sources" className="flex-1 flex flex-col overflow-hidden">
        <TabsList className="w-fit">
          <TabsTrigger value="data-sources">
            <Database className="w-4 h-4 mr-2" />
            Data Sources
          </TabsTrigger>
          <TabsTrigger value="global-tables">
            <Globe className="w-4 h-4 mr-2" />
            Global Tables
          </TabsTrigger>
          <TabsTrigger value="studio">
            <GitMerge className="w-4 h-4 mr-2" />
            Studio
          </TabsTrigger>
        </TabsList>

        {/* Data Sources Tab */}
        <TabsContent value="data-sources" className="flex-1 overflow-hidden mt-4 flex flex-col">
          {/* Top Bar */}
          <div className="border rounded-lg bg-background p-4 mb-4">
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
          <div className="flex-1 overflow-y-auto">
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
        </TabsContent>

        {/* Global Tables Tab */}
        <TabsContent value="global-tables" className="flex-1 overflow-hidden mt-4 flex flex-col">
          <div className="flex-1 overflow-y-auto">
            {isLoadingGlobalTables ? (
              <div className="text-muted-foreground">Loading global tables...</div>
            ) : globalTables.length === 0 ? (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                <p>No global tables found. Create table relations in the Studio tab to generate global tables.</p>
              </div>
            ) : (
              <div className="max-w-4xl">
                <div className="text-sm font-semibold text-muted-foreground mb-2 px-3">
                  Global Tables ({globalTables.length})
                </div>
                {globalTables.map((table) => (
                  <GlobalTableItem
                    key={table.Name}
                    tableName={table.Name}
                  />
                ))}
              </div>
            )}
          </div>
        </TabsContent>

        {/* Studio Tab */}
        <TabsContent value="studio" className="flex-1 overflow-hidden mt-4">
          <div className="h-full overflow-y-auto">
            <div className="max-w-5xl mx-auto space-y-4">
              {/* Magic Button Header */}
              <div className="mb-4 flex items-center justify-between">
                <div>
                  <h2 className="text-xl font-semibold">Table Relations</h2>
                  <p className="text-sm text-muted-foreground">
                    Manually create or auto-discover table relations
                  </p>
                </div>
                <Button
                  onClick={() => autoMatchMutation.mutate()}
                  disabled={autoMatchMutation.isPending}
                  className="gap-2"
                  variant="default"
                >
                  {autoMatchMutation.isPending ? (
                    <>
                      <Loader2 className="w-4 h-4 animate-spin" />
                      Analyzing...
                    </>
                  ) : (
                    <>
                      <Sparkles className="w-4 h-4" />
                      Auto-Match Relations
                    </>
                  )}
                </Button>
              </div>

              {/* Show results if available */}
              {matchResults && (
                <div className="mb-4 p-4 border rounded-lg bg-muted/30">
                  <h3 className="font-semibold mb-2">Auto-Match Results</h3>
                  <div className="space-y-2 text-sm">
                    <div className="flex items-center justify-between">
                      <span>Suggestions found:</span>
                      <span className="font-mono">{matchResults.suggestions.length}</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <span>Relations created:</span>
                      <span className="font-mono text-green-600">
                        {matchResults.createdRelations?.length || 0}
                      </span>
                    </div>
                    {matchResults.errors && matchResults.errors.length > 0 && (
                      <div className="mt-2">
                        <span className="text-destructive">Errors:</span>
                        <ul className="list-disc list-inside text-destructive text-xs mt-1">
                          {matchResults.errors.map((err, i) => (
                            <li key={i}>{err}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setMatchResults(null)}
                    className="mt-2"
                  >
                    Dismiss
                  </Button>
                </div>
              )}

              {relations.map((relation) => (
                <RelationCard
                  key={relation.id}
                  relation={relation}
                  relations={relations}
                  onDelete={() => handleDeleteRelation(relation.id)}
                  onClick={() => handleRelationClick(relation)}
                />
              ))}

              <AddRelationDialog
                onAdd={handleAddRelation}
                existingRelations={relations}
              />
            </div>
          </div>
        </TabsContent>
      </Tabs>

      {/* Relation Details Sidebar */}
      <RelationDetailsSidebar
        relation={selectedRelation}
        open={sidebarOpen}
        onOpenChange={setSidebarOpen}
      />
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
