'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter } from 'next/navigation'
import Comments from '@/components/Comments'

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
  likes?: number
  dislikes?: number
  user_reaction?: 'like' | 'dislike' | null
}

export default function Dashboard() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [newPost, setNewPost] = useState('')
  const [privacy, setPrivacy] = useState('public')
  const [loading, setLoading] = useState(true)
  const [selectedImage, setSelectedImage] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const [uploading, setUploading] = useState(false)
  const [unreadCount, setUnreadCount] = useState(0)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [notifications, setNotifications] = useState<Array<{id: number, message: string}>>([])
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const isConnectingRef = useRef(false)
  const MAX_RECONNECT_ATTEMPTS = 10

  useEffect(() => {
    fetchUser()
    fetchFeed()
    fetchUnreadCount()
    
    // Poll for new notifications every 30 seconds
    const interval = setInterval(fetchUnreadCount, 30000)
    return () => clearInterval(interval)
  }, [])

  useEffect(() => {
    if (user) {
      connectWebSocket()
    }
    
    return () => {
      // Cleanup on unmount
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
    }
  }, [user])

  const connectWebSocket = () => {
    // Only connect if user is authenticated
    if (!user) {
      console.log('⏸️ WebSocket connection skipped - user not authenticated')
      return
    }
    
    // Prevent multiple simultaneous connection attempts
    if (isConnectingRef.current) {
      console.log('⏳ Connection attempt already in progress')
      return
    }
    
    // Check if already connected
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('✅ WebSocket already connected')
      return
    }
    
    // Check reconnection attempts
    if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
      console.error('❌ Max reconnection attempts reached. Please refresh the page.')
      return
    }
    
    isConnectingRef.current = true
    console.log(`🔄 Connecting to WebSocket... (Attempt ${reconnectAttemptsRef.current + 1})`)
    
    try {
      const websocket = new WebSocket('ws://localhost:8080/ws')
      
      websocket.onopen = () => {
        console.log('✅ Dashboard WebSocket connected successfully')
        reconnectAttemptsRef.current = 0 // Reset counter on successful connection
        isConnectingRef.current = false
      }

      // Handle ping from server - browser automatically responds with pong
      websocket.addEventListener('ping', () => {
        console.log('🏓 Received ping from server')
      })
      
      websocket.onmessage = (event) => {
        // Handle multiple messages in one frame (separated by newlines)
        const messages = event.data.trim().split('\n')

        messages.forEach((messageStr: string) => {
          if (!messageStr.trim()) return // Skip empty messages

          try {
            const data = JSON.parse(messageStr)
            console.log('📥 Received:', data)

            // Handle new post notifications
            if (data.type === 'new_post') {
              // Refresh feed to show new post
              fetchFeed()
            }

            // Handle real-time notifications
            if (data.type === 'notification') {
              // Show toast notification
              const notifId = Date.now()
              setNotifications(prev => [...prev, { id: notifId, message: data.content }])

              // Update unread count
              fetchUnreadCount()

              // Auto-remove notification after 5 seconds
              setTimeout(() => {
                setNotifications(prev => prev.filter(n => n.id !== notifId))
              }, 5000)
            }
          } catch (error) {
            console.error('Failed to parse WebSocket message:', messageStr, error)
          }
        })
      }
      
      websocket.onerror = (error) => {
        console.error('❌ WebSocket connection error:', error)
        console.log('💡 This may be due to authentication issues or server unavailability')
        isConnectingRef.current = false
      }
      
      websocket.onclose = (event) => {
        console.log(`🔌 WebSocket closed: Code ${event.code} - ${event.reason || 'No reason provided'}`)
        
        // Check if it's an authentication error (code 1006 or 1008)
        if (event.code === 1006 || event.code === 1008) {
          console.warn('⚠️ WebSocket closed due to possible authentication issue')
          console.log('💡 Make sure you are logged in and have a valid session')
        }
        
        isConnectingRef.current = false
        wsRef.current = null
        setWs(null)
        
        // Only attempt to reconnect if user is still authenticated and we haven't exceeded max attempts
        if (user && reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
          reconnectAttemptsRef.current++
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current - 1), 30000)
          console.log(`⏳ Reconnecting in ${delay / 1000} seconds... (Attempt ${reconnectAttemptsRef.current}/${MAX_RECONNECT_ATTEMPTS})`)
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connectWebSocket()
          }, delay)
        } else if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
          console.error('❌ Max reconnection attempts reached. WebSocket will not reconnect automatically.')
          console.log('💡 Please refresh the page to retry the connection')
        }
      }
      
      wsRef.current = websocket
      setWs(websocket)
    } catch (error) {
      console.error('❌ Failed to create WebSocket connection:', error)
      isConnectingRef.current = false
    }
  }

  const fetchUnreadCount = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/notifications/unread', {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setUnreadCount(data.count || 0)
      }
    } catch (err) {
      console.error('Failed to fetch unread count:', err)
    }
  }

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
        setPosts(postsData || [])
      } else {
        setPosts([])
      }
    } catch (err) {
      console.error('Failed to fetch feed:', err)
      setPosts([])
    } finally {
      setLoading(false)
    }
  }

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      // Validate file type
      const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
      if (!validTypes.includes(file.type)) {
        alert('Please select a valid image file (JPEG, PNG, or GIF)')
        return
      }

      // Validate file size (max 10MB)
      if (file.size > 10 * 1024 * 1024) {
        alert('File size must be less than 10MB')
        return
      }

      setSelectedImage(file)
      
      // Create preview
      const reader = new FileReader()
      reader.onloadend = () => {
        setImagePreview(reader.result as string)
      }
      reader.readAsDataURL(file)
    }
  }

  const handleRemoveImage = () => {
    setSelectedImage(null)
    setImagePreview(null)
  }

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newPost.trim()) return

    setUploading(true)

    try {
      let imagePath: string | null = null

      // Upload image first if selected
      if (selectedImage) {
        const formData = new FormData()
        formData.append('image', selectedImage)

        const uploadResponse = await fetch('http://localhost:8080/api/upload', {
          method: 'POST',
          body: formData,
          credentials: 'include',
        })

        if (uploadResponse.ok) {
          const uploadData = await uploadResponse.json()
          imagePath = uploadData.path
        } else {
          alert('Failed to upload image')
          setUploading(false)
          return
        }
      }

      // Create post with or without image
      const response = await fetch('http://localhost:8080/api/posts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: newPost,
          privacy: privacy,
          image_path: imagePath,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        const newPostData = await response.json()
        
        // Add post to feed immediately
        setPosts(prev => [newPostData, ...prev])
        
        setNewPost('')
        setSelectedImage(null)
        setImagePreview(null)
      }
    } catch (err) {
      console.error('Failed to create post:', err)
    } finally {
      setUploading(false)
    }
  }

  const handlePostReaction = async (postId: number, reaction: 'like' | 'dislike') => {
    try {
      const response = await fetch('http://localhost:8080/api/posts/react', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          post_id: postId,
          reaction: reaction,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        // Refresh feed to get updated counts
        fetchFeed()
      }
    } catch (err) {
      console.error('Failed to react to post:', err)
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

  const removeNotification = (id: number) => {
    setNotifications(prev => prev.filter(n => n.id !== id))
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Toast Notifications */}
      <div className="fixed top-4 right-4 z-50 space-y-2">
        {notifications.map((notif) => (
          <div
            key={notif.id}
            className="bg-blue-600 text-white px-6 py-4 rounded-lg shadow-lg flex items-center justify-between min-w-[300px] animate-slide-in"
          >
            <div className="flex items-center">
              <span className="text-2xl mr-3">💬</span>
              <span className="font-medium">{notif.message}</span>
            </div>
            <button
              onClick={() => removeNotification(notif.id)}
              className="ml-4 text-white hover:text-gray-200"
            >
              ✕
            </button>
          </div>
        ))}
      </div>
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">Social Network</h1>
            <div className="flex items-center space-x-4">
              <span className="text-gray-700">
                Welcome, {user?.first_name} {user?.last_name}
              </span>
              
              {/* Notification Bell */}
              <button
                onClick={() => router.push('/notifications')}
                className="relative text-gray-600 hover:text-gray-900"
              >
                <span className="text-2xl">🔔</span>
                {unreadCount > 0 && (
                  <span className="absolute -top-1 -right-1 bg-red-600 text-white text-xs font-bold rounded-full h-5 w-5 flex items-center justify-center">
                    {unreadCount > 9 ? '9+' : unreadCount}
                  </span>
                )}
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
                
                {/* Image Preview */}
                {imagePreview && (
                  <div className="mt-4 relative">
                    <img 
                      src={imagePreview} 
                      alt="Preview" 
                      className="w-full max-h-64 object-cover rounded-md"
                    />
                    <button
                      type="button"
                      onClick={handleRemoveImage}
                      className="absolute top-2 right-2 bg-red-600 text-white rounded-full p-2 hover:bg-red-700"
                    >
                      ✕
                    </button>
                  </div>
                )}

                <div className="mt-4 flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <select
                      value={privacy}
                      onChange={(e) => setPrivacy(e.target.value)}
                      className="border border-gray-300 rounded-md px-3 py-2 focus:ring-blue-500 focus:border-blue-500"
                    >
                      <option value="public">Public</option>
                      <option value="almost_private">Followers Only</option>
                      <option value="private">Private</option>
                    </select>
                    
                    {/* Image Upload Button */}
                    <label className="cursor-pointer bg-gray-100 hover:bg-gray-200 text-gray-700 px-4 py-2 rounded-md flex items-center">
                      <span className="mr-2">📷</span>
                      <span>Image</span>
                      <input
                        type="file"
                        accept="image/jpeg,image/jpg,image/png,image/gif"
                        onChange={handleImageSelect}
                        className="hidden"
                      />
                    </label>
                  </div>
                  
                  <button
                    type="submit"
                    disabled={uploading}
                    className={`px-6 py-2 rounded-md ${
                      uploading 
                        ? 'bg-gray-400 cursor-not-allowed' 
                        : 'bg-blue-600 hover:bg-blue-700'
                    } text-white`}
                  >
                    {uploading ? 'Posting...' : 'Post'}
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
                      <img src={`http://localhost:8080${post.image_path}`} alt="Post image" className="w-full rounded-md mb-4" />
                    )}
                    <div className="flex items-center justify-between text-sm text-gray-500 mb-2">
                      <span className={`px-2 py-1 rounded-full text-xs ${
                        post.privacy === 'public' ? 'bg-green-100 text-green-800' :
                        post.privacy === 'almost_private' ? 'bg-yellow-100 text-yellow-800' :
                        'bg-red-100 text-red-800'
                      }`}>
                        {post.privacy === 'public' ? 'Public' :
                         post.privacy === 'almost_private' ? 'Followers Only' : 'Private'}
                      </span>
                      
                      {/* Post Like/Dislike Buttons */}
                      <div className="flex items-center gap-4">
                        <button
                          onClick={() => handlePostReaction(post.id, 'like')}
                          className={`flex items-center gap-1 ${
                            post.user_reaction === 'like' 
                              ? 'text-blue-600 font-semibold' 
                              : 'text-gray-600 hover:text-blue-600'
                          } transition`}
                        >
                          <span className="text-xl">👍</span>
                          <span>{post.likes || 0}</span>
                        </button>
                        <button
                          onClick={() => handlePostReaction(post.id, 'dislike')}
                          className={`flex items-center gap-1 ${
                            post.user_reaction === 'dislike' 
                              ? 'text-red-600 font-semibold' 
                              : 'text-gray-600 hover:text-red-600'
                          } transition`}
                        >
                          <span className="text-xl">👎</span>
                          <span>{post.dislikes || 0}</span>
                        </button>
                      </div>
                    </div>
                    
                    {/* Comments Component */}
                    <Comments postId={post.id} currentUserId={user?.id} />
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
                <button
                  onClick={() => router.push('/notifications')}
                  className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md flex items-center justify-between"
                >
                  <span>🔔 Notifications</span>
                  {unreadCount > 0 && (
                    <span className="bg-red-600 text-white text-xs font-bold px-2 py-1 rounded-full">
                      {unreadCount}
                    </span>
                  )}
                </button>
                <button
                  onClick={() => router.push('/groups')}
                  className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
                >
                  👥 Groups
                </button>
                <button
                  onClick={() => router.push('/messages')}
                  className="w-full text-left px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md"
                >
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
