import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { loginUser, registerUser } from "../api/auth";
import { setToken } from "../api/token";

export default function Register() {
  const navigate = useNavigate();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordRepeat, setPasswordRepeat] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(e) {
    e.preventDefault();
    setError("");

    const trimmedEmail = email.trim();

    if (!trimmedEmail) {
      setError("Укажи электронную почту");
      return;
    }

    if (password.length < 8) {
      setError("Пароль должен быть не короче 8 символов");
      return;
    }

    if (password !== passwordRepeat) {
      setError("Пароли не совпадают");
      return;
    }

    setLoading(true);

    try {
      await registerUser(trimmedEmail, password);
      const authData = await loginUser(trimmedEmail, password);
      setToken(authData.access_token);
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(err?.message || "Не удалось зарегистрировать пользователя");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <div className="auth-card">
        <div className="auth-badge">Новый аккаунт</div>
        <h1 className="auth-title">Регистрация</h1>
        <p className="auth-subtitle">
          Создай аккаунт и сразу перейди в свою финансовую панель.
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
              placeholder="Минимум 8 символов"
              autoComplete="new-password"
            />
          </label>

          <label className="field">
            <span className="field-label">Повтори пароль</span>
            <input
              type="password"
              value={passwordRepeat}
              onChange={(e) => setPasswordRepeat(e.target.value)}
              placeholder="Повтори пароль"
              autoComplete="new-password"
            />
          </label>

          <div className="form-actions">
            <button
              className="btn btn-primary btn-block"
              type="submit"
              disabled={loading}
            >
              {loading ? "Создаём аккаунт..." : "Зарегистрироваться"}
            </button>
          </div>

          {error && <div className="alert">{error}</div>}
        </form>

        <div style={{ marginTop: 18 }} className="helper-text">
          Уже есть аккаунт?{" "}
          <Link to="/login" style={{ color: "#a5b4fc", fontWeight: 700 }}>
            Войти
          </Link>
        </div>
      </div>
    </div>
  );
}
