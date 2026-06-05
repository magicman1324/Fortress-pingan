import { useEffect, useState } from 'react';
import { XCircle } from 'lucide-react';
import { api } from '../api/client';
import type { Session } from '../api/types';

export default function Sessions() {
  const [sessions, setSessions] = useState<Session[]>([]);

  const loadSessions = () => {
    api.getSessions(1, 50).then(setSessions).catch(() => {});
  };

  useEffect(() => { loadSessions(); }, []);

  const handleTerminate = async (id: number) => {
    if (!confirm('Force terminate this session?')) return;
    await api.terminateSession(id);
    loadSessions();
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-white">Sessions</h1>
        <button onClick={loadSessions} className="px-3 py-1.5 text-sm text-gray-400 hover:text-white bg-gray-800 rounded-lg">Refresh</button>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-800 text-gray-400 text-left">
              <th className="p-3">ID</th>
              <th className="p-3">User</th>
              <th className="p-3">Asset</th>
              <th className="p-3">Status</th>
              <th className="p-3">Client IP</th>
              <th className="p-3">Started</th>
              <th className="p-3">Ended</th>
              <th className="p-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {sessions.map((s) => (
              <tr key={s.id} className="border-b border-gray-800 text-gray-300 hover:bg-gray-800/50">
                <td className="p-3 font-mono text-xs">#{s.id}</td>
                <td className="p-3">{s.username}</td>
                <td className="p-3">{s.asset_name}</td>
                <td className="p-3">
                  <span className={`px-2 py-0.5 text-xs rounded ${s.status === 'active' ? 'bg-green-900/50 text-success' : 'bg-gray-800 text-gray-400'}`}>
                    {s.status}
                  </span>
                </td>
                <td className="p-3 font-mono text-xs">{s.client_ip}</td>
                <td className="p-3 text-xs text-gray-500">{new Date(s.started_at).toLocaleString()}</td>
                <td className="p-3 text-xs text-gray-500">{s.ended_at ? new Date(s.ended_at).toLocaleString() : '-'}</td>
                <td className="p-3">
                  {s.status === 'active' && (
                    <button onClick={() => handleTerminate(s.id)} className="p-1 text-gray-400 hover:text-critical" title="Terminate">
                      <XCircle className="w-4 h-4" />
                    </button>
                  )}
                </td>
              </tr>
            ))}
            {sessions.length === 0 && (
              <tr><td colSpan={8} className="p-3 text-center text-gray-500">No sessions</td></tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
