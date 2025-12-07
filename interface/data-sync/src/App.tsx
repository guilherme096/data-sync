import { RouterProvider, createRouter, createRootRoute, createRoute } from '@tanstack/react-router'
import { MainLayout } from './layouts/MainLayout'
import { QueryPage } from './pages/QueryPage'

// 1. Create the root route (wraps everything)
const rootRoute = createRootRoute({
  component: MainLayout,
})

// 2. Create the index route (Home / QueryPage)
const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: QueryPage,
})

// Placeholders for other routes
const schemaRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/schema',
  component: () => (
    <div className="flex flex-col items-center justify-center min-h-[50vh] space-y-4">
      <h2 className="text-2xl font-bold">Schema Studio</h2>
      <p className="text-muted-foreground">Mapping interface coming soon.</p>
    </div>
  ),
})

const inventoryRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/inventory',
  component: () => (
    <div className="flex flex-col items-center justify-center min-h-[50vh] space-y-4">
      <h2 className="text-2xl font-bold">Data Inventory</h2>
      <p className="text-muted-foreground">Catalog browser coming soon.</p>
    </div>
  ),
})

// 3. Create the route tree
const routeTree = rootRoute.addChildren([indexRoute, schemaRoute, inventoryRoute])

// 4. Create the router
const router = createRouter({ routeTree })

// 5. Register the router for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

function App() {
  return <RouterProvider router={router} />
}

export default App