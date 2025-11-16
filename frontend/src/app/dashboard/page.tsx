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

interface Post {
  id: number
  user_id: number
  content: string
  image_path?: string
  privacy: string
  created_at: string
  user?: User
}

export default function Dashboard() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [newPost, setNewPost] = useState('')
  const [privacy, setPrivacy] = useState('public')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchUser()
    fetchFeed()
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

  const fetchFeed = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/posts/feed', {
        credentials: 'include',
      })
      if (response.ok) {
        const postsData = await response.json()
        setPosts(postsData)
      }
    } catch (err) {
      console.error('Failed to fetch feed:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newPost.trim()) return

    try {
      const response = await fetch('http://localhost:8080/api/posts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: newPost,
          privacy: privacy,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        setNewPost('')
        fetchFeed() // Refresh feed
      }
    } catch (err) {
      console.error('Failed to create post:', err)
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
            <h1 className="text-2xl font-bold text-gray-900">Social Network</h1>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                Welcome, {user?.first_name} {user?.last_name}
              </span>
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
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Create Post */}
          <div className="lg:col-span-2">
            <div className="bg-white shadow rounded-lg p-6 mb-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-4">Create Post</h2>
              <form onSubmit={handleCreatePost}>
                <textarea
                  value={newPost}
                  onChange={(e) => setNewPost(e.target.value)}
                  placeholder="What's on your mind?"
                  className="w-full p-3 border border-gray-300 rounded-md resize-none focus:ring-blue-500 focus:border-blue-500"
                  rows={4}
                />
                <div className="mt-4 flex items-center justify-between">
                  <select
                    value={privacy}
                    onChange={(e) => setPrivacy(e.target.value)}
                    className="border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option value="public">Public</option>
                    <option value="almost_private">Followers Only</option>
                    <option value="private">Private</option>
                  </select>
                  <button
                    type="submit"
                    className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-md"
                  >
                    Post
                  </button>
                </div>
              </form>
            </div>

            {/* Feed */}
            <div className="space-y-6">
              {posts.length === 0 ? (
                <div className="bg-white shadow rounded-lg p-6 text-center">
                  <p className="text-gray-500">No posts yet. Follow some users to see their posts!</p>
                </div>
              ) : (
                posts.map((post) => (
                  <div key={post.id} className="bg-white shadow rounded-lg p-6">
                    <div className="flex items-center mb-4">
                      <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                        {post.user?.avatar_path ? (
                          <img src={post.user.avatar_path} alt="Avatar" className="w-full h-full rounded-full object-cover" />
                        ) : (
                          <span className="text-sm text-gray-600">
                            {post.user?.first_name[0]}{post.user?.last_name[0]}
                          </span>
                        )}
                      </div>
                      <div>
                        <p className="font-semibold text-gray-900">
                          {post.user?.first_name} {post.user?.last_name}
                        </p>
                        <p className="text-sm text-gray-500">
                          {new Date(post.created_at).toLocaleDateString()}
                        </p>
                      </div>
                    </div>
                    <p className="text-gray-700 mb-4">{post.content}</p>
                    {post.image_path && (
                      <img src={post.image_path} alt="Post image" className="w-full rounded-md mb-4" />
                    )}
                    <div className="flex items-center justify-between text-sm text-gray-500">
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        post.privacy === 'public' ? 'bg-green-100 text-green-800' :
                        post.privacy === 'almost_private' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {post.privacy === 'public' ? 'Public' :
                         post.privacy === 'almost_private' ? 'Followers Only' : 'Private'}
                      </span>
                      <button className="text-blue-600 hover:text-blue-800">
                        Comment
                      </button>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* User Info */}
            <div className="bg-white shadow rounded-lg p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Your Profile</h3>
              <div className="flex items-center">
                <div className="w-16 h-16 bg-gray-300 rounded-full flex items-center justify-center mr-4">
                  {user?.avatar_path ? (
                    <img src={user.avatar_path} alt="Avatar" className="w-full h-full rounded-full object-cover" />
                  ) : (
                    <span className="text-xl text-gray-600">
                      {user?.first_name[0]}{user?.last_name[0]}
                    </span>
                  )}
                </div>
                <div>
                  <p className="font-semibold text-gray-900">
                    {user?.first_name} {user?.last_name}
                  </p>
                  <p className="text-gray-600">{user?.email}</p>
                  {user?.nickname && (
                    <p className="text-gray-500">@{user.nickname}</p>
                  )}
                </div>
              </div>
            </div>

            {/* Quick Actions */}
            <div className="bg-white shadow rounded-lg p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
              <div className="space-y-2">
                <button
                  onClick={() => router.push('/search')}
                  className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
                >
                  🔍 Search Users
                </button>
                <button className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">
                  👥 Create Group
                </button>
                <button className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">
                  💬 Messages
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
