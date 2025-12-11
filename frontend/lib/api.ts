import { api } from '@/lib/api'
import type { User, LoginRequest, RegisterRequest, AuthResponse } from '../types'

export const authAPI = {
  // Register new user
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/register', data)
    return response.data
  },

  // Login user
  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const response = await api.post('/auth/login', data)
    return response.data
  },

  // Logout user
  logout: async (): Promise<void> => {
    await api.post('/auth/logout')
  },

  // Get current user
  me: async (): Promise<User> => {
    const response = await api.get('/auth/me')
    return response.data
  },
}
