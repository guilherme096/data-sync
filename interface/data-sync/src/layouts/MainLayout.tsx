import { Link, Outlet } from '@tanstack/react-router'
import { Database, Search, Workflow } from 'lucide-react'
import { cn } from '@/lib/utils'

export function MainLayout() {
  return (
    <div className="flex h-screen w-full bg-background text-foreground">
      {/* Sidebar */}
      <aside className="w-64 border-r bg-muted/10 hidden md:flex flex-col">
        <div className="p-6 flex items-center gap-2">
           <div className="h-6 w-6 rounded bg-primary" />
           <h1 className="text-xl font-bold tracking-tight">DataSync</h1>
        </div>
        <nav className="flex-1 px-4 space-y-2">
           <NavLink to="/" icon={<Search className="w-4 h-4" />} label="Query" />
           <NavLink to="/schema" icon={<Workflow className="w-4 h-4" />} label="Schema Studio" />
           <NavLink to="/inventory" icon={<Database className="w-4 h-4" />} label="Inventory" />
        </nav>
        <div className="p-4 border-t text-xs text-muted-foreground text-center">
          v0.1.0-alpha
        </div>
      </aside>
      
      {/* Main Content */}
      <main className="flex-1 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}

function NavLink({ to, icon, label }: { to: string; icon: React.ReactNode; label: string }) {
  return (
    <Link 
      to={to} 
      className="flex items-center gap-3 px-3 py-2 text-sm font-medium rounded-md text-muted-foreground hover:bg-muted hover:text-foreground transition-colors [&.active]:bg-muted [&.active]:text-foreground [&.active]:font-semibold"
    >
      {icon}
      {label}
    </Link>
  )
}
