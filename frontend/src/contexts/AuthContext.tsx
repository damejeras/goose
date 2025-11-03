import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { apiClient } from '../lib/apiClient'
import type { User as ApiUser } from '../lib/api/v1/api_pb'

interface User {
  id: string
  email: string
  googleId: string
  name: string
}

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  loginWithGoogle: (googleIdToken: string) => Promise<void>
  logout: () => Promise<void>
  isLoading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

// Convert API User (with bigint id) to our User interface
function convertApiUser(apiUser: ApiUser): User {
  return {
    id: apiUser.id.toString(),
    email: apiUser.email,
    googleId: apiUser.googleId,
    name: apiUser.name,
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    // Check if user is authenticated on mount
    const checkAuth = async () => {
      const token = localStorage.getItem('auth_token')
      if (!token) {
        setIsLoading(false)
        return
      }

      try {
        // Validate token and get current user from backend
        const response = await apiClient.auth.getCurrentUser({})
        if (response.user) {
          const userData = convertApiUser(response.user)
          setUser(userData)
          localStorage.setItem('auth_user', JSON.stringify(userData))
        }
      } catch (error) {
        console.error('Failed to get current user:', error)
        // Token is invalid, clear it
        apiClient.clearToken()
        localStorage.removeItem('auth_user')
      } finally {
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [])

  const loginWithGoogle = async (googleIdToken: string) => {
    try {
      // Call backend Login endpoint with Google ID token
      const response = await apiClient.auth.login({
        googleIdToken,
      })

      if (!response.user || !response.jwt) {
        throw new Error('Invalid login response')
      }

      // Store JWT token
      apiClient.setToken(response.jwt)

      // Store user data
      const userData = convertApiUser(response.user)
      setUser(userData)
      localStorage.setItem('auth_user', JSON.stringify(userData))
    } catch (error) {
      console.error('Login failed:', error)
      throw error
    }
  }

  const logout = async () => {
    try {
      // Call backend logout endpoint
      await apiClient.auth.logout({})
    } catch (error) {
      console.error('Logout failed:', error)
      // Continue with local logout even if backend call fails
    } finally {
      // Clear local state
      setUser(null)
      apiClient.clearToken()
      localStorage.removeItem('auth_user')
    }
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isAuthenticated: !!user,
        loginWithGoogle,
        logout,
        isLoading,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
