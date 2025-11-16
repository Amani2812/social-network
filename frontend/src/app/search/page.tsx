'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'

interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  nickname?: string
  avatar_path?: string
  is_public: boolean
}

export default function Search() {
  const [searchQuery, setSearchQuery] = useState('')
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const router = useRouter()

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!searchQuery.trim()) return

    setLoading(true)
    setError('')

    try {
      const response = await fetch(
        `http://localhost:8080/api/users/search?q=${encodeURIComponent(searchQuery)}`,
        {
          credentials: 'include',
        }
      )

      if (response.ok) {
        const data = await response.json()
        setUsers(data || [])
      } else {
        setError('Failed to search users')
      }
    } catch (err) {
      setError('Network error. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">Search Users</h1>
            <button
              onClick={() => router.push('/dashboard')}
              className="text-blue-600 hover:text-blue-800"
            >
              Back to Dashboard
            </button>
          </div>
        </div>
      </header>

      <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        {/* Search Form */}
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <form onSubmit={handleSearch}>
            <div className="flex gap-4">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search by name or email..."
                className="flex-1 p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
              />
              <button
                type="submit"
                disabled={loading}
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-md disabled:opacity-50"
              >
                {loading ? 'Searching...' : 'Search'}
              </button>
            </div>
          </form>

          {error && (
            <div className="mt-4 text-red-600 text-sm">
              {error}
            </div>
          )}
        </div>

        {/* Search Results */}
        <div className="bg-white shadow rounded-lg p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            Search Results
          </h2>

          {users.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              {searchQuery ? 'No users found' : 'Enter a search query to find users'}
            </div>
          ) : (
            <div className="space-y-4">
              {users.map((user) => (
                <div
                  key={user.id}
                  className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer"
                  onClick={() => router.push(`/profile/${user.id}`)}
                >
                  <div className="flex items-center">
                    <div className="w-12 h-12 bg-gray-300 rounded-full flex items-center justify-center mr-4">
                      {user.avatar_path ? (
                        <img
                          src={user.avatar_path}
                          alt="Avatar"
                          className="w-full h-full rounded-full object-cover"
                        />
                      ) : (
                        <span className="text-lg text-gray-600">
                          {user.first_name[0]}{user.last_name[0]}
                        </span>
                      )}
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900">
                        {user.first_name} {user.last_name}
                      </p>
                      {user.nickname && (
                        <p className="text-sm text-gray-500">@{user.nickname}</p>
                      )}
                      <p className="text-sm text-gray-600">{user.email}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <span
                      className={`text-xs px-2 py-1 rounded-full ${
                        user.is_public
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      }`}
                    >
                      {user.is_public ? 'Public' : 'Private'}
                    </span>
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        router.push(`/profile/${user.id}`)
                      }}
                      className="text-blue-600 hover:text-blue-800 text-sm font-medium"
                    >
                      View Profile →
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
