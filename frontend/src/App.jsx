import {
  BrowserRouter,
  Navigate,
  NavLink,
  Outlet,
  Route,
  Routes,
  useNavigate,
} from "react-router-dom";
import { clearToken, isAuthed } from "./api/token";
import "./App.css";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import Profile from "./pages/Profile";
import Register from "./pages/Register";
import Transactions from "./pages/Transactions";

function Protected() {
  if (!isAuthed()) return <Navigate to="/login" replace />;
  return <Outlet />;
}

function AppLayout() {
  const navigate = useNavigate();

  function handleLogout() {
    clearToken();
    navigate("/login", { replace: true });
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand-block">
          <div className="brand-mark">₽</div>
          <div>
            <div className="brand-title">Ledger Lite</div>
            <div className="brand-subtitle">Личные финансы</div>
          </div>
        </div>

        <nav className="nav-menu">
          <NavLink
            to="/dashboard"
            className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}
          >
            <span className="nav-link-icon">▣</span>
            <span>Дашборд</span>
          </NavLink>

          <NavLink
            to="/transactions"
            className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}
          >
            <span className="nav-link-icon">◫</span>
            <span>Транзакции</span>
          </NavLink>

          <NavLink
            to="/profile"
            className={({ isActive }) => `nav-link ${isActive ? "active" : ""}`}
          >
            <span className="nav-link-icon">◎</span>
            <span>Профиль</span>
          </NavLink>
        </nav>

        <div className="sidebar-footer">
          <div className="status-pill">Онлайн</div>
          <button className="btn btn-danger btn-block" onClick={handleLogout}>
            Выйти
          </button>
        </div>
      </aside>

      <div className="content-area">
        <header className="topbar">
          <div>
            <div className="topbar-title">Финансовая панель</div>
            <div className="topbar-subtitle">
              Доходы, расходы, категории и аналитика в одном месте
            </div>
          </div>
        </header>

        <main className="page-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

export default function App() {
  const defaultPath = isAuthed() ? "/dashboard" : "/login";

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            isAuthed() ? <Navigate to="/dashboard" replace /> : <Login />
          }
        />

        <Route
          path="/register"
          element={
            isAuthed() ? <Navigate to="/dashboard" replace /> : <Register />
          }
        />

        <Route element={<Protected />}>
          <Route element={<AppLayout />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/transactions" element={<Transactions />} />
            <Route path="/profile" element={<Profile />} />
          </Route>
        </Route>

        <Route path="*" element={<Navigate to={defaultPath} replace />} />
      </Routes>
    </BrowserRouter>
  );
}
