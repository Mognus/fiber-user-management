import { fetcher, mutateFetch } from '@/lib/api/fetcher'
import type { User, LoginRequest, AuthResponse } from '../types'

export const authAPI = {
  // Login user
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    return mutateFetch<AuthResponse>('/auth/login', {
      method: 'POST',
      body: data,
    })
  },

  // Logout user
  logout: async (): Promise<void> => {
    await mutateFetch<void>('/auth/logout', { method: 'POST' })
  },

  // Get current user
  me: async (): Promise<User> => {
    return fetcher<User>('/auth/me')
  },
}
