import { RouterProvider, createRouter, createRootRoute, createRoute } from '@tanstack/react-router'
import { MainLayout } from './layouts/MainLayout'
import { QueryPage } from './pages/QueryPage'
import { InventoryPage } from './pages/InventoryPage'
import { SchemaStudioPage } from './pages/SchemaStudioPage'

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

// Schema Studio route
const schemaRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/schema',
  component: SchemaStudioPage,
})

const inventoryRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/inventory',
  component: InventoryPage,
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