import { BrowserRouter, Routes, Route, Link, Navigate, useNavigate } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import Transactions from "./pages/Transactions";
import { isAuthed, clearToken } from "./api/token";

function Protected({ children }) {
  if (!isAuthed()) return <Navigate to="/login" replace />;
  return children;
}

function Layout({ children }) {
  const navigate = useNavigate();

  function handleLogout() {
    clearToken();
    navigate("/login", { replace: true });
  }

  return (
    <div>
      <div style={{ padding: 12, borderBottom: "1px solid #444", display: "flex", gap: 12 }}>
        <Link to="/dashboard">Dashboard</Link>
        <Link to="/transactions">Transactions</Link>
        <div style={{ flex: 1 }} />
        <button onClick={handleLogout}>Выйти</button>
      </div>
      {children}
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />

        <Route
          path="/dashboard"
          element={
            <Protected>
              <Layout>
                <Dashboard />
              </Layout>
            </Protected>
          }
        />

        <Route
          path="/transactions"
          element={
            <Protected>
              <Layout>
                <Transactions />
              </Layout>
            </Protected>
          }
        />

        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
}