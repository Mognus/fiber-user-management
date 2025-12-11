export interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  role: 'admin' | 'user' | 'guest'
  active: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface AuthError {
  error: string
  details?: Record<string, string>
}
