'use client'

import { useEffect, useState, useRef } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'

interface User {
  id: number
  first_name: string
  last_name: string
  avatar_path?: string
  nickname?: string
}

interface Message {
  id: number
  sender_id: number
  receiver_id?: number
  group_id?: number
  content: string
  created_at: string
  sender?: User
}

export default function Messages() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const selectedUserId = searchParams.get('user')
  
  const [currentUser, setCurrentUser] = useState<User | null>(null)
  const [conversations, setConversations] = useState<User[]>([])
  const [selectedUser, setSelectedUser] = useState<User | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [newMessage, setNewMessage] = useState('')
  const [loading, setLoading] = useState(true)
  const [ws, setWs] = useState<WebSocket | null>(null)
  const [connectionStatus, setConnectionStatus] = useState<'connecting' | 'connected' | 'disconnected'>('connecting')
  const [showEmojiPicker, setShowEmojiPicker] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const selectedUserRef = useRef<User | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const messageQueueRef = useRef<any[]>([])
  const hasConnectedRef = useRef(false)
  const isConnectingRef = useRef(false)
  const MAX_RECONNECT_ATTEMPTS = 10

  // Keep refs updated
  useEffect(() => {
    wsRef.current = ws
  }, [ws])

  useEffect(() => {
    selectedUserRef.current = selectedUser
  }, [selectedUser])

  useEffect(() => {
    fetchCurrentUser()
    
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
  }, [])

  useEffect(() => {
    // Only connect once when currentUser is available
    if (currentUser && !hasConnectedRef.current) {
      hasConnectedRef.current = true
      connectWebSocket()
    }
  }, [currentUser])

  useEffect(() => {
    if (selectedUserId && currentUser) {
      selectUserById(parseInt(selectedUserId))
    }
  }, [selectedUserId, currentUser])

  useEffect(() => {
    if (selectedUser) {
      fetchMessages(selectedUser.id)
    }
  }, [selectedUser])

  useEffect(() => {
    scrollToBottom()
  }, [messages])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }

  const connectWebSocket = () => {
    if (!currentUser) return
    
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
      
      // Send any queued messages
      if (messageQueueRef.current.length > 0) {
        console.log(`📤 Sending ${messageQueueRef.current.length} queued messages`)
        messageQueueRef.current.forEach(msg => {
          websocket.send(JSON.stringify(msg))
        })
        messageQueueRef.current = []
      }
    }

    // Handle ping from server - browser automatically responds with pong
    websocket.addEventListener('ping', () => {
      console.log('🏓 Received ping from server')
    })
    
    websocket.onmessage = (event) => {
      // Handle multiple messages in one frame (separated by newlines)
      const messageStrings = event.data.trim().split('\n')

      messageStrings.forEach((messageStr: string) => {
        if (!messageStr.trim()) return // Skip empty messages

        try {
          const data = JSON.parse(messageStr)
          console.log('📥 Messages page received:', data)

          if (data.type === 'private') {
            // Only process messages from OTHER users (not from ourselves)
            // Our own messages are handled optimistically
            if (data.sender_id === currentUser.id) {
              console.log('⏭️ Skipping own message (handled optimistically)')
              return
            }

            const newMsg: Message = {
              id: Date.now() + Math.random(),
              sender_id: data.sender_id,
              receiver_id: data.receiver_id,
              content: data.content,
              created_at: data.timestamp || new Date().toISOString(),
            }

            console.log('📨 Processing message from user:', newMsg.sender_id, 'to:', newMsg.receiver_id)
            console.log('📋 Current selected user:', selectedUserRef.current?.id)

            // Add message if it's relevant to current conversation
            // Check if we're in a conversation with the sender
            if (selectedUserRef.current && selectedUserRef.current.id === newMsg.sender_id) {
              console.log('✅ Message is for current conversation, adding...')
              setMessages(prev => {
                // Check for duplicates in last 10 messages
                const recentMessages = prev.slice(-10)
                const isDuplicate = recentMessages.some(m => 
                  m.content === newMsg.content && 
                  m.sender_id === newMsg.sender_id &&
                  Math.abs(new Date(m.created_at).getTime() - new Date(newMsg.created_at).getTime()) < 3000
                )
                if (isDuplicate) {
                  console.log('⚠️ Duplicate message detected, skipping')
                  return prev
                }
                console.log('✅ Adding new message from other user')
                return [...prev, newMsg]
              })
            } else {
              console.log('ℹ️ Message not for current conversation, skipping UI update')
            }

            // Add sender to conversations if new
            setConversations(prev => {
              if (!prev.some(u => u.id === newMsg.sender_id)) {
                fetchUserById(newMsg.sender_id).then(user => {
                  if (user) {
                    setConversations(p => [user, ...p])
                  }
                })
              }
              return prev
            })
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', messageStr, error)
        }
      })
    }
    
    websocket.onerror = (error) => {
      console.error('❌ WebSocket connection error:', error)
      console.log('💡 This may be due to authentication issues or server unavailability')
      setConnectionStatus('disconnected')
      isConnectingRef.current = false
    }
    
    websocket.onclose = (event) => {
      console.log(`🔌 WebSocket closed: Code ${event.code} - ${event.reason || 'No reason provided'}`)
      
      // Check if it's an authentication error (code 1006 or 1008)
      if (event.code === 1006 || event.code === 1008) {
        console.warn('⚠️ WebSocket closed due to possible authentication issue')
        console.log('💡 Make sure you are logged in and have a valid session')
      }
      
      setConnectionStatus('disconnected')
      isConnectingRef.current = false
      wsRef.current = null
      setWs(null)
      
      // Only attempt to reconnect if user is still authenticated and we haven't exceeded max attempts
      if (currentUser && reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
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
  }

  const fetchCurrentUser = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/auth/me', {
        credentials: 'include',
      })
      if (response.ok) {
        const userData = await response.json()
        setCurrentUser(userData)
        fetchConversations(userData.id)
      } else {
        router.push('/login')
      }
    } catch (err) {
      router.push('/login')
    } finally {
      setLoading(false)
    }
  }

  const fetchUserById = async (userId: number): Promise<User | null> => {
    try {
      const response = await fetch(`http://localhost:8080/api/users/${userId}`)
      if (response.ok) {
        return await response.json()
      }
      return null
    } catch (err) {
      return null
    }
  }
  
  const selectUserById = async (userId: number) => {
    const userData = await fetchUserById(userId)
    if (userData) {
      setSelectedUser(userData)
    }
  }

  const fetchConversations = async (userId: number) => {
    try {
      // Get all users who have messaged with this user
      const response = await fetch(`http://localhost:8080/api/messages/conversations`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setConversations(data || [])
      } else {
        // Fallback to followers if conversations endpoint doesn't exist
        const fallbackResponse = await fetch(`http://localhost:8080/api/follow/following?user_id=${userId}`, {
          credentials: 'include',
        })
        if (fallbackResponse.ok) {
          const data = await fallbackResponse.json()
          setConversations(data || [])
        }
      }
    } catch (err) {
      console.error('Failed to fetch conversations')
    }
  }

  const fetchMessages = async (userId: number) => {
    try {
      const response = await fetch(`http://localhost:8080/api/messages/private?user_id=${userId}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setMessages(data || [])
      }
    } catch (err) {
      setMessages([])
    }
  }

  const emojis = [
    '😀', '😃', '😄', '😁', '😆', '😅', '🤣', '😂', '🙂', '🙃', '😉', '😊', '😇',
    '🥰', '😍', '🤩', '😘', '😗', '😚', '😙', '🥲', '😋', '😛', '😜', '🤪', '😝',
    '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐', '😑', '😶', '😏', '😒', '🙄',
    '😬', '🤥', '😌', '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢', '🤮', '🤧',
    '🥵', '🥶', '🥴', '😵', '🤯', '🤠', '🥳', '🥸', '😎', '🤓', '🧐', '😕', '😟',
    '🙁', '☹️', '😮', '😯', '😲', '😳', '🥺', '😦', '😧', '😨', '😰', '😥', '😢',
    '😭', '😱', '😖', '😣', '😞', '😓', '😩', '😫', '🥱', '😤', '😡', '😠', '🤬',
    '😈', '👿', '💀', '☠️', '💩', '🤡', '👹', '👺', '👻', '👽', '👾', '🤖', '😺',
    '😸', '😹', '😻', '😼', '😽', '🙀', '😿', '😾', '❤️', '🧡', '💛', '💚', '💙',
    '💜', '🤎', '🖤', '🤍', '💔', '❤️‍🔥', '❤️‍🩹', '💕', '💞', '💓', '💗', '💖', '💘',
    '💝', '💟', '☮️', '✝️', '☪️', '🕉️', '☸️', '✡️', '🔯', '🕎', '☯️', '☦️', '🛐',
    '⛎', '♈', '♉', '♊', '♋', '♌', '♍', '♎', '♏', '♐', '♑', '♒', '♓', '🆔',
    '⚛️', '🉑', '☢️', '☣️', '📴', '📳', '🈶', '🈚', '🈸', '🈺', '🈷️', '✴️', '🆚',
    '💮', '🉐', '㊙️', '㊗️', '🈴', '🈵', '🈹', '🈲', '🅰️', '🅱️', '🆎', '🆑', '🅾️',
    '🆘', '❌', '⭕', '🛑', '⛔', '📛', '🚫', '💯', '💢', '♨️', '🚷', '🚯', '🚳',
    '🚱', '🔞', '📵', '🚭', '❗', '❕', '❓', '❔', '‼️', '⁉️', '🔅', '🔆', '〽️',
    '⚠️', '🚸', '🔱', '⚜️', '🔰', '♻️', '✅', '🈯', '💹', '❇️', '✳️', '❎', '🌐',
    '💠', 'Ⓜ️', '🌀', '💤', '🏧', '🚾', '♿', '🅿️', '🈳', '🈂️', '🛂', '🛃', '🛄',
    '🛅', '🚹', '🚺', '🚼', '⚧️', '🚻', '🚮', '🎦', '📶', '🈁', '🔣', 'ℹ️', '🔤',
    '🔡', '🔠', '🆖', '🆗', '🆙', '🆒', '🆕', '🆓', '0️⃣', '1️⃣', '2️⃣', '3️⃣', '4️⃣',
    '5️⃣', '6️⃣', '7️⃣', '8️⃣', '9️⃣', '🔟', '🔢', '#️⃣', '*️⃣', '⏏️', '▶️', '⏸️', '⏯️',
    '⏹️', '⏺️', '⏭️', '⏮️', '⏩', '⏪', '⏫', '⏬', '◀️', '🔼', '🔽', '➡️', '⬅️',
    '⬆️', '⬇️', '↗️', '↘️', '↙️', '↖️', '↕️', '↔️', '↪️', '↩️', '⤴️', '⤵️', '🔀',
    '🔁', '🔂', '🔄', '🔃', '🎵', '🎶', '➕', '➖', '➗', '✖️', '♾️', '💲', '💱',
    '™️', '©️', '®️', '〰️', '➰', '➿', '🔚', '🔙', '🔛', '🔝', '🔜', '✔️', '☑️',
    '🔘', '🔴', '🟠', '🟡', '🟢', '🔵', '🟣', '⚫', '⚪', '🟤', '🔺', '🔻', '🔸',
    '🔹', '🔶', '🔷', '🔳', '🔲', '▪️', '▫️', '◾', '◽', '◼️', '◻️', '🟥', '🟧',
    '🟨', '🟩', '🟦', '🟪', '⬛', '⬜', '🟫', '🔈', '🔇', '🔉', '🔊', '🔔', '🔕',
    '📣', '📢', '👁️‍🗨️', '💬', '💭', '🗯️', '♠️', '♣️', '♥️', '♦️', '🃏', '🎴', '🀄',
    '🕐', '🕑', '🕒', '🕓', '🕔', '🕕', '🕖', '🕗', '🕘', '🕙', '🕚', '🕛', '🕜',
    '🕝', '🕞', '🕟', '🕠', '🕡', '🕢', '🕣', '🕤', '🕥', '🕦', '🕧', '👍', '👎',
    '👌', '✌️', '🤞', '🤟', '🤘', '🤙', '👈', '👉', '👆', '👇', '☝️', '✋', '🤚',
    '🖐️', '🖖', '👋', '🤙', '💪', '🦾', '🖕', '✍️', '🙏', '🦶', '🦵', '🦿', '💄',
    '💋', '👄', '🦷', '👅', '👂', '🦻', '👃', '👣', '👁️', '👀', '🧠', '🫀', '🫁',
    '🦴', '👤', '👥', '🗣️', '👶', '👧', '🧒', '👦', '👩', '🧑', '👨', '👩‍🦱', '🧑‍🦱',
    '👨‍🦱', '👩‍🦰', '🧑‍🦰', '👨‍🦰', '👱‍♀️', '👱', '👱‍♂️', '👩‍🦳', '🧑‍🦳', '👨‍🦳', '👩‍🦲', '🧑‍🦲',
    '👨‍🦲', '🧔', '👵', '🧓', '👴', '👲', '👳‍♀️', '👳', '👳‍♂️', '🧕', '👮‍♀️', '👮', '👮‍♂️',
    '👷‍♀️', '👷', '👷‍♂️', '💂‍♀️', '💂', '💂‍♂️', '🕵️‍♀️', '🕵️', '🕵️‍♂️', '👩‍⚕️', '🧑‍⚕️', '👨‍⚕️',
    '👩‍🌾', '🧑‍🌾', '👨‍🌾', '👩‍🍳', '🧑‍🍳', '👨‍🍳', '👩‍🎓', '🧑‍🎓', '👨‍🎓', '👩‍🎤', '🧑‍🎤', '👨‍🎤',
    '👩‍🏫', '🧑‍🏫', '👨‍🏫', '👩‍🏭', '🧑‍🏭', '👨‍🏭', '👩‍💻', '🧑‍💻', '👨‍💻', '👩‍💼', '🧑‍💼', '👨‍💼',
    '🎃', '🎄', '🎆', '🎇', '🧨', '✨', '🎈', '🎉', '🎊', '🎋', '🎍', '🎎', '🎏',
    '🎐', '🎑', '🧧', '🎀', '🎁', '🎗️', '🎟️', '🎫', '🎖️', '🏆', '🏅', '🥇', '🥈',
    '🥉', '⚽', '⚾', '🥎', '🏀', '🏐', '🏈', '🏉', '🎾', '🥏', '🎳', '🏏', '🏑',
    '🏒', '🥍', '🏓', '🏸', '🥊', '🥋', '🥅', '⛳', '⛸️', '🎣', '🤿', '🎽', '🎿',
    '🛷', '🥌', '🎯', '🪀', '🪁', '🎱', '🔮', '🪄', '🧿', '🎮', '🕹️', '🎰', '🎲',
    '🧩', '🧸', '🪅', '🪆', '♟️', '🎭', '🎨', '🧵', '🪡', '🧶', '🪢'
  ]

  const addEmoji = (emoji: string) => {
    setNewMessage(prev => prev + emoji)
    setShowEmojiPicker(false)
  }

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!newMessage.trim() || !selectedUser || !currentUser) {
      return
    }

    const messageContent = newMessage.trim()
    const messageData = {
      type: 'private',
      receiver_id: selectedUser.id,
      content: messageContent,
    }

    // Optimistic UI update - add message immediately
    const optimisticMessage: Message = {
      id: Date.now() + Math.random(),
      sender_id: currentUser.id,
      receiver_id: selectedUser.id,
      content: messageContent,
      created_at: new Date().toISOString(),
    }

    setMessages(prev => [...prev, optimisticMessage])
    setNewMessage('')
    console.log('✅ Message added optimistically')

    // If WebSocket is connected, send immediately
    if (ws && ws.readyState === WebSocket.OPEN) {
      try {
        ws.send(JSON.stringify(messageData))
        console.log('📤 Message sent via WebSocket')
      } catch (error) {
        console.error('❌ Failed to send message:', error)
        // Queue message for retry
        messageQueueRef.current.push(messageData)
      }
    } else {
      // Queue message if not connected
      console.log('⏳ WebSocket not connected, queueing message')
      messageQueueRef.current.push(messageData)
      
      // Try to reconnect
      if (!ws || ws.readyState === WebSocket.CLOSED) {
        connectWebSocket()
      }
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-xl">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white">
      {/* Header - Instagram style */}
      <header className="border-b border-gray-300 bg-white sticky top-0 z-10">
        <div className="max-w-6xl mx-auto px-4">
          <div className="flex justify-between items-center h-16">
            <button
              onClick={() => router.push('/dashboard')}
              className="text-2xl font-semibold"
            >
              ← Back
            </button>
            <div className="flex items-center gap-3">
              <h1 className="text-xl font-semibold">{currentUser?.first_name}</h1>
              {/* Connection Status Indicator */}
              <div className="flex items-center gap-2">
                {connectionStatus === 'connected' && (
                  <div className="flex items-center gap-1">
                    <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                    <span className="text-xs text-green-600">Connected</span>
                  </div>
                )}
                {connectionStatus === 'connecting' && (
                  <div className="flex items-center gap-1">
                    <div className="w-2 h-2 bg-yellow-500 rounded-full animate-pulse"></div>
                    <span className="text-xs text-yellow-600">Connecting...</span>
                  </div>
                )}
                {connectionStatus === 'disconnected' && (
                  <div className="flex items-center gap-1">
                    <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                    <span className="text-xs text-red-600">Disconnected</span>
                  </div>
                )}
              </div>
            </div>
            <div className="w-20"></div>
          </div>
        </div>
      </header>

      <div className="max-w-6xl mx-auto">
        <div className="flex border-l border-r border-gray-300" style={{ height: 'calc(100vh - 64px)' }}>
          {/* Conversations List - Instagram style */}
          <div className="w-96 border-r border-gray-300 flex flex-col bg-white">
            <div className="p-4 border-b border-gray-300">
              <div className="flex justify-between items-center">
                <h2 className="text-xl font-semibold">{currentUser?.first_name}</h2>
                <button className="text-2xl">✏️</button>
              </div>
            </div>
            <div className="flex-1 overflow-y-auto">
              {conversations.length === 0 ? (
                <div className="p-8 text-center text-gray-500">
                  <p className="text-sm">No messages yet</p>
                  <p className="text-xs mt-2">Follow someone to start chatting</p>
                </div>
              ) : (
                conversations.map((user) => (
                  <button
                    key={user.id}
                    onClick={() => setSelectedUser(user)}
                    className={`w-full p-3 flex items-center hover:bg-gray-50 transition ${
                      selectedUser?.id === user.id ? 'bg-gray-100' : ''
                    }`}
                  >
                    <div className="w-14 h-14 rounded-full bg-gradient-to-tr from-yellow-400 via-red-500 to-purple-500 p-0.5 mr-3">
                      <div className="w-full h-full rounded-full bg-white p-0.5">
                        {user.avatar_path ? (
                          <img src={user.avatar_path} alt="" className="w-full h-full rounded-full object-cover" />
                        ) : (
                          <div className="w-full h-full rounded-full bg-gray-300 flex items-center justify-center">
                            <span className="text-lg font-semibold text-gray-600">
                              {user.first_name[0]}{user.last_name[0]}
                            </span>
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="flex-1 text-left">
                      <p className="font-semibold text-sm">
                        {user.first_name} {user.last_name}
                      </p>
                      {user.nickname && (
                        <p className="text-xs text-gray-500">@{user.nickname}</p>
                      )}
                    </div>
                  </button>
                ))
              )}
            </div>
          </div>

          {/* Messages Area - Instagram style */}
          <div className="flex-1 flex flex-col bg-white">
            {selectedUser ? (
              <>
                {/* Chat Header */}
                <div className="p-4 border-b border-gray-300 flex items-center">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-tr from-yellow-400 via-red-500 to-purple-500 p-0.5 mr-3">
                    <div className="w-full h-full rounded-full bg-white p-0.5">
                      {selectedUser.avatar_path ? (
                        <img src={selectedUser.avatar_path} alt="" className="w-full h-full rounded-full object-cover" />
                      ) : (
                        <div className="w-full h-full rounded-full bg-gray-300 flex items-center justify-center">
                          <span className="text-sm font-semibold text-gray-600">
                            {selectedUser.first_name[0]}{selectedUser.last_name[0]}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                  <div>
                    <p className="font-semibold text-sm">
                      {selectedUser.first_name} {selectedUser.last_name}
                    </p>
                  </div>
                </div>

                {/* Messages */}
                <div className="flex-1 overflow-y-auto p-4">
                  {messages.length === 0 ? (
                    <div className="flex flex-col items-center justify-center h-full text-center">
                      <div className="w-24 h-24 rounded-full border-2 border-black flex items-center justify-center mb-4">
                        <span className="text-4xl">
                          {selectedUser.first_name[0]}{selectedUser.last_name[0]}
                        </span>
                      </div>
                      <p className="font-semibold text-lg mb-1">
                        {selectedUser.first_name} {selectedUser.last_name}
                      </p>
                      {selectedUser.nickname && (
                        <p className="text-gray-500 text-sm mb-4">@{selectedUser.nickname}</p>
                      )}
                      <p className="text-gray-500 text-sm mb-4">
                        Send a message to start the conversation
                      </p>
                    </div>
                  ) : (
                    <div className="space-y-2">
                      {messages.map((message, index) => {
                        const isOwn = message.sender_id === currentUser?.id
                        const showAvatar = index === 0 || messages[index - 1].sender_id !== message.sender_id
                        
                        return (
                          <div
                            key={`${message.id}-${index}`}
                            className={`flex items-end ${isOwn ? 'justify-end' : 'justify-start'}`}
                          >
                            {!isOwn && showAvatar && (
                              <div className="w-6 h-6 rounded-full bg-gray-300 flex items-center justify-center mr-2 flex-shrink-0">
                                <span className="text-xs">
                                  {selectedUser.first_name[0]}
                                </span>
                              </div>
                            )}
                            {!isOwn && !showAvatar && <div className="w-6 mr-2"></div>}
                            <div
                              className={`max-w-xs px-4 py-2 rounded-3xl ${
                                isOwn
                                  ? 'bg-blue-500 text-white'
                                  : 'bg-gray-100 text-black border border-gray-300'
                              }`}
                            >
                              <p className="text-sm">{message.content}</p>
                            </div>
                          </div>
                        )
                      })}
                      <div ref={messagesEndRef} />
                    </div>
                  )}
                </div>

                {/* Message Input - Instagram style */}
                <div className="p-4 border-t border-gray-300 relative">
                  {/* Emoji Picker */}
                  {showEmojiPicker && (
                    <div className="absolute bottom-20 left-4 bg-white border border-gray-300 rounded-lg shadow-lg p-4 w-80 max-h-64 overflow-y-auto z-50">
                      <div className="grid grid-cols-8 gap-2">
                        {emojis.map((emoji, index) => (
                          <button
                            key={index}
                            type="button"
                            onClick={() => addEmoji(emoji)}
                            className="text-2xl hover:bg-gray-100 rounded p-1 transition"
                          >
                            {emoji}
                          </button>
                        ))}
                      </div>
                    </div>
                  )}
                  
                  <form onSubmit={sendMessage} className="flex items-center">
                    <button 
                      type="button" 
                      onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                      className="text-2xl mr-3 hover:scale-110 transition"
                    >
                      😊
                    </button>
                    <input
                      type="text"
                      value={newMessage}
                      onChange={(e) => setNewMessage(e.target.value)}
                      placeholder="Message..."
                      className="flex-1 px-4 py-2 border border-gray-300 rounded-full focus:outline-none focus:border-gray-400"
                    />
                    <button
                      type="submit"
                      disabled={!newMessage.trim()}
                      className={`ml-3 font-semibold ${
                        newMessage.trim() ? 'text-blue-500' : 'text-blue-300'
                      }`}
                    >
                      Send
                    </button>
                  </form>
                </div>
              </>
            ) : (
              <div className="flex-1 flex flex-col items-center justify-center text-center p-8">
                <div className="w-24 h-24 rounded-full border-2 border-black flex items-center justify-center mb-4">
                  <span className="text-4xl">💬</span>
                </div>
                <h2 className="text-2xl font-light mb-2">Your Messages</h2>
                <p className="text-gray-500 text-sm">
                  Send private messages to a friend
                </p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
