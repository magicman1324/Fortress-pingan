import { Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Hosts from './pages/Hosts';
import Terminal from './pages/Terminal';
import Sessions from './pages/Sessions';
import AuditLog from './pages/AuditLog';

function isAuthenticated() {
  return !!localStorage.getItem('token');
}

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/terminal/:assetId" element={
        <ProtectedRoute><Terminal /></ProtectedRoute>
      } />
      <Route element={
        <ProtectedRoute><Layout /></ProtectedRoute>
      }>
        <Route path="/" element={<Dashboard />} />
        <Route path="/hosts" element={<Hosts />} />
        <Route path="/sessions" element={<Sessions />} />
        <Route path="/audit" element={<AuditLog />} />
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}
