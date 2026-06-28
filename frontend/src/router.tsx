import { lazy } from 'react'
import type { ReactNode } from 'react'
import { createBrowserRouter, Navigate } from 'react-router-dom'
import { WithSuspense } from './WithSuspense'

const Login = lazy(() => import('./pages/login/App'))
const AppLayout = lazy(() => import('./pages/layout/AppLayout'))
const Dashboard = lazy(() => import('./pages/dashboard/App'))
const Orders = lazy(() => import('./pages/orders/App'))
const Users = lazy(() => import('./pages/users/App'))
const ServiceTypes = lazy(() => import('./pages/serviceTypes/App'))
const Services = lazy(() => import('./pages/services/App'))
const Accounts = lazy(() => import('./pages/accounts/App'))
const Shops = lazy(() => import('./pages/shops/App'))
const FAQs = lazy(() => import('./pages/faqs/App'))
const Chat = lazy(() => import('./pages/chat/App'))
const Profile = lazy(() => import('./pages/profile/App'))

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
        path: 'service-types',
        element: wrap(<ServiceTypes />)
      },
      {
        path: 'services',
        element: wrap(<Services />)
      },
      {
        path: 'accounts',
        element: wrap(<Accounts />)
      },
      {
        path: 'shops',
        element: wrap(<Shops />)
      },
      {
        path: 'faqs',
        element: wrap(<FAQs />)
      },
      {
        path: 'chat',
        element: wrap(<Chat />)
      },
      {
        path: 'profile',
        element: wrap(<Profile />)
      }
    ]
  },
  {
    path: '*',
    element: <Navigate to="/admin" replace />
  }
])

export default router
