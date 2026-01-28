export interface Role {
  id: number
  name: 'admin' | 'user' | 'guest'
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  role_id: number
  role: Role
  active: boolean
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface AuthError {
  error: string
  details?: Record<string, string>
}
