import { useState } from 'react'
import { ChevronDown, ChevronUp, Database } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'

interface QueryResultTableProps {
  data: {
    generatedSQL?: string
    rows?: Array<Record<string, any>>
    rowCount?: number
    executionTime?: string
  }
}

export function QueryResultTable({ data }: QueryResultTableProps) {
  const [showSQL, setShowSQL] = useState(false)
  const { generatedSQL, rows = [], rowCount = 0, executionTime } = data

  // Get column names from first row
  const columns = rows.length > 0 ? Object.keys(rows[0]) : []

  return (
    <Card className="mt-2 border-blue-200 dark:border-blue-800">
      <CardContent className="pt-4">
        {/* Generated SQL section (collapsible) */}
        {generatedSQL && (
          <div className="mb-4">
            <button
              onClick={() => setShowSQL(!showSQL)}
              className="flex items-center gap-2 text-sm font-medium text-blue-600 dark:text-blue-400 hover:underline"
            >
              <Database className="h-4 w-4" />
              Generated SQL
              {showSQL ? (
                <ChevronUp className="h-4 w-4" />
              ) : (
                <ChevronDown className="h-4 w-4" />
              )}
            </button>
            {showSQL && (
              <pre className="mt-2 rounded-md bg-muted p-3 text-xs overflow-x-auto">
                <code>{generatedSQL}</code>
              </pre>
            )}
          </div>
        )}

        {/* Results table */}
        {rows.length > 0 ? (
          <>
            <div className="rounded-md border overflow-auto max-h-96">
              <Table>
                <TableHeader>
                  <TableRow>
                    {columns.map((col) => (
                      <TableHead key={col} className="font-semibold">
                        {col}
                      </TableHead>
                    ))}
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {rows.map((row, idx) => (
                    <TableRow key={idx}>
                      {columns.map((col) => (
                        <TableCell key={col} className="font-mono text-xs">
                          {formatCellValue(row[col])}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>

            {/* Footer with stats */}
            <div className="mt-3 flex items-center gap-3 text-xs text-muted-foreground">
              <Badge variant="outline" className="font-normal">
                {rowCount} {rowCount === 1 ? 'row' : 'rows'}
              </Badge>
              {executionTime && (
                <Badge variant="outline" className="font-normal">
                  {executionTime}
                </Badge>
              )}
            </div>
          </>
        ) : (
          <div className="text-sm text-muted-foreground py-8 text-center">
            No results to display
          </div>
        )}
      </CardContent>
    </Card>
  )
}

function formatCellValue(value: any): string {
  if (value === null || value === undefined) {
    return 'NULL'
  }
  if (typeof value === 'object') {
    return JSON.stringify(value)
  }
  if (typeof value === 'boolean') {
    return value ? 'true' : 'false'
  }
  return String(value)
}
