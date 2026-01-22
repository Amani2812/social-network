'use client'

import { useEffect, useState } from 'react'

interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  avatar_path?: string
  nickname?: string
}

interface Comment {
  id: number
  post_id: number
  user_id: number
  content: string
  image_path?: string
  created_at: string
  user?: User
  likes?: number
  dislikes?: number
  user_reaction?: 'like' | 'dislike' | null
}

interface CommentsProps {
  postId: number
  currentUserId?: number
}

export default function Comments({ postId, currentUserId }: CommentsProps) {
  const [comments, setComments] = useState<Comment[]>([])
  const [newComment, setNewComment] = useState('')
  const [loading, setLoading] = useState(true)
  const [showComments, setShowComments] = useState(false)
  const [selectedImage, setSelectedImage] = useState<File | null>(null)
  const [imagePreview, setImagePreview] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const handleReaction = async (commentId: number, reaction: 'like' | 'dislike') => {
    try {
      const response = await fetch('http://localhost:8080/api/comments/react', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          comment_id: commentId,
          reaction: reaction,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        // Refresh comments to get updated counts
        fetchComments()
      }
    } catch (err) {
      console.error('Failed to react to comment:', err)
    }
  }

  useEffect(() => {
    if (showComments) {
      fetchComments()
    }
  }, [showComments, postId])

  const fetchComments = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/comments/get?post_id=${postId}`, {
        credentials: 'include',
      })
      if (response.ok) {
        const data = await response.json()
        setComments(data || [])
      }
    } catch (err) {
      console.error('Failed to fetch comments:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleImageSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
      if (!validTypes.includes(file.type)) {
        alert('Please select a valid image file (JPEG, PNG, or GIF)')
        return
      }

      if (file.size > 10 * 1024 * 1024) {
        alert('File size must be less than 10MB')
        return
      }

      setSelectedImage(file)
      
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

  const handleSubmitComment = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!newComment.trim()) return

    setSubmitting(true)

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
          setSubmitting(false)
          return
        }
      }

      // Create comment
      const response = await fetch('http://localhost:8080/api/comments', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          post_id: postId,
          content: newComment,
          image_path: imagePath,
        }),
        credentials: 'include',
      })

      if (response.ok) {
        setNewComment('')
        setSelectedImage(null)
        setImagePreview(null)
        fetchComments() // Refresh comments
      }
    } catch (err) {
      console.error('Failed to create comment:', err)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="mt-4">
      {/* Toggle Comments Button */}
      <button
        onClick={() => setShowComments(!showComments)}
        className="text-blue-600 hover:text-blue-800 text-sm font-medium"
      >
        {showComments ? '▼ Hide Comments' : `▶ Show Comments (${comments.length})`}
      </button>

      {/* Comments Section */}
      {showComments && (
        <div className="mt-4 space-y-4">
          {/* Comment Form */}
          {currentUserId && (
            <form onSubmit={handleSubmitComment} className="border-t pt-4">
              <textarea
                value={newComment}
                onChange={(e) => setNewComment(e.target.value)}
                placeholder="Write a comment..."
                className="w-full p-2 border border-gray-300 rounded-md resize-none focus:ring-blue-500 focus:border-blue-500 text-sm"
                rows={2}
              />

              {/* Image Preview */}
              {imagePreview && (
                <div className="mt-2 relative">
                  <img 
                    src={imagePreview} 
                    alt="Preview" 
                    className="w-full max-h-32 object-cover rounded-md"
                  />
                  <button
                    type="button"
                    onClick={handleRemoveImage}
                    className="absolute top-1 right-1 bg-red-600 text-white rounded-full p-1 text-xs hover:bg-red-700"
                  >
                    ✕
                  </button>
                </div>
              )}

              <div className="mt-2 flex items-center justify-between">
                <label className="cursor-pointer bg-gray-100 hover:bg-gray-200 text-gray-700 px-3 py-1 rounded-md text-sm flex items-center">
                  <span className="mr-1">📷</span>
                  <span>Image</span>
                  <input
                    type="file"
                    accept="image/jpeg,image/jpg,image/png,image/gif"
                    onChange={handleImageSelect}
                    className="hidden"
                  />
                </label>

                <button
                  type="submit"
                  disabled={submitting}
                  className={`px-4 py-1 rounded-md text-sm ${
                    submitting 
                      ? 'bg-gray-400 cursor-not-allowed' 
                      : 'bg-blue-600 hover:bg-blue-700'
                  } text-white`}
                >
                  {submitting ? 'Posting...' : 'Comment'}
                </button>
              </div>
            </form>
          )}

          {/* Comments List */}
          {loading ? (
            <div className="text-center py-4 text-gray-500 text-sm">Loading comments...</div>
          ) : comments.length === 0 ? (
            <div className="text-center py-4 text-gray-500 text-sm">No comments yet. Be the first to comment!</div>
          ) : (
            <div className="space-y-3">
              {comments.map((comment) => (
                <div key={comment.id} className="bg-gray-50 rounded-lg p-3">
                  <div className="flex items-start">
                    <div className="w-8 h-8 bg-gray-300 rounded-full flex items-center justify-center mr-2 flex-shrink-0">
                      {comment.user?.avatar_path ? (
                        <img 
                          src={comment.user.avatar_path} 
                          alt="Avatar" 
                          className="w-full h-full rounded-full object-cover" 
                        />
                      ) : (
                        <span className="text-xs text-gray-600">
                          {comment.user?.first_name[0]}{comment.user?.last_name[0]}
                        </span>
                      )}
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <p className="font-semibold text-sm text-gray-900">
                          {comment.user?.first_name} {comment.user?.last_name}
                        </p>
                        <p className="text-xs text-gray-500">
                          {new Date(comment.created_at).toLocaleDateString()}
                        </p>
                      </div>
                      <p className="text-sm text-gray-700 mt-1">{comment.content}</p>
                      {comment.image_path && (
                        <img 
                          src={comment.image_path} 
                          alt="Comment image" 
                          className="mt-2 w-full max-h-48 object-cover rounded-md" 
                        />
                      )}
                      
                      {/* Like/Dislike Buttons */}
                      <div className="flex items-center gap-4 mt-2">
                        <button
                          onClick={() => handleReaction(comment.id, 'like')}
                          className={`flex items-center gap-1 text-sm ${
                            comment.user_reaction === 'like' 
                              ? 'text-blue-600 font-semibold' 
                              : 'text-gray-600 hover:text-blue-600'
                          } transition`}
                        >
                          <span className="text-lg">👍</span>
                          <span>{comment.likes || 0}</span>
                        </button>
                        <button
                          onClick={() => handleReaction(comment.id, 'dislike')}
                          className={`flex items-center gap-1 text-sm ${
                            comment.user_reaction === 'dislike' 
                              ? 'text-red-600 font-semibold' 
                              : 'text-gray-600 hover:text-red-600'
                          } transition`}
                        >
                          <span className="text-lg">👎</span>
                          <span>{comment.dislikes || 0}</span>
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
