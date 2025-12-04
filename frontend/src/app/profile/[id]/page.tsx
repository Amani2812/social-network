'use client'

import { useEffect, useState, useRef } from 'react'
import { useParams, useRouter } from 'next/navigation'

interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  date_of_birth: string
  avatar_path?: string
  nickname?: string
  about_me?: string
  is_public: boolean
  created_at: string
  updated_at: string
}

interface Post {
  id: number
  user_id: number
  content: string
  image_path?: string
  privacy: string
  created_at: string
}

interface Follow {
  id: number
  follower_id: number
  following_id: number
  status: string
  created_at: string
  updated_at: string
}

export default function Profile() {
  const params = useParams()
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [posts, setPosts] = useState<Post[]>([])
  const [followers, setFollowers] = useState<User[]>([])
  const [following, setFollowing] = useState<User[]>([])
  const [isFollowing, setIsFollowing] = useState(false)
  const [hasPendingRequest, setHasPendingRequest] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showFollowersModal, setShowFollowersModal] = useState(false)
  const [showFollowingModal, setShowFollowingModal] = useState(false)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const isConnectingRef = useRef(false)
  const MAX_RECONNECT_ATTEMPTS = 10

  useEffect(() => {
    fetchProfile()
    fetchCurrentUser()
  }, [params.id])

  useEffect(() => {
    if (currentUser) {
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
  }, [currentUser])

  const connectWebSocket = () => {
    if (!currentUser) return
    
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
        console.log('✅ Profile WebSocket connected')
        reconnectAttemptsRef.current = 0 // Reset counter on successful connection
        isConnectingRef.current = false
      }
      
      websocket.onmessage = (event) => {
        const messages = event.data.trim().split('\n')

        messages.forEach((messageStr: string) => {
          if (!messageStr.trim()) return

          try {
            const data = JSON.parse(messageStr)
            console.log('📥 Profile received:', data)

            // Handle follow status updates
            if (data.type === 'follow_accepted' || data.type === 'follow_status_update') {
              // Refresh follow data to show updated status
              fetchFollowData()
            }
          } catch (error) {
            console.error('Failed to parse WebSocket message:', messageStr, error)
          }
        })
      }
      
      websocket.onerror = (error) => {
        console.error('❌ WebSocket error:', error)
        isConnectingRef.current = false
      }
      
      websocket.onclose = (event) => {
        console.log(`🔌 WebSocket closed: ${event.code} - ${event.reason || 'No reason provided'}`)
        isConnectingRef.current = false
        wsRef.current = null
        setWs(null)
        
        // Attempt to reconnect with exponential backoff
        if (currentUser && reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
          reconnectAttemptsRef.current++
          const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current - 1), 30000)
          console.log(`⏳ Reconnecting in ${delay / 1000} seconds...`)
          
          reconnectTimeoutRef.current = setTimeout(() => {
            connectWebSocket()
          }, delay)
        }
      }
      
      wsRef.current = websocket
      setWs(websocket)
    } catch (error) {
      console.error('❌ Failed to create WebSocket connection:', error)
      isConnectingRef.current = false
    }
  }

  const fetchProfile = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/users/${params.id}`)
      if (response.ok) {
        const userData = await response.json()
        setUser(userData)
        fetchUserPosts(userData.id)
      } else {
        setError('Profile not found')
      }
    } catch (err) {
      setError('Failed to load profile')
    }
  }

  const fetchUserPosts = async (userId: number) => {
    try {
      const response = await fetch(`http://localhost:8080/api/posts/user?user_id=${userId}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const postsData = await response.json()
        setPosts(postsData || [])
      } else {
        setPosts([])
      }
    } catch (err) {
      console.error('Failed to fetch posts:', err)
      setPosts([])
    }
  }

  const fetchCurrentUser = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/auth/me', {
        credentials: 'include',
      })
      if (response.ok) {
        const userData = await response.json()
        setCurrentUser(userData)
      }
    } catch (err) {
      // User not logged in
    }
  }

  const fetchFollowData = async () => {
    if (!user || !currentUser) return

    try {
      const [followersRes, followingRes, followStatusRes] = await Promise.all([
        fetch(`http://localhost:8080/api/follow/followers?user_id=${user.id}`),
        fetch(`http://localhost:8080/api/follow/following?user_id=${user.id}`),
        fetch(`http://localhost:8080/api/follow/status?follower_id=${currentUser.id}&following_id=${user.id}`)
      ])

      if (followersRes.ok) {
        const followersData = await followersRes.json()
        setFollowers(followersData || [])
      } else {
        setFollowers([])
      }

      if (followingRes.ok) {
        const followingData = await followingRes.json()
        setFollowing(followingData || [])
      } else {
        setFollowing([])
      }

      if (followStatusRes.ok) {
        const statusData = await followStatusRes.json()
        setIsFollowing(statusData.is_following)
        setHasPendingRequest(statusData.has_pending_request)
      }
    } catch (err) {
      console.error('Failed to fetch follow data:', err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (user && currentUser) {
      fetchFollowData()
    }
  }, [user, currentUser])

  const handleFollow = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/follow/request', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ following_id: user!.id }),
        credentials: 'include',
      })

      if (response.ok) {
        // Refresh follow status after sending request
        await fetchFollowData()
      }
    } catch (err) {
      console.error('Failed to send follow request:', err)
    }
  }

  const handleUnfollow = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/follow/unfollow?following_id=${user!.id}`, {
        method: 'POST',
        credentials: 'include',
      })

      if (response.ok) {
        // Refresh follow status after unfollowing
        await fetchFollowData()
      }
    } catch (err) {
      console.error('Failed to unfollow:', err)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    )
  }

  if (error || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-red-600">{error || 'Profile not found'}</div>
      </div>
    )
  }

  const isOwnProfile = currentUser && currentUser.id === user.id
  const canViewProfile = user.is_public || isFollowing || isOwnProfile

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="bg-white shadow rounded-lg overflow-hidden">
          {/* Profile Header */}
          <div className="bg-gradient-to-r from-blue-500 to-purple-600 h-32"></div>
          <div className="relative px-6 pb-6">
            <div className="flex flex-col sm:flex-row sm:items-end sm:space-x-5">
              <div className="-mt-16">
                <div className="w-32 h-32 bg-gray-300 rounded-full border-4 border-white flex items-center justify-center">
                  {user.avatar_path ? (
                    <img src={user.avatar_path} alt="Avatar" className="w-full h-full rounded-full object-cover" />
                  ) : (
                    <span className="text-4xl text-gray-600">{user.first_name[0]}{user.last_name[0]}</span>
                  )}
                </div>
              </div>
              <div className="mt-6 sm:mt-0 sm:flex-1">
                <h1 className="text-3xl font-bold text-gray-900">
                  {user.first_name} {user.last_name}
                  {user.nickname && <span className="text-xl text-gray-500 ml-2">({user.nickname})</span>}
                </h1>
                <p className="text-gray-600">{user.email}</p>
                <div className="mt-4 flex items-center space-x-4">
                  <button
                    onClick={() => setShowFollowersModal(true)}
                    className="text-sm text-gray-500 hover:text-blue-600 hover:underline"
                  >
                    <span className="font-semibold text-gray-900">{followers.length}</span> followers
                  </button>
                  <button
                    onClick={() => setShowFollowingModal(true)}
                    className="text-sm text-gray-500 hover:text-blue-600 hover:underline"
                  >
                    <span className="font-semibold text-gray-900">{following.length}</span> following
                  </button>
                  <span className={`text-sm px-2 py-1 rounded-full ${user.is_public ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                    {user.is_public ? 'Public' : 'Private'}
                  </span>
                </div>
              </div>
              <div className="mt-6 sm:mt-0">
                {isOwnProfile ? (
                  <button 
                    onClick={() => router.push('/profile/edit')}
                    className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
                  >
                    Edit Profile
                  </button>
                ) : currentUser && (
                  <div className="flex gap-2">
                    {isFollowing ? (
                      <>
                        <button
                          onClick={handleUnfollow}
                          className="bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-md"
                        >
                          Unfollow
                        </button>
                        <button
                          onClick={() => router.push(`/messages?user=${user.id}`)}
                          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
                        >
                          💬 Message
                        </button>
                      </>
                    ) : hasPendingRequest ? (
                      <button className="bg-yellow-600 text-white px-4 py-2 rounded-md cursor-not-allowed" disabled>
                        Request Pending
                      </button>
                    ) : (
                      <button
                        onClick={handleFollow}
                        className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
                      >
                        Follow
                      </button>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Profile Content */}
          {canViewProfile ? (
            <div className="px-6 py-6">
              {user.about_me && (
                <div className="mb-6">
                  <h2 className="text-xl font-semibold text-gray-900 mb-2">About</h2>
                  <p className="text-gray-700">{user.about_me}</p>
                </div>
              )}

              <div className="mb-6">
                <h2 className="text-xl font-semibold text-gray-900 mb-4">Posts</h2>
                {posts.length === 0 ? (
                  <div className="text-center py-8 text-gray-500">
                    No posts yet
                  </div>
                ) : (
                  <div className="space-y-4">
                    {posts.map((post) => (
                      <div key={post.id} className="bg-gray-50 rounded-lg p-4">
                        <p className="text-gray-700 mb-2">{post.content}</p>
                        {post.image_path && (
                          <img src={post.image_path} alt="Post image" className="w-full rounded-md mb-2" />
                        )}
                        <div className="flex items-center justify-between text-sm text-gray-500">
                          <span>{new Date(post.created_at).toLocaleDateString()}</span>
                          <span className={`px-2 py-1 rounded-full text-xs ${
                            post.privacy === 'public' ? 'bg-green-100 text-green-800' :
                            post.privacy === 'almost_private' ? 'bg-yellow-100 text-yellow-800' :
                            'bg-red-100 text-red-800'
                          }`}>
                            {post.privacy === 'public' ? 'Public' :
                             post.privacy === 'almost_private' ? 'Followers Only' : 'Private'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ) : (
            <div className="px-6 py-12 text-center">
              <div className="text-gray-500 text-lg">
                This profile is private. Follow this user to see their posts.
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Followers Modal */}
      {showFollowersModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={() => setShowFollowersModal(false)}>
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 max-h-96 overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">Followers</h2>
              <button onClick={() => setShowFollowersModal(false)} className="text-gray-500 hover:text-gray-700">
                ✕
              </button>
            </div>
            {followers.length === 0 ? (
              <p className="text-center text-gray-500 py-4">No followers yet</p>
            ) : (
              <div className="space-y-3">
                {followers.map((follower) => (
                  <button
                    key={follower.id}
                    onClick={() => {
                      setShowFollowersModal(false)
                      router.push(`/profile/${follower.id}`)
                    }}
                    className="w-full flex items-center p-3 hover:bg-gray-50 rounded-lg"
                  >
                    <div className="w-12 h-12 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                      {follower.avatar_path ? (
                        <img src={follower.avatar_path} alt="Avatar" className="w-full h-full rounded-full object-cover" />
                      ) : (
                        <span className="text-lg text-gray-600">
                          {follower.first_name[0]}{follower.last_name[0]}
                        </span>
                      )}
                    </div>
                    <div className="flex-1 text-left">
                      <p className="font-semibold text-gray-900">
                        {follower.first_name} {follower.last_name}
                      </p>
                      {follower.nickname && (
                        <p className="text-sm text-gray-500">@{follower.nickname}</p>
                      )}
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Following Modal */}
      {showFollowingModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" onClick={() => setShowFollowingModal(false)}>
          <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4 max-h-96 overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-xl font-bold text-gray-900">Following</h2>
              <button onClick={() => setShowFollowingModal(false)} className="text-gray-500 hover:text-gray-700">
                ✕
              </button>
            </div>
            {following.length === 0 ? (
              <p className="text-center text-gray-500 py-4">Not following anyone yet</p>
            ) : (
              <div className="space-y-3">
                {following.map((followedUser) => (
                  <button
                    key={followedUser.id}
                    onClick={() => {
                      setShowFollowingModal(false)
                      router.push(`/profile/${followedUser.id}`)
                    }}
                    className="w-full flex items-center p-3 hover:bg-gray-50 rounded-lg"
                  >
                    <div className="w-12 h-12 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                      {followedUser.avatar_path ? (
                        <img src={followedUser.avatar_path} alt="Avatar" className="w-full h-full rounded-full object-cover" />
                      ) : (
                        <span className="text-lg text-gray-600">
                          {followedUser.first_name[0]}{followedUser.last_name[0]}
                        </span>
                      )}
                    </div>
                    <div className="flex-1 text-left">
                      <p className="font-semibold text-gray-900">
                        {followedUser.first_name} {followedUser.last_name}
                      </p>
                      {followedUser.nickname && (
                        <p className="text-sm text-gray-500">@{followedUser.nickname}</p>
                      )}
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
