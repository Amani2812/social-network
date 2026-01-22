'use client'

import { useEffect, useState, useRef } from 'react'
import { useParams, useRouter } from 'next/navigation'

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

interface GroupPost {
  id: number
  group_id: number
  user_id: number
  content: string
  image_path?: string
  created_at: string
  user?: User
}

interface Event {
  id: number
  group_id: number
  creator_id: number
  title: string
  description?: string
  event_time: string
  created_at: string
  user_response?: string
}

export default function GroupDetail() {
  const params = useParams()
  const router = useRouter()
  const [user, setUser] = useState<User | null>(null)
  const [group, setGroup] = useState<Group | null>(null)
  const [posts, setPosts] = useState<GroupPost[]>([])
  const [events, setEvents] = useState<Event[]>([])
  const [loading, setLoading] = useState(true)
  const [activeTab, setActiveTab] = useState<'posts' | 'events' | 'members'>('posts')
  interface JoinRequest {
    id: number
    group_id: number
    user_id: number
    status: string
    role: string
    created_at: string
    user?: User
  }

  const [joinRequests, setJoinRequests] = useState<JoinRequest[]>([])
  const [isAdmin, setIsAdmin] = useState(false)

  // Post creation
  const [newPost, setNewPost] = useState('')
  const [posting, setPosting] = useState(false)

  // Event creation
  const [showEventForm, setShowEventForm] = useState(false)
  const [newEventTitle, setNewEventTitle] = useState('')
  const [newEventDescription, setNewEventDescription] = useState('')
  const [newEventTime, setNewEventTime] = useState('')
  const [creatingEvent, setCreatingEvent] = useState(false)

  // WebSocket state
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected'>('connecting')
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const MAX_RECONNECT_ATTEMPTS = 10
  const hasConnectedRef = useRef(false)
  const isConnectingRef = useRef(false)

  // Emoji picker state
  const [showEmojiPicker, setShowEmojiPicker] = useState(false)

  // Keep refs updated
  useEffect(() => {
    wsRef.current = ws
  }, [ws])

  useEffect(() => {
    fetchUser()
    fetchGroup()
    fetchPosts()
    fetchEvents()
    fetchJoinRequests()

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
      if (wsRef.current) {
        wsRef.current.close()
        wsRef.current = null
      }
      hasConnectedRef.current = false
    }
  }, [params.id])

  useEffect(() => {
    // Only connect once when user is available
    if (user && !hasConnectedRef.current) {
      hasConnectedRef.current = true
      connectWebSocket()
    }
  }, [user])

  const connectWebSocket = () => {
    if (!user) return

    // Prevent multiple simultaneous connection attempts
    if (isConnectingRef.current) {
      console.log('⏳ Connection attempt already in progress')
      return
    }

    // Don't create a new connection if one already exists and is open
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      console.log('✅ WebSocket already connected')
      return
    }

    // Check reconnection attempts
    if (reconnectAttemptsRef.current >= MAX_RECONNECT_ATTEMPTS) {
      console.error('❌ Max reconnection attempts reached. Please refresh the page.')
      return
    }

    // Close existing connection if it's in a bad state
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }

    // Clear any existing reconnect timeout
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }

    isConnectingRef.current = true
    setConnectionStatus('connecting')
    console.log(`🔄 Connecting to WebSocket... (Attempt ${reconnectAttemptsRef.current + 1})`)

    const websocket = new WebSocket('ws://localhost:8080/ws')

    websocket.onopen = () => {
      console.log('✅ WebSocket Connected')
      setConnectionStatus('connected')
      reconnectAttemptsRef.current = 0 // Reset counter on successful connection
      isConnectingRef.current = false
    }

    websocket.onmessage = (event) => {
      // Handle multiple messages in one frame (separated by newlines)
      const messageStrings = event.data.trim().split('\n')

      messageStrings.forEach((messageStr: string) => {
        if (!messageStr.trim()) return // Skip empty messages

        try {
          const data = JSON.parse(messageStr)
          console.log('📥 Group page received:', data)

          if (data.type === 'group' && data.group_id === parseInt(params.id as string)) {
            // Only process messages from OTHER users (not from ourselves)
            // Our own messages are handled optimistically
            if (data.sender_id === user?.id) {
              console.log('⏭️ Skipping own group message (handled optimistically)')
              return
            }

            const newPost: GroupPost = {
              id: Date.now() + Math.random(),
              group_id: data.group_id,
              user_id: data.sender_id,
              content: data.content,
              created_at: data.timestamp || new Date().toISOString(),
            }

            console.log('📨 Processing group message from user:', newPost.user_id, 'in group:', newPost.group_id)

            // Add the post to the UI
            setPosts(prev => {
              // Check for duplicates in last 10 posts
              const recentPosts = prev.slice(-10)
              const isDuplicate = recentPosts.some(p =>
                p.content === newPost.content &&
                p.user_id === newPost.user_id &&
                Math.abs(new Date(p.created_at).getTime() - new Date(newPost.created_at).getTime()) < 3000
              )
              if (isDuplicate) {
                console.log('⚠️ Duplicate post detected, skipping')
                return prev
              }
              console.log('✅ Adding new group post from other user')
              return [...prev, newPost]
            })
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', messageStr, error)
        }
      })
    }

    websocket.onerror = (error) => {
      console.log('⚠️ WebSocket connection error (this is normal during reconnection):', error)
      setConnectionStatus('disconnected')
      isConnectingRef.current = false
    }

    websocket.onclose = (event) => {
      console.log(`🔌 WebSocket closed: ${event.code} - ${event.reason || 'No reason provided'}`)
      setConnectionStatus('disconnected')
      isConnectingRef.current = false
      wsRef.current = null
      setWs(null)

      // Attempt to reconnect with exponential backoff
      if (reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS && user) {
        reconnectAttemptsRef.current++
        const delay = Math.min(1000 * Math.pow(2, reconnectAttemptsRef.current - 1), 30000)
        console.log(`⏳ Reconnecting in ${delay / 1000} seconds...`)

        reconnectTimeoutRef.current = setTimeout(() => {
          connectWebSocket()
        }, delay)
      }
    }

    setWs(websocket)
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

  const fetchGroup = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/get?id=${params.id}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setGroup(data)
      }
    } catch (err) {
      console.error('Failed to fetch group:', err)
    } finally {
      setLoading(false)
    }
  }

  const fetchPosts = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/posts/get?group_id=${params.id}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setPosts(data || [])
      } else {
        setPosts([])
      }
    } catch (err) {
      console.error('Failed to fetch posts:', err)
      setPosts([])
    }
  }

  const fetchEvents = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/events/get?group_id=${params.id}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setEvents(data || [])
      } else {
        setEvents([])
      }
    } catch (err) {
      console.error('Failed to fetch events:', err)
      setEvents([])
    }
  }

  const fetchJoinRequests = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/groups/join/requests?group_id=${params.id}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setJoinRequests(data || [])
        // If we can fetch join requests successfully, we're an admin
        setIsAdmin(true)
      } else if (response.status === 403 || response.status === 401) {
        // Explicitly not an admin
        console.log('ℹ️ Not authorized to view join requests (not an admin)')
        setJoinRequests([])
        setIsAdmin(false)
      } else {
        // Other errors - check if user is group creator
        console.log('⚠️ Error fetching join requests, checking group ownership')
        setJoinRequests([])
        // Check if current user is the group creator
        if (group && user && group.creator_id === user.id) {
          console.log('✅ User is group creator, setting as admin')
          setIsAdmin(true)
        } else {
          setIsAdmin(false)
        }
      }
    } catch (err) {
      console.log('⚠️ Failed to fetch join requests:', err)
      setJoinRequests([])
      // Check if current user is the group creator as fallback
      if (group && user && group.creator_id === user.id) {
        console.log('✅ User is group creator (fallback check), setting as admin')
        setIsAdmin(true)
      } else {
        setIsAdmin(false)
      }
    }
  }

  const handleJoinRequest = async (userId: number, accept: boolean) => {
    try {
      const response = await fetch('http://localhost:8080/api/groups/join/respond', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          group_id: parseInt(params.id as string),
          user_id: userId,
          accept: accept,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        alert(accept ? 'Request accepted! ✓' : 'Request declined ✗')
        fetchJoinRequests() // Refresh the list
      } else {
        const errorData = await response.json()
        alert(`Failed: ${errorData.error || 'Unknown error'}`)
      }
    } catch (err) {
      console.error('Failed to respond to join request:', err)
      alert('Failed to respond to request')
    }
  }

  const handleCreatePost = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newPost.trim()) return

    setPosting(true)

    const messageContent = newPost.trim()
    const messageData = {
      type: 'group',
      group_id: parseInt(params.id as string),
      content: messageContent,
    }

    // Optimistic UI update - add post immediately
    const optimisticPost: GroupPost = {
      id: Date.now() + Math.random(),
      group_id: parseInt(params.id as string),
      user_id: user?.id || 0,
      content: messageContent,
      created_at: new Date().toISOString(),
      user: user || undefined,
    }

    setPosts(prev => [...prev, optimisticPost])
    setNewPost('')
    console.log('✅ Post added optimistically')

    // If WebSocket is connected, send immediately
    if (ws && ws.readyState === WebSocket.OPEN) {
      try {
        ws.send(JSON.stringify(messageData))
        console.log('📤 Post sent via WebSocket')
      } catch (error) {
        console.error('❌ Failed to send post:', error)
        // Remove optimistic post on failure
        setPosts(prev => prev.filter(p => p.id !== optimisticPost.id))
        setNewPost(messageContent) // Restore the message
      }
    } else {
      // Send via HTTP API as fallback
      console.log('⏳ WebSocket not connected, sending via HTTP')
      try {
        const response = await fetch('http://localhost:8080/api/groups/posts/create', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            group_id: parseInt(params.id as string),
            content: messageContent,
          }),
          credentials: 'include',
        })

        if (!response.ok) {
          const errorData = await response.json()
          alert(`Failed to create post: ${errorData.error || 'Unknown error'}`)
          // Remove optimistic post on failure
          setPosts(prev => prev.filter(p => p.id !== optimisticPost.id))
          setNewPost(messageContent) // Restore the message
        } else {
          alert('Post created successfully! ✓')
        }
      } catch (err) {
        console.error('Failed to create post:', err)
        alert('Failed to create post. You might not be a member of this group.')
        // Remove optimistic post on failure
        setPosts(prev => prev.filter(p => p.id !== optimisticPost.id))
        setNewPost(messageContent) // Restore the message
      }
    }

    setPosting(false)
  }

  const handleCreateEvent = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newEventTitle.trim() || !newEventTime) return

    setCreatingEvent(true)

    try {
      const response = await fetch('http://localhost:8080/api/events/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          group_id: parseInt(params.id as string),
          title: newEventTitle,
          description: newEventDescription || null,
          event_time: new Date(newEventTime).toISOString(),
        }),
        credentials: 'include',
      })

      if (response.ok) {
        setNewEventTitle('')
        setNewEventDescription('')
        setNewEventTime('')
        setShowEventForm(false)
        fetchEvents()
      } else {
        alert('Failed to create event')
      }
    } catch (err) {
      console.error('Failed to create event:', err)
      alert('Failed to create event')
    } finally {
      setCreatingEvent(false)
    }
  }

  const handleEventResponse = async (eventId: number, response: string) => {
    try {
      const res = await fetch('http://localhost:8080/api/events/respond', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          event_id: eventId,
          response: response,
        }),
        credentials: 'include',
      })

      if (res.ok) {
        alert(`You responded: ${response === 'going' ? 'Going ✓' : 'Not Going ✗'}`)
        // Refresh events to show updated response
        await fetchEvents()
      } else {
        const errorData = await res.json()
        alert(`Failed to respond: ${errorData.error || 'Unknown error'}`)
      }
    } catch (err) {
      console.error('Failed to respond to event:', err)
      alert('Failed to respond to event')
    }
  }

  const insertEmoji = (emoji: string) => {
    setNewPost(prev => prev + emoji)
    setShowEmojiPicker(false)
  }

  const emojis = ['😀', '😂', '❤️', '👍', '👎', '🔥', '🎉', '😢', '😮', '😍', '🤔', '😎', '🙄', '😴', '🤗', '🤩', '🥳', '😇', '🤖', '👻', '💀', '👽', '🎃', '🎄', '🎁', '💯', '✨', '🌟', '⭐', '🌙', '☀️', '🌈', '🌊', '🌸', '🌺', '🌻', '🌹', '🍀', '🌿', '🍎', '🍊', '🍋', '🍌', '🍉', '🍇', '🍓', '🍑', '🍒', '🥝', '🍅', '🥕', '🌽', '🥔', '🍞', '🥖', '🍕', '🍔', '🍟', '🌭', '🍿', '🍫', '🍬', '🍭', '🍰', '🎂', '🍪', '☕', '🍵', '🍺', '🍻', '🥂', '🍷', '🥃', '🍸', '🍹', '🧊', '🏆', '🎖️', '🏅', '🥇', '🥈', '🥉', '🏅', '🎗️', '🏵️', '🎀', '🎁', '🎈', '🎆', '🎇', '✨', '🎉', '🎊', '🎈', '🎁', '🏮', '🎋', '🎍', '🎎', '🎏', '🎐', '🎑', '🧧', '🎀', '🎁', '🎈', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁', '🎀', '🎉', '🎊', '🎈', '🎁',
