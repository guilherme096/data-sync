import { useState, useRef, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Send, Loader2, Search, Database, Table, Layers } from 'lucide-react'
import { api } from '@/lib/api'

type Message = {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
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
        content: response,
        timestamp: new Date(),
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
      {/* Fixed Header with Search */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="mx-auto max-w-3xl px-4 py-3">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              ref={inputRef}
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder="Ask a follow-up question..."
              className="h-11 pl-10 pr-12 border-2"
              disabled={isLoading}
            />
            <Button
              onClick={handleSend}
              disabled={!input.trim() || isLoading}
              size="icon"
              className="absolute right-1.5 top-1/2 h-8 w-8 -translate-y-1/2"
            >
              <Send className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </div>

      {/* Thread Messages */}
      <div className="flex-1 overflow-y-auto">
        <div className="mx-auto max-w-3xl px-4 py-8 space-y-8">
          {messages.map((message, index) => (
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
                    <div className="prose prose-sm max-w-none">
                      <p className="text-sm leading-7 text-foreground whitespace-pre-wrap">
                        {message.content}
                      </p>
                    </div>
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
    </div>
  )
}
