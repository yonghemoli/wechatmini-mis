import { configureStore, createSlice, PayloadAction } from '@reduxjs/toolkit'

// ==================== Auth Slice ====================
interface UserInfo {
  id: number
  username: string
  email: string
  avatar: string
  name?: string
  status?: string
  isSuperAdmin?: boolean
  roleId?: number | null
}

interface AuthState {
  loggedIn: boolean
  user: UserInfo | null
}

const authSlice = createSlice({
  name: 'auth',
  initialState: {
    loggedIn: !!localStorage.getItem('mis:logged_in'),
    user: null
  } as AuthState,
  reducers: {
    setLoggedIn(state, action: PayloadAction<UserInfo>) {
      state.loggedIn = true
      state.user = action.payload
      localStorage.setItem('mis:logged_in', '1')
    },
    setLoggedOut(state) {
      state.loggedIn = false
      state.user = null
      localStorage.removeItem('mis:logged_in')
    }
  }
})

export const { setLoggedIn, setLoggedOut } = authSlice.actions

// ==================== Store ====================
const store = configureStore({
  reducer: {
    auth: authSlice.reducer
  }
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
export default store
