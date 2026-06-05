import { useEffect, useState } from 'react';
import { Plus, Trash2, TestTube, Terminal } from 'lucide-react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import type { Asset } from '../api/types';

export default function Hosts() {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [editing, setEditing] = useState<Asset | null>(null);
  const [form, setForm] = useState({ name: '', host: '', port: 22, username: 'root', auth_type: 'password', credential: '' });
  const [testing, setTesting] = useState<number | null>(null);
  const [testResult, setTestResult] = useState<string | null>(null);
  const navigate = useNavigate();

  const loadAssets = () => {
    api.getAssets().then(setAssets).catch(() => {});
  };

  useEffect(() => { loadAssets(); }, []);

  const openCreate = () => {
    setEditing(null);
    setForm({ name: '', host: '', port: 22, username: 'root', auth_type: 'password', credential: '' });
    setShowForm(true);
  };

  const openEdit = (a: Asset) => {
    setEditing(a);
    setForm({ name: a.name, host: a.host, port: a.port, username: a.username, auth_type: a.auth_type, credential: '' });
    setShowForm(true);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (editing) {
      await api.updateAsset(editing.id, form);
    } else {
      await api.createAsset(form);
    }
    setShowForm(false);
    loadAssets();
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this host?')) return;
    await api.deleteAsset(id);
    loadAssets();
  };

  const handleTest = async (id: number) => {
    setTesting(id);
    setTestResult(null);
    try {
      const res = await api.testAsset(id);
      setTestResult(res.status === 'connected' ? 'Connected OK' : `Failed: ${res.error}`);
    } catch (err: unknown) {
      setTestResult(`Failed: ${err instanceof Error ? err.message : 'Error'}`);
    }
    setTesting(null);
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-white">Managed Hosts</h1>
        <button onClick={openCreate} className="flex items-center gap-2 px-4 py-2 bg-success text-black text-sm font-semibold rounded-lg hover:bg-green-400 transition-colors">
          <Plus className="w-4 h-4" /> Add Host
        </button>
      </div>

      <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-800 text-gray-400 text-left">
              <th className="p-3">Name</th>
              <th className="p-3">Host</th>
              <th className="p-3">Port</th>
              <th className="p-3">Username</th>
              <th className="p-3">Auth</th>
              <th className="p-3">Status</th>
              <th className="p-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {assets.map((a) => (
              <tr key={a.id} className="border-b border-gray-800 text-gray-300 hover:bg-gray-800/50">
                <td className="p-3 font-medium">{a.name}</td>
                <td className="p-3 font-mono">{a.host}</td>
                <td className="p-3">{a.port}</td>
                <td className="p-3">{a.username}</td>
                <td className="p-3"><span className="px-2 py-0.5 text-xs rounded bg-gray-800">{a.auth_type}</span></td>
                <td className="p-3"><span className={`w-2 h-2 rounded-full inline-block mr-1 ${a.status === 'online' ? 'bg-success' : 'bg-gray-600'}`} /> {a.status}</td>
                <td className="p-3">
                  <div className="flex items-center gap-2">
                    <button onClick={() => navigate(`/terminal/${a.id}`)} className="p-1.5 text-gray-400 hover:text-success transition-colors" title="Connect">
                      <Terminal className="w-4 h-4" />
                    </button>
                    <button onClick={() => handleTest(a.id)} disabled={testing === a.id} className="p-1.5 text-gray-400 hover:text-info transition-colors" title="Test">
                      <TestTube className="w-4 h-4" />
                    </button>
                    <button onClick={() => openEdit(a)} className="p-1.5 text-gray-400 hover:text-white transition-colors" title="Edit">
                      <span className="text-xs">Edit</span>
                    </button>
                    <button onClick={() => handleDelete(a.id)} className="p-1.5 text-gray-400 hover:text-critical transition-colors" title="Delete">
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            ))}
            {assets.length === 0 && (
              <tr><td colSpan={7} className="p-3 text-center text-gray-500">No hosts configured</td></tr>
            )}
          </tbody>
        </table>
      </div>
      {testResult && (
        <div className={`mt-2 p-2 text-sm rounded ${testResult.startsWith('Connected') ? 'bg-green-900/50 text-green-400' : 'bg-red-900/50 text-red-400'}`}>
          {testResult}
        </div>
      )}

      {showForm && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50">
          <div className="bg-gray-900 border border-gray-700 rounded-xl p-6 w-full max-w-md">
            <h2 className="text-lg font-semibold text-white mb-4">{editing ? 'Edit Host' : 'Add Host'}</h2>
            <form onSubmit={handleSubmit} className="space-y-3">
              <input type="text" placeholder="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm" required />
              <input type="text" placeholder="Host/IP" value={form.host} onChange={(e) => setForm({ ...form, host: e.target.value })} className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm" required />
              <div className="flex gap-2">
                <input type="number" placeholder="Port" value={form.port} onChange={(e) => setForm({ ...form, port: +e.target.value })} className="w-24 px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm" />
                <input type="text" placeholder="Username" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm" required />
              </div>
              <select value={form.auth_type} onChange={(e) => setForm({ ...form, auth_type: e.target.value })} className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm">
                <option value="password">Password</option>
                <option value="key">SSH Key</option>
              </select>
              <textarea placeholder={form.auth_type === 'password' ? 'Password' : 'SSH Private Key'} value={form.credential} onChange={(e) => setForm({ ...form, credential: e.target.value })} className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded text-white text-sm h-20" required />
              <div className="flex justify-end gap-2 pt-2">
                <button type="button" onClick={() => setShowForm(false)} className="px-4 py-2 text-sm text-gray-400 hover:text-white">Cancel</button>
                <button type="submit" className="px-4 py-2 bg-success text-black text-sm font-semibold rounded hover:bg-green-400">{editing ? 'Save' : 'Create'}</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}
