import { useEffect, useState } from "react";
import { getProfile } from "../api/profile";

function formatDate(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return date.toLocaleString("ru-RU", {
    dateStyle: "long",
    timeStyle: "short",
  });
}

export default function Profile() {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    async function load() {
      setLoading(true);
      setError("");

      try {
        const data = await getProfile();
        setProfile(data);
      } catch (e) {
        setError(e.message || "Не удалось загрузить профиль");
      } finally {
        setLoading(false);
      }
    }

    load();
  }, []);

  return (
    <div className="page">
      <div className="page-header">
        <div className="page-title-block">
          <h1 className="page-title">Профиль</h1>
          <div className="page-subtitle">
            Основная информация о твоём аккаунте и системе категорий.
          </div>
        </div>
      </div>

      {error && <div className="alert">{error}</div>}

      <div className="grid-2">
        <section className="surface-card section-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Данные аккаунта</h2>
              <div className="section-subtitle">
                Здесь отображаются данные текущего пользователя.
              </div>
            </div>
          </div>

          {loading ? (
            <div className="helper-text">Загрузка профиля...</div>
          ) : profile ? (
            <div className="form-grid">
              <div className="field">
                <span className="field-label">ID пользователя</span>
                <div className="mono">{profile.id}</div>
              </div>

              <div className="field">
                <span className="field-label">Электронная почта</span>
                <div>{profile.email}</div>
              </div>

              <div className="field">
                <span className="field-label">Дата регистрации</span>
                <div>{formatDate(profile.created_at)}</div>
              </div>
            </div>
          ) : (
            <div className="helper-text">Нет данных профиля</div>
          )}
        </section>

        <section className="surface-card section-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Системные категории</h2>
              <div className="section-subtitle">
                Они автоматически доступны при создании транзакций.
              </div>
            </div>
          </div>

          <div className="form-grid">
            <div>
              <div className="field-label" style={{ marginBottom: 10 }}>
                Доходы
              </div>
              <div className="form-actions">
                <span className="type-badge income">Зарплата</span>
                <span className="type-badge income">Инвестиции</span>
                <span className="type-badge income">Фриланс</span>
                <span className="type-badge income">Подарки</span>
              </div>
            </div>

            <div>
              <div className="field-label" style={{ marginBottom: 10 }}>
                Расходы
              </div>
              <div className="form-actions">
                <span className="type-badge expense">Продукты</span>
                <span className="type-badge expense">Транспорт</span>
                <span className="type-badge expense">Игры</span>
                <span className="type-badge expense">Дом</span>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  );
}
