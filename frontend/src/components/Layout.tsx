import { useState } from 'react';
import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { ShieldCheck, LayoutDashboard, Server, TerminalSquare, ScrollText, Moon, Sun, LogOut } from 'lucide-react';

const navItems = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/hosts', icon: Server, label: 'Hosts' },
  { to: '/sessions', icon: TerminalSquare, label: 'Sessions' },
  { to: '/audit', icon: ScrollText, label: 'Audit Log' },
];

export default function Layout() {
  const [dark, setDark] = useState(true);
  const navigate = useNavigate();

  const toggleDark = () => {
    const next = !dark;
    setDark(next);
    document.documentElement.classList.toggle('dark', next);
  };

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  return (
    <div className="flex h-screen bg-gray-950">
      {/* Sidebar */}
      <aside className="w-56 flex-shrink-0 bg-gray-900 border-r border-gray-800 flex flex-col">
        <div className="p-4 border-b border-gray-800 flex items-center gap-2">
          <ShieldCheck className="w-6 h-6 text-success" />
          <span className="text-lg font-bold text-white">Bastion</span>
        </div>

        <nav className="flex-1 p-2 space-y-1">
          {navItems.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-colors ${
                  isActive
                    ? 'bg-gray-800 text-success'
                    : 'text-gray-400 hover:text-white hover:bg-gray-800'
                }`
              }
            >
              <Icon className="w-5 h-5" />
              {label}
            </NavLink>
          ))}
        </nav>

        <div className="p-3 border-t border-gray-800 space-y-2">
          <button
            onClick={toggleDark}
            className="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-400 hover:text-white rounded-lg hover:bg-gray-800 transition-colors"
          >
            {dark ? <Sun className="w-4 h-4" /> : <Moon className="w-4 h-4" />}
            {dark ? 'Light' : 'Dark'} Mode
          </button>
          <button
            onClick={handleLogout}
            className="flex items-center gap-2 w-full px-3 py-2 text-sm text-gray-400 hover:text-red-400 rounded-lg hover:bg-gray-800 transition-colors"
          >
            <LogOut className="w-4 h-4" />
            Logout
          </button>
        </div>
      </aside>

      {/* Main */}
      <main className="flex-1 overflow-auto p-6">
        <Outlet />
      </main>
    </div>
  );
}
