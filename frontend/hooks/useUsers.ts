'use client'

import { useState, useEffect } from 'react'
import { authAPI } from '../lib/api'
import type { User, UsersQueryParams, UsersResponse } from '../types'

export function useUsers(params?: UsersQueryParams) {
  const [data, setData] = useState<UsersResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchUsers = async () => {
    try {
      setLoading(true)
      setError(null)
      const response = await authAPI.getUsers(params)
      setData(response)
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [params?.page, params?.limit, params?.role, params?.active])

  return {
    users: data?.users || [],
    total: data?.total || 0,
    page: data?.page || 1,
    limit: data?.limit || 20,
    loading,
    error,
    refetch: fetchUsers,
  }
}

export function useUser(id: number | null) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id) {
      setUser(null)
      return
    }

    const fetchUser = async () => {
      try {
        setLoading(true)
        setError(null)
        const data = await authAPI.getUser(id)
        setUser(data)
      } catch (err: any) {
        setError(err.response?.data?.error || 'Failed to fetch user')
      } finally {
        setLoading(false)
      }
    }

    fetchUser()
  }, [id])

  return { user, loading, error }
}
