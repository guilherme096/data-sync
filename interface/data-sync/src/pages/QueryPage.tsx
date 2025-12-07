import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Play, Eraser, Send, Bot, User, Sparkles } from 'lucide-react'
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "@/components/ui/resizable"
import CodeMirror from '@uiw/react-codemirror'
import { sql as sqlLang } from '@codemirror/lang-sql'

export function QueryPage() {
  const [sqlCode, setSqlCode] = useState('-- Write your SQL query here\nSELECT * FROM global_users LIMIT 10;')
  const [chatInput, setChatInput] = useState('')

  return (
    <div className="h-full w-full bg-muted/20">
      <ResizablePanelGroup direction="horizontal" className="h-full w-full rounded-lg border">
        
        {/* LEFT PANEL: SQL Workstation */}
        <ResizablePanel defaultSize={70} minSize={30}>
          <div className="h-full p-4 flex flex-col gap-4 min-w-0">
            {/* Toolbar */}
            <div className="flex items-center justify-between bg-muted/50 p-2 rounded-t-md border-b">
               <div className="flex items-center gap-2 px-2">
                  <h2 className="text-sm font-medium tracking-tight text-muted-foreground">SQL Editor</h2>
               </div>
               <div className="flex items-center gap-2">
                  <Button variant="ghost" size="sm" className="h-7 text-muted-foreground hover:text-foreground" onClick={() => setSqlCode('')}>
                     <Eraser className="w-3.5 h-3.5 mr-2" />
                     Clear
                  </Button>
                  <Button size="sm" className="h-7">
                     <Play className="w-3.5 h-3.5 mr-2" />
                     Run
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
                      <span className="text-xs text-muted-foreground">0 rows</span>
                  </div>
                  <div className="flex-1 p-4 overflow-auto">
                      <div className="flex items-center justify-center h-full text-sm text-muted-foreground font-mono">
                         No query executed yet.
                      </div>
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
                        <Button variant="secondary" size="xs" className="w-full h-7 text-xs">Insert into Editor</Button>
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
