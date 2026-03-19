import { useEffect, useMemo, useState } from "react";
import { getCategories } from "../api/categories";
import {
  createTransaction,
  deleteTransactionRequest,
  getTransactions,
} from "../api/transactions";

function todayISO() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function toIntCents(v) {
  const n = Number(v);
  if (!Number.isFinite(n)) return null;
  if (!Number.isInteger(n)) return null;
  if (n < 0) return null;
  return n;
}

function formatMoney(cents) {
  const amount = Number(cents || 0) / 100;
  return amount.toLocaleString("ru-RU", {
    style: "currency",
    currency: "RUB",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function formatType(type) {
  return type === "income" ? "Доход" : "Расход";
}

export default function Transactions() {
  const [loading, setLoading] = useState(false);
  const [loadingCats, setLoadingCats] = useState(false);
  const [error, setError] = useState("");

  const [categories, setCategories] = useState([]);

  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [q, setQ] = useState("");
  const [categoryId, setCategoryId] = useState("");

  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);

  const [occurredAt, setOccurredAt] = useState(todayISO());
  const [type, setType] = useState("expense");
  const [amountCents, setAmountCents] = useState("");
  const [merchant, setMerchant] = useState("");
  const [note, setNote] = useState("");
  const [createCategoryId, setCreateCategoryId] = useState("");

  const categoriesByType = useMemo(() => {
    return categories.filter((c) => c.type === type);
  }, [categories, type]);

  const categoryMap = useMemo(() => {
    const map = new Map();
    categories.forEach((c) => {
      map.set(c.id, c);
    });
    return map;
  }, [categories]);

  function buildFilters(next = {}) {
    const finalFrom = next.from ?? from;
    const finalTo = next.to ?? to;
    const finalQ = next.q ?? q;
    const finalCategoryId = next.categoryId ?? categoryId;

    const params = {};
    if (finalFrom) params.from = finalFrom;
    if (finalTo) params.to = finalTo;
    if (finalQ) params.q = finalQ;
    if (finalCategoryId) params.category_id = finalCategoryId;

    return params;
  }

  async function loadTransactions(params = buildFilters()) {
    setLoading(true);
    setError("");

    try {
      const data = await getTransactions(params);
      setItems(data.items || []);
      setTotal(data.total || 0);
    } catch (e) {
      setError(e.message || "Не удалось загрузить транзакции");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    async function loadCats() {
      setLoadingCats(true);
      setError("");

      try {
        const data = await getCategories();
        setCategories(data.items || []);
      } catch (e) {
        setError(e.message || "Не удалось загрузить категории");
      } finally {
        setLoadingCats(false);
      }
    }

    loadCats();
  }, []);

  useEffect(() => {
    const hasSelected = categoriesByType.some((c) => c.id === createCategoryId);
    if (!hasSelected) {
      setCreateCategoryId(categoriesByType[0]?.id || "");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [type, categoriesByType.length]);

  useEffect(() => {
    loadTransactions();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleApplyFilters(e) {
    e.preventDefault();
    await loadTransactions();
  }

  async function handleClearFilters() {
    setFrom("");
    setTo("");
    setQ("");
    setCategoryId("");
    await loadTransactions({});
  }

  async function handleCreate(e) {
    e.preventDefault();
    setError("");

    if (!occurredAt) {
      setError("Укажи дату операции");
      return;
    }

    if (!merchant.trim()) {
      setError("Укажи место или источник операции");
      return;
    }

    if (!createCategoryId) {
      setError("Выбери категорию");
      return;
    }

    const cents = toIntCents(amountCents);
    if (cents === null) {
      setError("Сумма должна быть целым числом в копейках");
      return;
    }

    setLoading(true);

    try {
      await createTransaction({
        occurred_at: occurredAt,
        type,
        amount_cents: cents,
        category_id: createCategoryId,
        merchant: merchant.trim(),
        note: note ? note.trim() : "",
      });

      setAmountCents("");
      setMerchant("");
      setNote("");

      await loadTransactions();
    } catch (e) {
      setError(e.message || "Не удалось создать транзакцию");
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id) {
    setError("");
    setLoading(true);

    try {
      await deleteTransactionRequest(id);
      await loadTransactions();
    } catch (e) {
      setError(e.message || "Не удалось удалить транзакцию");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <div className="page-title-block">
          <h1 className="page-title">Транзакции</h1>
          <div className="page-subtitle">
            Добавляй операции, фильтруй список и отслеживай движения денег.
          </div>
        </div>

        <div className="page-actions">
          <div className="kpi-inline">
            <span>Всего записей:</span>
            <strong>{total}</strong>
          </div>

          <button
            className="btn btn-secondary"
            onClick={() => loadTransactions()}
            disabled={loading}
          >
            {loading ? "Обновляем..." : "Обновить"}
          </button>
        </div>
      </div>

      {error && <div className="alert">{error}</div>}

      <section className="surface-card section-card">
        <div className="section-head">
          <div>
            <h2 className="section-title">Фильтры</h2>
            <div className="section-subtitle">
              Можно отбирать записи по датам, поисковой фразе и категории.
            </div>
          </div>
        </div>

        <form className="form-grid" onSubmit={handleApplyFilters}>
          <div className="form-grid columns-4">
            <label className="field">
              <span className="field-label">Дата от</span>
              <input
                type="date"
                value={from}
                onChange={(e) => setFrom(e.target.value)}
              />
            </label>

            <label className="field">
              <span className="field-label">Дата до</span>
              <input
                type="date"
                value={to}
                onChange={(e) => setTo(e.target.value)}
              />
            </label>

            <label className="field">
              <span className="field-label">Поиск</span>
              <input
                value={q}
                onChange={(e) => setQ(e.target.value)}
                placeholder="Например: Starbucks"
              />
            </label>

            <label className="field">
              <span className="field-label">Категория</span>
              <select
                value={categoryId}
                onChange={(e) => setCategoryId(e.target.value)}
                disabled={loadingCats}
              >
                <option value="">Все категории</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.type === "expense" ? "Расход" : "Доход"} · {c.name}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div className="form-actions">
            <button
              className="btn btn-primary"
              type="submit"
              disabled={loading}
            >
              Применить
            </button>
            <button
              className="btn btn-ghost"
              type="button"
              onClick={handleClearFilters}
              disabled={loading}
            >
              Сбросить
            </button>
          </div>
        </form>
      </section>

      <section className="surface-card section-card">
        <div className="section-head">
          <div>
            <h2 className="section-title">Новая транзакция</h2>
            <div className="section-subtitle">
              Все суммы вводятся в копейках, чтобы не терять точность.
            </div>
          </div>
        </div>

        <form className="form-grid" onSubmit={handleCreate}>
          <div className="form-grid columns-3">
            <label className="field">
              <span className="field-label">Дата операции</span>
              <input
                type="date"
                value={occurredAt}
                onChange={(e) => setOccurredAt(e.target.value)}
              />
            </label>

            <label className="field">
              <span className="field-label">Тип</span>
              <select value={type} onChange={(e) => setType(e.target.value)}>
                <option value="expense">Расход</option>
                <option value="income">Доход</option>
              </select>
            </label>

            <label className="field">
              <span className="field-label">Сумма (в копейках)</span>
              <input
                value={amountCents}
                onChange={(e) => setAmountCents(e.target.value)}
                placeholder="Например: 1299"
              />
            </label>

            <label className="field">
              <span className="field-label">Категория</span>
              <select
                value={createCategoryId}
                onChange={(e) => setCreateCategoryId(e.target.value)}
                disabled={loadingCats}
              >
                {categoriesByType.length === 0 && (
                  <option value="">Нет категорий</option>
                )}
                {categoriesByType.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            </label>

            <label className="field">
              <span className="field-label">Магазин / источник</span>
              <input
                value={merchant}
                onChange={(e) => setMerchant(e.target.value)}
                placeholder="Например: Starbucks"
              />
            </label>

            <label className="field">
              <span className="field-label">Комментарий</span>
              <input
                value={note}
                onChange={(e) => setNote(e.target.value)}
                placeholder="Необязательно"
              />
            </label>
          </div>

          <div className="form-actions">
            <button
              className="btn btn-primary"
              type="submit"
              disabled={loading || loadingCats || !createCategoryId}
            >
              Добавить транзакцию
            </button>
          </div>
        </form>
      </section>

      <section className="surface-card section-card">
        <div className="section-head">
          <div>
            <h2 className="section-title">Список операций</h2>
            <div className="section-subtitle">
              Здесь отображаются все транзакции по текущим фильтрам.
            </div>
          </div>
        </div>

        <div className="table-wrap">
          <table className="table">
            <thead>
              <tr>
                <th>Дата</th>
                <th>Тип</th>
                <th>Сумма</th>
                <th>Место / источник</th>
                <th>Комментарий</th>
                <th>Категория</th>
                <th>Действие</th>
              </tr>
            </thead>

            <tbody>
              {items.length > 0 ? (
                items.map((t) => {
                  const category = categoryMap.get(t.category_id);
                  const isIncome = t.type === "income";

                  return (
                    <tr key={t.id}>
                      <td>{t.occurred_at}</td>
                      <td>
                        <span
                          className={`type-badge ${isIncome ? "income" : "expense"}`}
                        >
                          {formatType(t.type)}
                        </span>
                      </td>
                      <td
                        className={
                          isIncome ? "amount-positive" : "amount-negative"
                        }
                      >
                        {formatMoney(t.amount_cents)}
                      </td>
                      <td>{t.merchant}</td>
                      <td>{t.note || "—"}</td>
                      <td>
                        <div>{category?.name || "Неизвестная категория"}</div>
                        <div className="mono">{t.category_id}</div>
                      </td>
                      <td>
                        <button
                          className="btn btn-danger"
                          type="button"
                          onClick={() => handleDelete(t.id)}
                          disabled={loading}
                        >
                          Удалить
                        </button>
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td className="table-empty" colSpan={7}>
                    {loading ? "Загрузка..." : "Пока нет транзакций"}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="helper-text" style={{ marginTop: 12 }}>
          Подсказка: фильтрация поддерживает даты, строку поиска и категорию.
        </div>
      </section>
    </div>
  );
}
