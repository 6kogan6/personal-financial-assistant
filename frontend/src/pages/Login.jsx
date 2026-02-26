import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { apiFetch } from "../api/http";
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
      // логин
      const data = await apiFetch("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });

      // сохраняем токен
      setToken(data.access_token);

      // редирект
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(err?.message || "Ошибка входа");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ padding: 24, maxWidth: 420 }}>
      <h2>Вход</h2>

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 12 }}>
          <div>Email</div>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            style={{ width: "100%", padding: 8 }}
            placeholder="email"
          />
        </div>

        <div style={{ marginBottom: 12 }}>
          <div>Password</div>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            style={{ width: "100%", padding: 8 }}
            placeholder="password"
          />
        </div>

        <button type="submit" disabled={loading}>
          {loading ? "Входим..." : "Войти"}
        </button>

        {error && <div style={{ color: "crimson", marginTop: 12 }}>{error}</div>}
      </form>
    </div>
  );
}