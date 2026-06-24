import { lazy } from 'react'
import type { ReactNode } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { WithSuspense } from './WithSuspense'

const Login = lazy(() => import('./pages/login/App'))
const AppLayout = lazy(() => import('./pages/layout/AppLayout'))
const Dashboard = lazy(() => import('./pages/dashboard/App'))
const Orders = lazy(() => import('./pages/orders/App'))
const Users = lazy(() => import('./pages/users/App'))
const Content = lazy(() => import('./pages/content/App'))
const Reports = lazy(() => import('./pages/reports/App'))

const wrap = (node: ReactNode) => <WithSuspense>{node}</WithSuspense>

const router = createBrowserRouter([
  {
    path: '/',
    element: <Navigate to="/admin" replace />
  },
  {
    path: '/admin/login',
    element: wrap(<Login />)
  },
  {
    path: '/admin',
    element: wrap(<AppLayout />),
    children: [
      {
        index: true,
        element: <Navigate to="/admin/dashboard" replace />
      },
      {
        path: 'dashboard',
        element: wrap(<Dashboard />)
      },
      {
        path: 'orders',
        element: wrap(<Orders />)
      },
      {
        path: 'users',
        element: wrap(<Users />)
      },
      {
        path: 'content',
        element: wrap(<Content />)
      },
      {
        path: 'reports',
        element: wrap(<Reports />)
      }
    ]
  },
  {
    path: '*',
    element: <Navigate to="/admin" replace />
  }
])

export default router
