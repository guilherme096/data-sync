import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Play, Eraser, Send, Bot, User, Sparkles, Database, Globe } from 'lucide-react'
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import CodeMirror from '@uiw/react-codemirror'
import { sql as sqlLang } from '@codemirror/lang-sql'
import { githubLight, githubDark } from '@uiw/codemirror-theme-github'
import { useTheme } from '@/components/theme-provider'

type QueryMode = 'direct' | 'global'

interface QueryResult {
  generatedSQL?: string
  rows: Array<Record<string, any>>
  rowCount: number
  executionTime?: string
}

const STORAGE_KEY = 'data-sync-sql-query'
const DEFAULT_QUERY = '-- Write your SQL query here\nSELECT * FROM global_users LIMIT 10;'

export function QueryPage() {
  const [sqlCode, setSqlCode] = useState(() => {
    // Load from localStorage on mount
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      return saved || DEFAULT_QUERY
    } catch {
      return DEFAULT_QUERY
    }
  })
  const [chatInput, setChatInput] = useState('')
  const [queryMode, setQueryMode] = useState<QueryMode>('global')
  const [queryResult, setQueryResult] = useState<QueryResult | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const { theme } = useTheme()
  const [isDark, setIsDark] = useState(false)

  // Save SQL code to localStorage whenever it changes
  useEffect(() => {
    try {
      localStorage.setItem(STORAGE_KEY, sqlCode)
    } catch {
      // Ignore localStorage errors
    }
  }, [sqlCode])

  // Determine if we should use dark theme
  useEffect(() => {
    const updateTheme = () => {
      if (theme === 'dark') {
        setIsDark(true)
      } else if (theme === 'light') {
        setIsDark(false)
      } else {
        // system theme
        setIsDark(window.matchMedia('(prefers-color-scheme: dark)').matches)
      }
    }

    updateTheme()

    // Listen for system preference changes when in system mode
    if (theme === 'system') {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
      const handleChange = (e: MediaQueryListEvent) => {
        setIsDark(e.matches)
      }
      mediaQuery.addEventListener('change', handleChange)
      return () => mediaQuery.removeEventListener('change', handleChange)
    }
  }, [theme])

  const executeQuery = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const endpoint = queryMode === 'global' ? '/query/global' : '/query'
      const response = await fetch(`http://localhost:8081${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ query: sqlCode }),
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(errorText || 'Query execution failed')
      }

      const data = await response.json()

      if (queryMode === 'global') {
        setQueryResult({
          generatedSQL: data.generatedSQL,
          rows: data.rows || [],
          rowCount: data.rowCount || 0,
          executionTime: data.executionTime,
        })
      } else {
        setQueryResult({
          rows: data.rows || [],
          rowCount: data.rows?.length || 0,
        })
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error occurred')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="h-full w-full bg-muted/20">
      <ResizablePanelGroup direction="horizontal" className="h-full w-full rounded-lg border">
        
        {/* LEFT PANEL: SQL Workstation */}
        <ResizablePanel defaultSize={70} minSize={30}>
          <div className="h-full p-4 flex flex-col gap-4 min-w-0">
            {/* Toolbar */}
            <div className="flex items-center justify-between bg-muted/50 p-2 rounded-t-md border-b">
               <div className="flex items-center gap-3 px-2">
                  <h2 className="text-sm font-medium tracking-tight text-muted-foreground">SQL Editor</h2>
                  <Select value={queryMode} onValueChange={(value) => setQueryMode(value as QueryMode)}>
                    <SelectTrigger className="h-7 w-[160px] text-xs">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="global">
                        <div className="flex items-center gap-2">
                          <Globe className="w-3.5 h-3.5" />
                          <span>Global Query</span>
                        </div>
                      </SelectItem>
                      <SelectItem value="direct">
                        <div className="flex items-center gap-2">
                          <Database className="w-3.5 h-3.5" />
                          <span>Direct SQL</span>
                        </div>
                      </SelectItem>
                    </SelectContent>
                  </Select>
                  {queryMode === 'global' && (
                    <Badge variant="secondary" className="h-5 text-[10px] px-2">
                      Uses global tables
                    </Badge>
                  )}
               </div>
               <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-7 text-muted-foreground hover:text-foreground"
                    onClick={() => {
                      setSqlCode(DEFAULT_QUERY)
                      setQueryResult(null)
                      setError(null)
                    }}
                  >
                     <Eraser className="w-3.5 h-3.5 mr-2" />
                     Clear
                  </Button>
                  <Button
                    size="sm"
                    className="h-7"
                    onClick={executeQuery}
                    disabled={isLoading || !sqlCode.trim()}
                  >
                     <Play className="w-3.5 h-3.5 mr-2" />
                     {isLoading ? 'Running...' : 'Run'}
                  </Button>
               </div>
            </div>

            <Card className="flex-1 flex flex-col overflow-hidden shadow-sm rounded-t-none mt-[-1rem] border-t-0">
               <div className="flex-1 overflow-auto bg-card">
                  <CodeMirror
                    value={sqlCode}
                    height="100%"
                    extensions={[sqlLang()]}
                    onChange={(val) => setSqlCode(val)}
                    theme={isDark ? githubDark : githubLight}
                    className="h-full text-[13px]"
                    basicSetup={{
                        lineNumbers: true,
                        highlightActiveLineGutter: true,
                        history: true,
                        indentOnInput: true,
                        bracketMatching: true,
                        closeBrackets: true,
                        autocompletion: true,
                        highlightActiveLine: true,
                    }}
                  />
               </div>
               
               <Separator />

               {/* Results Area */}
               <div className="h-1/3 bg-muted/10 p-0 overflow-hidden flex flex-col">
                  <div className="px-4 py-2 bg-muted/30 border-b flex justify-between items-center">
                      <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">Results</span>
                      <div className="flex items-center gap-3">
                        {queryResult?.executionTime && (
                          <span className="text-xs text-muted-foreground">{queryResult.executionTime}</span>
                        )}
                        <span className="text-xs text-muted-foreground">
                          {queryResult ? `${queryResult.rowCount} rows` : '0 rows'}
                        </span>
                      </div>
                  </div>

                  {/* Generated SQL Banner (for global queries) */}
                  {queryResult?.generatedSQL && (
                    <div className="bg-blue-50 dark:bg-blue-950/30 border-b border-blue-200 dark:border-blue-800 px-4 py-2">
                      <div className="flex items-start gap-2">
                        <Globe className="w-3.5 h-3.5 text-blue-600 dark:text-blue-400 mt-0.5" />
                        <div className="flex-1 min-w-0">
                          <p className="text-xs font-medium text-blue-900 dark:text-blue-100 mb-1">Generated Trino SQL:</p>
                          <code className="text-xs text-blue-800 dark:text-blue-200 font-mono block overflow-x-auto">
                            {queryResult.generatedSQL}
                          </code>
                        </div>
                      </div>
                    </div>
                  )}

                  <div className="flex-1 p-4 overflow-auto">
                      {error ? (
                        <div className="flex items-center justify-center h-full">
                          <div className="bg-destructive/10 text-destructive px-4 py-3 rounded-md border border-destructive/20 max-w-2xl">
                            <p className="text-sm font-medium mb-1">Query Error</p>
                            <p className="text-xs font-mono">{error}</p>
                          </div>
                        </div>
                      ) : !queryResult ? (
                        <div className="flex items-center justify-center h-full text-sm text-muted-foreground font-mono">
                          No query executed yet.
                        </div>
                      ) : queryResult.rows.length === 0 ? (
                        <div className="flex items-center justify-center h-full text-sm text-muted-foreground">
                          Query returned no results.
                        </div>
                      ) : (
                        <Table>
                          <TableHeader>
                            <TableRow>
                              {Object.keys(queryResult.rows[0]).map((column) => (
                                <TableHead key={column} className="font-mono text-xs">
                                  {column}
                                </TableHead>
                              ))}
                            </TableRow>
                          </TableHeader>
                          <TableBody>
                            {queryResult.rows.map((row, idx) => (
                              <TableRow key={idx}>
                                {Object.values(row).map((value, colIdx) => (
                                  <TableCell key={colIdx} className="font-mono text-xs">
                                    {value === null ? (
                                      <span className="text-muted-foreground italic">null</span>
                                    ) : typeof value === 'object' ? (
                                      JSON.stringify(value)
                                    ) : (
                                      String(value)
                                    )}
                                  </TableCell>
                                ))}
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      )}
                  </div>
               </div>
            </Card>
          </div>
        </ResizablePanel>

        <ResizableHandle withHandle />

        {/* RIGHT PANEL: AI Agent */}
        <ResizablePanel defaultSize={30} minSize={20}>
          <div className="h-full p-4 flex flex-col gap-4">
            <div className="flex items-center gap-2 h-9">
               <Sparkles className="w-4 h-4 text-primary" />
               <h2 className="text-lg font-semibold tracking-tight">Data Assistant</h2>
            </div>

            <Card className="flex-1 flex flex-col overflow-hidden shadow-sm border-muted-foreground/20">
               {/* Chat History */}
               <div className="flex-1 overflow-y-auto p-4 space-y-4 bg-muted/5">
                  
                  {/* Bot Welcome Message */}
                  <div className="flex gap-3">
                     <Avatar className="h-8 w-8 border">
                        <AvatarImage src="/bot-avatar.png" />
                        <AvatarFallback className="bg-primary/10 text-primary"><Bot className="w-4 h-4" /></AvatarFallback>
                     </Avatar>
                     <div className="bg-card border p-3 rounded-lg rounded-tl-none shadow-sm text-sm max-w-[85%]">
                        <p>Hello! I can help you query your global data. Ask me anything, and I'll generate the SQL for you.</p>
                     </div>
                  </div>

                  {/* Example User Message */}
                  <div className="flex gap-3 flex-row-reverse">
                     <Avatar className="h-8 w-8 border">
                        <AvatarFallback className="bg-muted"><User className="w-4 h-4" /></AvatarFallback>
                     </Avatar>
                     <div className="bg-primary text-primary-foreground p-3 rounded-lg rounded-tr-none shadow-sm text-sm max-w-[85%]">
                        <p>Show me the top 5 customers by revenue.</p>
                     </div>
                  </div>

                   {/* Example Bot Response */}
                   <div className="flex gap-3">
                     <Avatar className="h-8 w-8 border">
                        <AvatarFallback className="bg-primary/10 text-primary"><Bot className="w-4 h-4" /></AvatarFallback>
                     </Avatar>
                     <div className="bg-card border p-3 rounded-lg rounded-tl-none shadow-sm text-sm max-w-[85%] space-y-2">
                        <p>Here is the query for the top 5 customers:</p>
                        <div className="bg-muted p-2 rounded border font-mono text-xs overflow-x-auto">
                           SELECT * FROM global_customers ORDER BY revenue DESC LIMIT 5;
                        </div>
                        <Button variant="secondary" size="sm" className="w-full h-7 text-xs">Insert into Editor</Button>
                     </div>
                  </div>

               </div>

               <Separator />

               {/* Input Area */}
               <div className="p-4 bg-card">
                  <div className="relative">
                     <Input 
                        placeholder="Ask a question..." 
                        className="pr-10"
                        value={chatInput}
                        onChange={(e) => setChatInput(e.target.value)}
                     />
                     <Button size="icon" variant="ghost" className="absolute right-0 top-0 h-full w-10 text-muted-foreground hover:text-primary">
                        <Send className="w-4 h-4" />
                     </Button>
                  </div>
               </div>
            </Card>
          </div>
        </ResizablePanel>

      </ResizablePanelGroup>
    </div>
  )
}
