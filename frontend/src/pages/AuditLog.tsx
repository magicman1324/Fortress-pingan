import { useEffect, useState, useCallback } from 'react';
import { Search } from 'lucide-react';
import { api } from '../api/client';
import type { AuditLog as AuditLogType } from '../api/types';

export default function AuditLog() {
  const [logs, setLogs] = useState<AuditLogType[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [keyword, setKeyword] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const pageSize = 20;

  const loadLogs = useCallback(async () => {
    const res = await api.getAuditLogs({ page, size: pageSize, keyword });
    setLogs(res.items);
    setTotal(res.total);
  }, [page, keyword]);

  useEffect(() => { loadLogs(); }, [loadLogs]);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    setKeyword(searchInput);
    setPage(1);
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div>
      <h1 className="text-2xl font-bold text-white mb-6">Audit Log</h1>

      <form onSubmit={handleSearch} className="mb-4 flex gap-2">
        <div className="relative flex-1 max-w-md">
          <Search className="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" />
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder="Search commands..."
            className="w-full pl-9 pr-3 py-2 bg-gray-900 border border-gray-700 rounded-lg text-white text-sm focus:outline-none focus:border-success"
          />
        </div>
        <button type="submit" className="px-4 py-2 bg-gray-800 text-white text-sm rounded-lg hover:bg-gray-700">Search</button>
      </form>

      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-800 text-gray-400 text-left">
              <th className="p-3">Time</th>
              <th className="p-3">User</th>
              <th className="p-3">Asset</th>
              <th className="p-3">Session</th>
              <th className="p-3">Command</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log) => (
              <tr key={log.id} className="border-b border-gray-800 text-gray-300 hover:bg-gray-800/50">
                <td className="p-3 text-xs text-gray-500 whitespace-nowrap">{new Date(log.executed_at).toLocaleString()}</td>
                <td className="p-3">{log.username}</td>
                <td className="p-3">{log.asset_name}</td>
                <td className="p-3 font-mono text-xs">#{log.session_id}</td>
                <td className="p-3 font-mono text-xs max-w-md truncate">{log.command}</td>
              </tr>
            ))}
            {logs.length === 0 && (
              <tr><td colSpan={5} className="p-3 text-center text-gray-500">No audit logs found</td></tr>
            )}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between mt-4">
          <span className="text-sm text-gray-400">Total: {total} records</span>
          <div className="flex gap-1">
            <button onClick={() => setPage(p => Math.max(1, p - 1))} disabled={page <= 1} className="px-3 py-1 text-sm rounded bg-gray-800 text-gray-400 disabled:opacity-50 hover:text-white">Prev</button>
            <span className="px-3 py-1 text-sm text-gray-400">Page {page} / {totalPages}</span>
            <button onClick={() => setPage(p => Math.min(totalPages, p + 1))} disabled={page >= totalPages} className="px-3 py-1 text-sm rounded bg-gray-800 text-gray-400 disabled:opacity-50 hover:text-white">Next</button>
          </div>
        </div>
      )}
    </div>
  );
}
