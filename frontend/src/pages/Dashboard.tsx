import { useEffect, useState } from 'react';
import { Activity, Server, TerminalSquare, ScrollText } from 'lucide-react';
import { api } from '../api/client';
import type { DashboardStats } from '../api/types';

export default function Dashboard() {
  const [stats, setStats] = useState<DashboardStats | null>(null);

  useEffect(() => {
    api.getDashboardStats().then(setStats).catch(() => {});
  }, []);

  const cards = [
    { label: 'Active Sessions', value: stats?.active_sessions ?? '-', icon: Activity, color: 'text-success' },
    { label: 'Total Sessions', value: stats?.total_sessions ?? '-', icon: TerminalSquare, color: 'text-info' },
    { label: 'Managed Hosts', value: '-', icon: Server, color: 'text-warning' },
    { label: 'Today\'s Audits', value: stats?.today_audit_count ?? '-', icon: ScrollText, color: 'text-purple-400' },
  ];

  return (
    <div>
      <h1 className="text-2xl font-bold text-white mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {cards.map(({ label, value, icon: Icon, color }) => (
          <div key={label} className="bg-gray-900 border border-gray-800 rounded-xl p-4">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm text-gray-400">{label}</span>
              <Icon className={`w-5 h-5 ${color}`} />
            </div>
            <div className="text-2xl font-bold text-white">{value}</div>
          </div>
        ))}
      </div>

      <h2 className="text-lg font-semibold text-white mb-4">Recent Audit Logs</h2>
      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-800 text-gray-400 text-left">
              <th className="p-3">Time</th>
              <th className="p-3">User</th>
              <th className="p-3">Asset</th>
              <th className="p-3">Command</th>
            </tr>
          </thead>
          <tbody>
            {stats?.recent_audits?.map((log) => (
              <tr key={log.id} className="border-b border-gray-800 text-gray-300">
                <td className="p-3 text-gray-500">{new Date(log.executed_at).toLocaleString()}</td>
                <td className="p-3">{log.username}</td>
                <td className="p-3">{log.asset_name}</td>
                <td className="p-3 font-mono text-xs max-w-xs truncate">{log.command}</td>
              </tr>
            )) || (
              <tr>
                <td colSpan={4} className="p-3 text-center text-gray-500">No audit logs yet</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
