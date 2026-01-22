'use client'

import { useEffect, useState } from 'react'
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
  going_count: number
  not_going_count: number
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

  useEffect(() => {
    fetchUser()
    fetchGroup()
    fetchPosts()
    fetchEvents()
    fetchJoinRequests()
  }, [params.id])

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

    try {
      const response = await fetch('http://localhost:8080/api/groups/posts/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          group_id: parseInt(params.id as string),
          content: newPost,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        setNewPost('')
        await fetchPosts()
        alert('Post created successfully! ✓')
      } else {
        const errorData = await response.json()
        alert(`Failed to create post: ${errorData.error || 'Unknown error'}`)
      }
    } catch (err) {
      console.error('Failed to create post:', err)
      alert('Failed to create post. You might not be a member of this group.')
    } finally {
      setPosting(false)
    }
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

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl">Loading...</div>
      </div>
    )
  }

  if (!group) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-xl text-red-600">Group not found</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <h1 className="text-2xl font-bold text-gray-900">{group.title}</h1>
            <div className="flex items-center space-x-4">
              <button
                onClick={() => router.push('/groups')}
                className="text-blue-600 hover:text-blue-800"
              >
                All Groups
              </button>
              <button
                onClick={() => router.push('/dashboard')}
                className="text-blue-600 hover:text-blue-800"
              >
                Dashboard
              </button>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-4xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        {/* Group Info */}
        <div className="bg-white shadow rounded-lg p-6 mb-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-2">{group.title}</h2>
          {group.description && (
            <p className="text-gray-600 mb-4">{group.description}</p>
          )}
          <p className="text-sm text-gray-500">
            Created {new Date(group.created_at).toLocaleDateString()}
          </p>
        </div>

        {/* Tabs */}
        <div className="bg-white shadow rounded-lg mb-6">
          <div className="border-b border-gray-200">
            <nav className="flex -mb-px">
              <button
                onClick={() => setActiveTab('posts')}
                className={`py-4 px-6 text-sm font-medium ${
                  activeTab === 'posts'
                    ? 'border-b-2 border-blue-600 text-blue-600'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                📝 Posts
              </button>
              <button
                onClick={() => setActiveTab('events')}
                className={`py-4 px-6 text-sm font-medium ${
                  activeTab === 'events'
                    ? 'border-b-2 border-blue-600 text-blue-600'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                📅 Events
              </button>
              {isAdmin && (
                <button
                  onClick={() => setActiveTab('members')}
                  className={`py-4 px-6 text-sm font-medium ${
                    activeTab === 'members'
                      ? 'border-b-2 border-blue-600 text-blue-600'
                      : 'text-gray-500 hover:text-gray-700'
                  }`}
                >
                  👥 Manage Members
                  {joinRequests.length > 0 && (
                    <span className="ml-2 bg-red-600 text-white text-xs font-bold px-2 py-1 rounded-full">
                      {joinRequests.length}
                    </span>
                  )}
                </button>
              )}
            </nav>
          </div>

          {/* Posts Tab */}
          {activeTab === 'posts' && (
            <div className="p-6">
              {/* Create Post Form */}
              <form onSubmit={handleCreatePost} className="mb-6">
                <textarea
                  value={newPost}
                  onChange={(e) => setNewPost(e.target.value)}
                  placeholder="Share something with the group..."
                  className="w-full p-3 border border-gray-300 rounded-md resize-none focus:ring-blue-500 focus:border-blue-500"
                  rows={3}
                />
                <div className="mt-2 flex justify-end">
                  <button
                    type="submit"
                    disabled={posting}
                    className={`px-6 py-2 rounded-md ${
                      posting
                        ? 'bg-gray-400 cursor-not-allowed'
                        : 'bg-blue-600 hover:bg-blue-700'
                    } text-white`}
                  >
                    {posting ? 'Posting...' : 'Post'}
                  </button>
                </div>
              </form>

              {/* Posts List */}
              {posts.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  No posts yet. Be the first to post!
                </div>
              ) : (
                <div className="space-y-4">
                  {posts.map((post) => (
                    <div key={post.id} className="border border-gray-200 rounded-lg p-4">
                      <div className="flex items-center mb-3">
                        <div className="w-10 h-10 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                          <span className="text-sm text-gray-600">
                            {post.user?.first_name[0]}{post.user?.last_name[0]}
                          </span>
                        </div>
                        <div>
                          <p className="font-semibold text-gray-900">
                            {post.user?.first_name} {post.user?.last_name}
                          </p>
                          <p className="text-xs text-gray-500">
                            {new Date(post.created_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                      <p className="text-gray-700">{post.content}</p>
                      {post.image_path && (
                        <img
                          src={`http://localhost:8080${post.image_path}`}
                          alt="Post"
                          className="mt-3 w-full rounded-md"
                        />
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Members Tab */}
          {activeTab === 'members' && isAdmin && (
            <div className="p-6">
              <h3 className="text-lg font-semibold mb-4">Join Requests</h3>
              
              {joinRequests.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  No pending join requests
                </div>
              ) : (
                <div className="space-y-3">
                  {joinRequests.map((request) => (
                    <div key={request.user_id} className="border border-gray-200 rounded-lg p-4 flex items-center justify-between">
                      <div className="flex items-center">
                        <div className="w-12 h-12 bg-gray-300 rounded-full flex items-center justify-center mr-3">
                          <span className="text-lg text-gray-600">
                            {request.user?.first_name?.[0]}{request.user?.last_name?.[0]}
                          </span>
                        </div>
                        <div>
                          <p className="font-semibold text-gray-900">
                            {request.user?.first_name} {request.user?.last_name}
                          </p>
                          {request.user?.nickname && (
                            <p className="text-sm text-gray-600">@{request.user.nickname}</p>
                          )}
                          <p className="text-sm text-gray-500">
                            Requested {new Date(request.created_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleJoinRequest(request.user_id, true)}
                          className="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium"
                        >
                          ✓ Accept
                        </button>
                        <button
                          onClick={() => handleJoinRequest(request.user_id, false)}
                          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md text-sm font-medium"
                        >
                          ✗ Decline
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Events Tab */}
          {activeTab === 'events' && (
            <div className="p-6">
              {/* Create Event Button */}
              <div className="mb-6">
                <button
                  onClick={() => setShowEventForm(!showEventForm)}
                  className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md"
                >
                  {showEventForm ? '✕ Cancel' : '+ Create Event'}
                </button>
              </div>

              {/* Create Event Form */}
              {showEventForm && (
                <form onSubmit={handleCreateEvent} className="mb-6 border border-gray-200 rounded-lg p-4">
                  <h3 className="text-lg font-semibold mb-4">Create New Event</h3>
                  <div className="mb-3">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Event Title *
                    </label>
                    <input
                      type="text"
                      value={newEventTitle}
                      onChange={(e) => setNewEventTitle(e.target.value)}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                    />
                  </div>
                  <div className="mb-3">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Description
                    </label>
                    <textarea
                      value={newEventDescription}
                      onChange={(e) => setNewEventDescription(e.target.value)}
                      className="w-full p-2 border border-gray-300 rounded-md resize-none"
                      rows={2}
                    />
                  </div>
                  <div className="mb-3">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Event Date & Time *
                    </label>
                    <input
                      type="datetime-local"
                      value={newEventTime}
                      onChange={(e) => setNewEventTime(e.target.value)}
                      className="w-full p-2 border border-gray-300 rounded-md"
                      required
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={creatingEvent}
                    className={`w-full py-2 rounded-md ${
                      creatingEvent
                        ? 'bg-gray-400 cursor-not-allowed'
                        : 'bg-blue-600 hover:bg-blue-700'
                    } text-white`}
                  >
                    {creatingEvent ? 'Creating...' : 'Create Event'}
                  </button>
                </form>
              )}

              {/* Events List */}
              {events.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  No events yet. Create one!
                </div>
              ) : (
                <div className="space-y-4">
                  {events.map((event) => (
                    <div key={event.id} className="border border-gray-200 rounded-lg p-4">
                      <h4 className="text-lg font-semibold text-gray-900 mb-2">
                        {event.title}
                      </h4>
                      {event.description && (
                        <p className="text-gray-600 mb-2">{event.description}</p>
                      )}
                      <p className="text-sm text-gray-500 mb-2">
                        📅 {new Date(event.event_time).toLocaleString()}
                      </p>
                      
                      {/* Attendance Counts */}
                      <div className="flex items-center gap-4 mb-3 text-sm">
                        <span className="text-green-600 font-medium">
                          ✓ {event.going_count} Going
                        </span>
                        <span className="text-red-600 font-medium">
                          ✗ {event.not_going_count} Not Going
                        </span>
                      </div>
                      
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleEventResponse(event.id, 'going')}
                          disabled={event.user_response === 'going'}
                          className={`px-4 py-1 rounded-md text-sm font-medium ${
                            event.user_response === 'going'
                              ? 'bg-green-700 text-white cursor-default'
                              : 'bg-green-600 hover:bg-green-700 text-white'
                          }`}
                        >
                          {event.user_response === 'going' ? '✓ Going (You)' : '✓ Going'}
                        </button>
                        <button
                          onClick={() => handleEventResponse(event.id, 'not_going')}
                          disabled={event.user_response === 'not_going'}
                          className={`px-4 py-1 rounded-md text-sm font-medium ${
                            event.user_response === 'not_going'
                              ? 'bg-red-700 text-white cursor-default'
                              : 'bg-red-600 hover:bg-red-700 text-white'
                          }`}
                        >
                          {event.user_response === 'not_going' ? '✗ Not Going (You)' : '✗ Not Going'}
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Back Button */}
        <div className="text-center">
          <button
            onClick={() => router.push('/groups')}
            className="text-blue-600 hover:text-blue-800"
          >
            ← Back to Groups
          </button>
        </div>
      </div>
    </div>
  )
}
