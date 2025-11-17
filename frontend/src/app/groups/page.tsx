'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'

interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  avatar_path?: string
  nickname?: string
}

interface Group {
  id: number
  creator_id: number
  title: string
  description?: string
  created_at: string
}

export default function Groups() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [groups, setGroups] = useState<Group[]>([])
  const [loading, setLoading] = useState(true)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [newGroupTitle, setNewGroupTitle] = useState('')
  const [newGroupDescription, setNewGroupDescription] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    fetchUser()
    fetchGroups()
  }, [])

  const fetchUser = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/auth/me', {
        credentials: 'include',
      })
      if (response.ok) {
        const userData = await response.json()
        setUser(userData)
      } else {
        router.push('/login')
      }
    } catch (err) {
      router.push('/login')
    }
  }

  const fetchGroups = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/groups/user', {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setGroups(data || [])
      } else {
        setGroups([])
      }
    } catch (err) {
      console.error('Failed to fetch groups:', err)
      setGroups([])
    } finally {
      setLoading(false)
    }
  }

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newGroupTitle.trim()) return

    setCreating(true)

    try {
      const response = await fetch('http://localhost:8080/api/groups/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          title: newGroupTitle,
          description: newGroupDescription || null,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        setNewGroupTitle('')
        setNewGroupDescription('')
        setShowCreateForm(false)
        fetchGroups() // Refresh groups list
      } else {
        alert('Failed to create group')
      }
    } catch (err) {
      console.error('Failed to create group:', err)
      alert('Failed to create group')
    } finally {
      setCreating(false)
    }
  }

  const handleLogout = async () => {
    try {
      await fetch('http://localhost:8080/api/auth/logout', {
        method: 'POST',
        credentials: 'include',
      })
      router.push('/')
    } catch (err) {
      console.error('Failed to logout:', err)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">Groups</h1>
            <div className="flex items-center space-x-4">
              <button
                onClick={() => router.push('/dashboard')}
                className="text-blue-600 hover:text-blue-800"
              >
                Dashboard
              </button>
              <button
                onClick={() => router.push(`/profile/${user?.id}`)}
                className="text-blue-600 hover:text-blue-800"
              >
                Profile
              </button>
              <button
                onClick={handleLogout}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        {/* Create Group Button */}
        <div className="mb-6">
          <button
            onClick={() => setShowCreateForm(!showCreateForm)}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-md font-semibold"
          >
            {showCreateForm ? '✕ Cancel' : '+ Create New Group'}
          </button>
        </div>

        {/* Create Group Form */}
        {showCreateForm && (
          <div className="bg-white shadow rounded-lg p-6 mb-6">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">Create New Group</h2>
            <form onSubmit={handleCreateGroup}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Group Name *
                </label>
                <input
                  type="text"
                  value={newGroupTitle}
                  onChange={(e) => setNewGroupTitle(e.target.value)}
                  placeholder="Enter group name"
                  className="w-full p-3 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500"
                  required
                />
              </div>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Description (Optional)
                </label>
                <textarea
                  value={newGroupDescription}
                  onChange={(e) => setNewGroupDescription(e.target.value)}
                  placeholder="Describe your group..."
                  className="w-full p-3 border border-gray-300 rounded-md resize-none focus:ring-blue-500 focus:border-blue-500"
                  rows={3}
                />
              </div>
              <button
                type="submit"
                disabled={creating}
                className={`w-full py-3 rounded-md font-semibold ${
                  creating
                    ? 'bg-gray-400 cursor-not-allowed'
                    : 'bg-blue-600 hover:bg-blue-700'
                } text-white`}
              >
                {creating ? 'Creating...' : 'Create Group'}
              </button>
            </form>
          </div>
        )}

        {/* Groups List */}
        <div className="bg-white shadow rounded-lg">
          <div className="px-6 py-4 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900">Your Groups</h2>
          </div>

          {groups.length === 0 ? (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-400 text-6xl mb-4">👥</div>
              <p className="text-gray-500 text-lg">No groups yet</p>
              <p className="text-gray-400 text-sm mt-2">
                Create a group to connect with others!
              </p>
            </div>
          ) : (
            <div className="divide-y divide-gray-200">
              {groups.map((group) => (
                <div
                  key={group.id}
                  onClick={() => router.push(`/groups/${group.id}`)}
                  className="px-6 py-4 hover:bg-gray-50 cursor-pointer transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="text-lg font-semibold text-gray-900 mb-1">
                        {group.title}
                      </h3>
                      {group.description && (
                        <p className="text-gray-600 text-sm mb-2">{group.description}</p>
                      )}
                      <p className="text-xs text-gray-500">
                        Created {new Date(group.created_at).toLocaleDateString()}
                      </p>
                    </div>
                    <div className="ml-4">
                      <span className="text-2xl">→</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Back Button */}
        <div className="mt-6 text-center">
          <button
            onClick={() => router.push('/dashboard')}
            className="text-blue-600 hover:text-blue-800"
          >
            ← Back to Dashboard
          </button>
        </div>
      </div>
    </div>
  )
}
