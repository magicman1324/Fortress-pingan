import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ShieldCheck } from 'lucide-react';
import { api } from '../api/client';

export default function Login() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      const res = await api.login(username, password);
      localStorage.setItem('token', res.token);
      localStorage.setItem('username', res.username);
      localStorage.setItem('role', res.role);
      navigate('/');
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-950">
      <div className="w-full max-w-sm p-8 bg-gray-900 rounded-xl border border-gray-800">
        <div className="flex items-center justify-center gap-3 mb-6">
          <ShieldCheck className="w-10 h-10 text-success" />
          <h1 className="text-2xl font-bold text-white">Bastion Host</h1>
        </div>
        <p className="text-center text-gray-400 text-sm mb-8">轻量级运维堡垒机</p>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm text-gray-400 mb-1">Username</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-success"
              placeholder="admin"
              required
            />
          </div>
          <div>
            <label className="block text-sm text-gray-400 mb-1">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:outline-none focus:border-success"
              placeholder="••••••••"
              required
            />
          </div>
          {error && (
            <div className="text-critical text-sm">{error}</div>
          )}
          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-success text-black font-semibold rounded-lg hover:bg-green-400 transition-colors disabled:opacity-50"
          >
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  );
}
