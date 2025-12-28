'use client'

import { useState } from 'react'
import { authAPI } from '../lib/api'

interface UserDeleteButtonProps {
  userId: number
  userName: string
  onDeleted: () => void
}

export function UserDeleteButton({ userId, userName, onDeleted }: UserDeleteButtonProps) {
  const [showConfirm, setShowConfirm] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleDelete = async () => {
    setLoading(true)
    setError(null)

    try {
      await authAPI.deleteUser(userId)
      setShowConfirm(false)
      onDeleted()
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete user')
    } finally {
      setLoading(false)
    }
  }

  if (!showConfirm) {
    return (
      <button
        onClick={() => setShowConfirm(true)}
        className="text-red-600 hover:text-red-900"
      >
        Delete
      </button>
    )
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg max-w-md w-full p-6">
        <h2 className="text-xl font-bold mb-4">Confirm Delete</h2>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 p-3 rounded mb-4">
            {error}
          </div>
        )}

        <p className="text-gray-700 mb-6">
          Are you sure you want to delete user <strong>{userName}</strong>? This action cannot be
          undone.
        </p>

        <div className="flex gap-2">
          <button
            onClick={handleDelete}
            disabled={loading}
            className="flex-1 bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700 disabled:opacity-50"
          >
            {loading ? 'Deleting...' : 'Delete User'}
          </button>
          <button
            onClick={() => setShowConfirm(false)}
            className="flex-1 bg-gray-200 text-gray-800 px-4 py-2 rounded hover:bg-gray-300"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  )
}
