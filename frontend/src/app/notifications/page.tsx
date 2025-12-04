'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter } from 'next/navigation'

interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  avatar_path?: string
  nickname?: string
}

interface Notification {
  id: number
  user_id: number
  type: string
  content: string
  related_id?: number
  is_read: boolean
  created_at: string
}

export default function Notifications() {
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [loading, setLoading] = useState(true)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [unreadCount, setUnreadCount] = useState(0)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const isConnectingRef = useRef(false)
  const MAX_RECONNECT_ATTEMPTS = 10

  useEffect(() => {
    fetchUser()
    fetchNotifications()
    fetchUnreadCount()
    connectWebSocket()

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
  }, [])

  const connectWebSocket = () => {
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
        console.log('✅ Notifications WebSocket connected')
        reconnectAttemptsRef.current = 0 // Reset counter on successful connection
        isConnectingRef.current = false
      }

      websocket.onmessage = (event) => {
        const messages = event.data.trim().split('\n')

        messages.forEach((messageStr: string) => {
          if (!messageStr.trim()) return

          try {
            const data = JSON.parse(messageStr)
            console.log('📥 Notification received:', data)

            // Handle real-time notifications
            if (data.type === 'notification') {
              // Refresh notifications list
              fetchNotifications()
              fetchUnreadCount()
            }
          } catch (error) {
            console.error('Failed to parse WebSocket message:', messageStr, error)
          }
        })
      }

      websocket.onerror = (error) => {
        console.log('⚠️ WebSocket connection error (this is normal during reconnection):', error)
        isConnectingRef.current = false
      }

      websocket.onclose = (event) => {
        console.log(`🔌 WebSocket closed: ${event.code} - ${event.reason || 'No reason provided'}`)
        isConnectingRef.current = false
        wsRef.current = null
        setWs(null)

        // Attempt to reconnect with exponential backoff
        if (reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
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

  const fetchNotifications = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/notifications', {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setNotifications(data || [])
      } else {
        setNotifications([])
      }
    } catch (err) {
      console.error('Failed to fetch notifications:', err)
      setNotifications([])
    } finally {
      setLoading(false)
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

  const markAsRead = async (notificationId: number) => {
    try {
      const response = await fetch('http://localhost:8080/api/notifications/read', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ notification_id: notificationId }),
        credentials: 'include',
      })

      if (response.ok) {
        // Update local state
        setNotifications(notifications.map(n => 
          n.id === notificationId ? { ...n, is_read: true } : n
        ))
        setUnreadCount(Math.max(0, unreadCount - 1))
      }
    } catch (err) {
      console.error('Failed to mark notification as read:', err)
    }
  }

  const markAllAsRead = async () => {
    try {
      // Mark all unread notifications as read
      const unreadNotifications = notifications.filter(n => !n.is_read)
      for (const notification of unreadNotifications) {
        await fetch('http://localhost:8080/api/notifications/read', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ notification_id: notification.id }),
          credentials: 'include',
        })
      }

      // Update local state
      setNotifications(notifications.map(n => ({ ...n, is_read: true })))
      setUnreadCount(0)
    } catch (err) {
      console.error('Failed to mark all as read:', err)
    }
  }

  const handleAcceptFollowRequest = async (notification: Notification, e: React.MouseEvent) => {
    e.stopPropagation() // Prevent notification click
    
    if (!notification.related_id) return

    try {
      const response = await fetch('http://localhost:8080/api/follow/respond', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          follower_id: notification.related_id,
          accept: true,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        // Mark notification as read and refresh
        await markAsRead(notification.id)
        await fetchNotifications()
        await fetchUnreadCount()
        alert('Follow request accepted! ✓')
      } else {
        const errorData = await response.json()
        alert(`Failed to accept: ${errorData.error || 'Unknown error'}`)
      }
    } catch (err) {
      console.error('Failed to accept follow request:', err)
      alert('Failed to accept follow request')
    }
  }

  const handleDeclineFollowRequest = async (notification: Notification, e: React.MouseEvent) => {
    e.stopPropagation() // Prevent notification click
    
    if (!notification.related_id) return

    try {
      const response = await fetch('http://localhost:8080/api/follow/respond', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          follower_id: notification.related_id,
          accept: false,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        // Mark notification as read and refresh
        await markAsRead(notification.id)
        await fetchNotifications()
        await fetchUnreadCount()
        alert('Follow request declined ✗')
      } else {
        const errorData = await response.json()
        alert(`Failed to decline: ${errorData.error || 'Unknown error'}`)
      }
    } catch (err) {
      console.error('Failed to decline follow request:', err)
      alert('Failed to decline follow request')
    }
  }

  const handleNotificationClick = (notification: Notification) => {
    // Don't navigate for follow requests (they have action buttons)
    if (notification.type === 'follow_request') {
      return
    }

    // Mark as read
    if (!notification.is_read) {
      markAsRead(notification.id)
    }

    // Navigate based on notification type
    if (notification.type === 'group_invite' && notification.related_id) {
      router.push(`/groups/${notification.related_id}`)
    } else if (notification.type === 'event_invite' && notification.related_id) {
      router.push(`/events/${notification.related_id}`)
    } else if (notification.type === 'message' && notification.related_id) {
      router.push(`/messages?user=${notification.related_id}`)
    }
  }

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'follow_request':
        return '👤'
      case 'group_invite':
        return '👥'
      case 'event_invite':
        return '📅'
      case 'new_post':
        return '📝'
      case 'message':
        return '💬'
      default:
        return '🔔'
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
            <h1 className="text-2xl font-bold text-gray-900">Notifications</h1>
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
        <div className="bg-white shadow rounded-lg">
          {/* Notifications Header */}
          <div className="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
            <div className="flex items-center">
              <h2 className="text-xl font-semibold text-gray-900">
                All Notifications
              </h2>
              {unreadCount > 0 && (
                <span className="ml-3 bg-blue-600 text-white text-xs font-bold px-2 py-1 rounded-full">
                  {unreadCount} new
                </span>
              )}
            </div>
            {unreadCount > 0 && (
              <button
                onClick={markAllAsRead}
                className="text-sm text-blue-600 hover:text-blue-800"
              >
                Mark all as read
              </button>
            )}
          </div>

          {/* Notifications List */}
          <div className="divide-y divide-gray-200">
            {notifications.length === 0 ? (
              <div className="px-6 py-12 text-center">
                <div className="text-gray-400 text-6xl mb-4">🔔</div>
                <p className="text-gray-500 text-lg">No notifications yet</p>
                <p className="text-gray-400 text-sm mt-2">
                  You'll see notifications here when someone follows you, invites you to a group, or more!
                </p>
              </div>
            ) : (
              notifications.map((notification) => (
                <div
                  key={notification.id}
                  onClick={() => handleNotificationClick(notification)}
                  className={`px-6 py-4 transition-colors ${
                    notification.type === 'follow_request' ? '' : 'cursor-pointer'
                  } ${
                    notification.is_read
                      ? 'bg-white hover:bg-gray-50'
                      : 'bg-blue-50 hover:bg-blue-100'
                  }`}
                >
                  <div className="flex items-start">
                    <div className="flex-shrink-0 text-3xl mr-4">
                      {getNotificationIcon(notification.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <p className={`text-sm ${
                        notification.is_read ? 'text-gray-700' : 'text-gray-900 font-semibold'
                      }`}>
                        {notification.content}
                      </p>
                      <p className="text-xs text-gray-500 mt-1">
                        {new Date(notification.created_at).toLocaleString()}
                      </p>
                      
                      {/* Accept/Decline buttons for follow requests */}
                      {notification.type === 'follow_request' && !notification.is_read && (
                        <div className="flex space-x-2 mt-3">
                          <button
                            onClick={(e) => handleAcceptFollowRequest(notification, e)}
                            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors"
                          >
                            ✓ Accept
                          </button>
                          <button
                            onClick={(e) => handleDeclineFollowRequest(notification, e)}
                            className="bg-gray-300 hover:bg-gray-400 text-gray-700 px-4 py-2 rounded-md text-sm font-medium transition-colors"
                          >
                            ✗ Decline
                          </button>
                        </div>
                      )}
                    </div>
                    {!notification.is_read && (
                      <div className="flex-shrink-0 ml-4">
                        <div className="w-2 h-2 bg-blue-600 rounded-full"></div>
                      </div>
                    )}
                  </div>
                </div>
              ))
            )}
          </div>
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
