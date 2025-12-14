import { useState, useRef, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Send, Loader2, Search, Database, Table, Layers } from 'lucide-react'
import { api, type ToolResult } from '@/lib/api'
import { QueryResultTable } from '@/components/QueryResultTable'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

type Message = {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
  toolResults?: ToolResult[]
}

export function ChatPage() {
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  // Auto-scroll to bottom when new messages arrive
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView({ behavior: 'smooth', block: 'end' })
    }
  }, [messages])

  const handleSend = async () => {
    if (!input.trim() || isLoading) return

    const userMessage: Message = {
      id: Date.now().toString(),
      role: 'user',
      content: input.trim(),
      timestamp: new Date(),
    }

    setMessages((prev) => [...prev, userMessage])
    setInput('')
    setIsLoading(true)

    try {
      const response = await api.sendChatMessage(input.trim(), messages)
      const assistantMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: response.message,
        timestamp: new Date(),
        toolResults: response.toolResults,
      }
      setMessages((prev) => [...prev, assistantMessage])
    } catch (error) {
      console.error('Chat error:', error)
      const errorMessage: Message = {
        id: (Date.now() + 1).toString(),
        role: 'assistant',
        content: 'Sorry, something went wrong. Please try again.',
        timestamp: new Date(),
      }
      setMessages((prev) => [...prev, errorMessage])
    } finally {
      setIsLoading(false)
      // Focus input after response
      inputRef.current?.focus()
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSend()
    }
  }

  const handleSuggestionClick = (suggestion: string) => {
    setInput(suggestion)
    inputRef.current?.focus()
  }

  // Initial centered view (no messages yet)
  if (messages.length === 0 && !isLoading) {
    return (
      <div className="flex h-screen w-full flex-col">
        {/* Hero Section */}
        <div className="flex flex-1 flex-col items-center justify-center p-8">
          <div className="w-full max-w-3xl space-y-12">
            {/* Title */}
            <div className="space-y-4 text-center">
           <img src="/data-sync.png" alt="DataSync logo" className="mx-auto h-18 w-18 mb-5" />
              <h1 className="text-5xl font-bold tracking-tight">
                Ask the Data
              </h1>
              <p className="text-lg text-muted-foreground">
                Explore your data sources with natural language
              </p>
            </div>

            {/* Search Input */}
            <div className="relative">
              <Search className="absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-muted-foreground" />
              <Input
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Ask anything about your data..."
                className="h-14 pl-12 pr-14 text-base shadow-lg border-2 focus-visible:ring-2"
                autoFocus
              />
              <Button
                onClick={handleSend}
                disabled={!input.trim()}
                size="icon"
                className="absolute right-2 top-1/2 h-10 w-10 -translate-y-1/2"
              >
                <Send className="h-4 w-4" />
              </Button>
            </div>

            {/* Suggestions */}
            <div className="space-y-4">
              <p className="text-sm font-medium text-muted-foreground">
                Try asking:
              </p>
              <div className="grid gap-3 sm:grid-cols-2">
                <button
                  onClick={() => handleSuggestionClick('Show me all customers')}
                  className="group flex items-start gap-3 rounded-lg border bg-card p-4 text-left transition-all hover:border-primary/50 hover:bg-accent"
                >
                  <Database className="mt-0.5 h-5 w-5 text-primary" />
                  <div className="space-y-1">
                    <p className="text-sm font-medium">Show me all customers</p>
                    <p className="text-xs text-muted-foreground">Query customer data</p>
                  </div>
                </button>
                <button
                  onClick={() => handleSuggestionClick('What tables do I have?')}
                  className="group flex items-start gap-3 rounded-lg border bg-card p-4 text-left transition-all hover:border-primary/50 hover:bg-accent"
                >
                  <Table className="mt-0.5 h-5 w-5 text-primary" />
                  <div className="space-y-1">
                    <p className="text-sm font-medium">What tables do I have?</p>
                    <p className="text-xs text-muted-foreground">Explore metadata</p>
                  </div>
                </button>
                <button
                  onClick={() => handleSuggestionClick('List all schemas in PostgreSQL')}
                  className="group flex items-start gap-3 rounded-lg border bg-card p-4 text-left transition-all hover:border-primary/50 hover:bg-accent"
                >
                  <Layers className="mt-0.5 h-5 w-5 text-primary" />
                  <div className="space-y-1">
                    <p className="text-sm font-medium">List all schemas</p>
                    <p className="text-xs text-muted-foreground">Browse structure</p>
                  </div>
                </button>
                <button
                  onClick={() => handleSuggestionClick('How many products are there?')}
                  className="group flex items-start gap-3 rounded-lg border bg-card p-4 text-left transition-all hover:border-primary/50 hover:bg-accent"
                >
                  <Database className="mt-0.5 h-5 w-5 text-primary" />
                  <div className="space-y-1">
                    <p className="text-sm font-medium">Count products</p>
                    <p className="text-xs text-muted-foreground">Aggregate data</p>
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  // Thread view (after first message)
  return (
    <div className="flex h-screen flex-col">
      {/* Thread Messages */}
      <div className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-5xl px-4 py-8 space-y-8">
          {messages.map((message) => (
            <div key={message.id} className="space-y-4">
              {/* User Query */}
              {message.role === 'user' && (
                <div className="flex items-start gap-3">
                  <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                    <Search className="h-4 w-4 text-primary" />
                  </div>
                  <div className="flex-1 pt-1">
                    <p className="text-lg font-medium">{message.content}</p>
                  </div>
                </div>
              )}

              {/* Assistant Response */}
              {message.role === 'assistant' && (
                <Card className="border-0 shadow-none bg-transparent">
                  <div className="space-y-4">
                    <Separator />
                    <div className="prose prose-lg max-w-none dark:prose-invert">
                      <ReactMarkdown
                        remarkPlugins={[remarkGfm]}
                        components={{
                          p: ({ children }) => <p className="mb-4 text-lg leading-8">{children}</p>,
                          ul: ({ children }) => <ul className="mb-4 ml-6 list-disc space-y-2 text-lg">{children}</ul>,
                          ol: ({ children }) => <ol className="mb-4 ml-6 list-decimal space-y-2 text-lg">{children}</ol>,
                          li: ({ children }) => <li className="text-lg leading-8">{children}</li>,
                          h1: ({ children }) => <h1 className="mb-4 mt-6 text-3xl font-bold">{children}</h1>,
                          h2: ({ children }) => <h2 className="mb-3 mt-5 text-2xl font-semibold">{children}</h2>,
                          h3: ({ children }) => <h3 className="mb-2 mt-4 text-xl font-semibold">{children}</h3>,
                          code: ({ className, children }) => {
                            const isInline = !className;
                            return isInline ? (
                              <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-base">{children}</code>
                            ) : (
                              <code className={className}>{children}</code>
                            );
                          },
                          pre: ({ children }) => <pre className="mb-4 overflow-x-auto rounded-lg bg-muted p-4 text-base">{children}</pre>,
                          blockquote: ({ children }) => <blockquote className="border-l-4 border-muted-foreground pl-4 italic text-lg">{children}</blockquote>,
                          a: ({ href, children }) => <a href={href} className="text-primary underline hover:no-underline text-lg">{children}</a>,
                          strong: ({ children }) => <strong className="font-semibold">{children}</strong>,
                          em: ({ children }) => <em className="italic">{children}</em>,
                          table: ({ children }) => (
                            <div className="my-6 w-full overflow-x-auto">
                              <table className="w-full border-collapse border border-border text-base">
                                {children}
                              </table>
                            </div>
                          ),
                          thead: ({ children }) => <thead className="bg-muted">{children}</thead>,
                          tbody: ({ children }) => <tbody>{children}</tbody>,
                          tr: ({ children }) => <tr className="border-b border-border">{children}</tr>,
                          th: ({ children }) => <th className="border border-border px-4 py-3 text-left font-semibold text-base">{children}</th>,
                          td: ({ children }) => <td className="border border-border px-4 py-3 text-base">{children}</td>,
                        }}
                      >
                        {message.content}
                      </ReactMarkdown>
                    </div>
                    {/* Render tool results (query results, metadata, etc.) */}
                    {message.toolResults && message.toolResults.length > 0 && (
                      <div className="space-y-3">
                        {message.toolResults.map((toolResult, idx) => (
                          <div key={idx}>
                            {toolResult.toolName === 'executeGlobalQuery' && (
                              <QueryResultTable data={toolResult.data} />
                            )}
                            {toolResult.toolName === 'listGlobalTables' && toolResult.data.tables && (
                              <Card className="mt-2 border-purple-200 dark:border-purple-800">
                                <div className="p-4">
                                  <h4 className="text-sm font-semibold mb-3 flex items-center gap-2">
                                    <Table className="h-4 w-4" />
                                    Available Global Tables ({toolResult.data.count})
                                  </h4>
                                  <div className="space-y-2">
                                    {toolResult.data.tables.map((table: any, tableIdx: number) => (
                                      <div key={tableIdx} className="flex items-start gap-2 text-sm border-l-2 border-purple-300 dark:border-purple-700 pl-3 py-1">
                                        <Database className="h-4 w-4 mt-0.5 text-purple-600 dark:text-purple-400" />
                                        <div>
                                          <div className="font-medium">{table.name}</div>
                                          {table.description && (
                                            <div className="text-xs text-muted-foreground">{table.description}</div>
                                          )}
                                        </div>
                                      </div>
                                    ))}
                                  </div>
                                </div>
                              </Card>
                            )}
                            {toolResult.toolName === 'discoverMetadata' && (
                              <Card className="mt-2 border-green-200 dark:border-green-800">
                                <div className="p-4">
                                  <pre className="text-xs overflow-x-auto">
                                    {JSON.stringify(toolResult.data, null, 2)}
                                  </pre>
                                </div>
                              </Card>
                            )}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </Card>
              )}
            </div>
          ))}

          {/* Loading State */}
          {isLoading && (
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10">
                  <Loader2 className="h-4 w-4 animate-spin text-primary" />
                </div>
                <div className="flex-1 pt-1">
                  <div className="flex items-center gap-2">
                    <div className="h-2 w-2 animate-pulse rounded-full bg-primary" />
                    <div className="h-2 w-2 animate-pulse rounded-full bg-primary delay-75" />
                    <div className="h-2 w-2 animate-pulse rounded-full bg-primary delay-150" />
                  </div>
                </div>
              </div>
            </div>
          )}

          <div ref={scrollRef} />
        </div>
      </div>

      {/* Fixed Input at Bottom */}
      <div className="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="mx-auto max-w-5xl px-4 py-4">
          <div className="relative">
            <Input
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a follow-up question..."
              className="h-12 pl-4 pr-12 text-base border-2 rounded-xl"
              disabled={isLoading}
            />
            <Button
              onClick={handleSend}
              disabled={!input.trim() || isLoading}
              size="icon"
              className="absolute right-2 top-1/2 h-9 w-9 -translate-y-1/2 rounded-lg"
            >
              {isLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Send className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
