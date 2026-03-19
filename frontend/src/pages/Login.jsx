import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { loginUser } from "../api/auth";
import { setToken } from "../api/token";

export default function Login() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("cat@test.com");
  const [password, setPassword] = useState("password123");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const data = await loginUser(email, password);
      setToken(data.access_token);
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(err?.message || "Не удалось войти в систему");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-badge">Добро пожаловать</div>
        <h1 className="auth-title">Вход в приложение</h1>
        <p className="auth-subtitle">
          Войди в аккаунт, чтобы посмотреть аналитику и управлять транзакциями.
        </p>

        <form className="form-grid" onSubmit={handleSubmit}>
          <label className="field">
            <span className="field-label">Электронная почта</span>
            <input
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com"
              autoComplete="email"
            />
          </label>

          <label className="field">
            <span className="field-label">Пароль</span>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Введите пароль"
              autoComplete="current-password"
            />
          </label>

          <div className="form-actions">
            <button
              className="btn btn-primary btn-block"
              type="submit"
              disabled={loading}
            >
              {loading ? "Входим..." : "Войти"}
            </button>
          </div>

          {error && <div className="alert">{error}</div>}
        </form>

        <div style={{ marginTop: 18 }} className="helper-text">
          Ещё нет аккаунта?{" "}
          <Link to="/register" style={{ color: "#a5b4fc", fontWeight: 700 }}>
            Зарегистрироваться
          </Link>
        </div>
      </div>
    </div>
  );
}
