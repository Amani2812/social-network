import Link from 'next/link'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center">
          <h1 className="text-4xl md:text-6xl font-bold text-gray-900 mb-6">
            Welcome to Social Network
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-2xl mx-auto">
            Connect with friends, share your thoughts, and discover amazing communities.
            Join our social network today!
          </p>
          <div className="space-x-4">
            <Link
              href="/register"
              className="inline-block bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-6 rounded-lg transition duration-300"
            >
              Get Started
            </Link>
            <Link
              href="/login"
              className="inline-block bg-white hover:bg-gray-50 text-indigo-600 font-bold py-3 px-6 rounded-lg border border-indigo-600 transition duration-300"
            >
              Sign In
            </Link>
          </div>
        </div>

        <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <div className="text-indigo-600 text-3xl mb-4">👥</div>
            <h3 className="text-xl font-semibold mb-2">Connect</h3>
            <p className="text-gray-600">Follow friends and discover new people with similar interests.</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md">
            <div className="text-indigo-600 text-3xl mb-4">📝</div>
            <h3 className="text-xl font-semibold mb-2">Share</h3>
            <p className="text-gray-600">Post your thoughts, photos, and experiences with your network.</p>
          </div>
          <div className="bg-white p-6 rounded-lg shadow-md">
            <div className="text-indigo-600 text-3xl mb-4">🎉</div>
            <h3 className="text-xl font-semibold mb-2">Discover</h3>
            <p className="text-gray-600">Join groups, attend events, and explore amazing communities.</p>
          </div>
        </div>
      </div>
    </div>
  )
}
