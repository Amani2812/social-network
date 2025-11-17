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
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)
  const selectedUserRef = useRef<User | null>(null)

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
      if (wsRef.current) {
        wsRef.current.close()
      }
    }
  }, [])

  useEffect(() => {
    if (currentUser) {
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
    
    const websocket = new WebSocket('ws://localhost:8080/ws')
    
    websocket.onopen = () => {
      console.log('✅ Connected')
    }
    
    websocket.onmessage = (event) => {
      const data = JSON.parse(event.data)
      
      if (data.type === 'private') {
        const newMsg: Message = {
          id: Date.now(),
          sender_id: data.sender_id,
          receiver_id: data.receiver_id,
          content: data.content,
          created_at: data.timestamp || new Date().toISOString(),
        }
        
        // Add to messages immediately
        setMessages(prev => [...prev, newMsg])
        
        // Add sender to conversations if new
        if (newMsg.sender_id !== currentUser.id) {
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
      }
    }
    
    websocket.onerror = (error) => {
      console.error('WebSocket error:', error)
    }
    
    websocket.onclose = () => {
      setTimeout(() => {
        if (currentUser) {
          connectWebSocket()
        }
      }, 3000)
    }
    
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
      const response = await fetch(`http://localhost:8080/api/follow/following?user_id=${userId}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setConversations(data || [])
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

  const sendMessage = async (e: React.FormEvent) => {
    e.preventDefault()
    
    if (!newMessage.trim() || !selectedUser || !ws || ws.readyState !== WebSocket.OPEN) {
      return
    }

    const messageData = {
      type: 'private',
      receiver_id: selectedUser.id,
      content: newMessage,
    }

    try {
      ws.send(JSON.stringify(messageData))
      setNewMessage('')
    } catch (error) {
      console.error('Failed to send')
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
            <h1 className="text-xl font-semibold">{currentUser?.first_name}</h1>
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
                <div className="p-4 border-t border-gray-300">
                  <form onSubmit={sendMessage} className="flex items-center">
                    <button type="button" className="text-2xl mr-3">😊</button>
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
