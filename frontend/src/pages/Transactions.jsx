import { useEffect, useMemo, useState } from "react";
import { getCategories } from "../api/categories";
import { createTransaction, getTransactions } from "../api/transactions";

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

export default function Transactions() {
  const [loading, setLoading] = useState(false);
  const [loadingCats, setLoadingCats] = useState(false);
  const [error, setError] = useState("");

  const [categories, setCategories] = useState([]);

  // filters for list
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [q, setQ] = useState("");
  const [categoryId, setCategoryId] = useState("");

  const [items, setItems] = useState([]);
  const [total, setTotal] = useState(0);

  // create form
  const [occurredAt, setOccurredAt] = useState(todayISO());
  const [type, setType] = useState("expense");
  const [amountCents, setAmountCents] = useState("1299");
  const [merchant, setMerchant] = useState("Starbucks");
  const [note, setNote] = useState("латте");
  const [createCategoryId, setCreateCategoryId] = useState("");

  const categoriesByType = useMemo(() => {
    return categories.filter((c) => c.type === type);
  }, [categories, type]);

  useEffect(() => {
    // load categories once
    async function loadCats() {
      setLoadingCats(true);
      setError("");
      try {
        const data = await getCategories();
        setCategories(data.items || []);
      } catch (e) {
        setError(e.message || "Ошибка загрузки категорий");
      } finally {
        setLoadingCats(false);
      }
    }
    loadCats();
  }, []);

  useEffect(() => {
    // when type changes, adjust selected category for create form
    const hasSelected = categoriesByType.some((c) => c.id === createCategoryId);
    if (!hasSelected) {
      setCreateCategoryId(categoriesByType[0]?.id || "");
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [type, categoriesByType.length]);

  async function loadTransactions() {
    setLoading(true);
    setError("");
    try {
      const params = {};
      if (from) params.from = from;
      if (to) params.to = to;
      if (q) params.q = q;
      if (categoryId) params.category_id = categoryId;

      const data = await getTransactions(params);
      setItems(data.items || []);
      setTotal(data.total || 0);
    } catch (e) {
      setError(e.message || "Ошибка загрузки транзакций");
    } finally {
      setLoading(false);
    }
  }

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
    // после setState сразу перезагрузим "без фильтров"
    setTimeout(() => loadTransactions(), 0);
  }

  async function handleCreate(e) {
    e.preventDefault();
    setError("");

    if (!occurredAt) {
      setError("Укажи дату (occurred_at)");
      return;
    }
    if (!merchant.trim()) {
      setError("Укажи merchant");
      return;
    }
    if (!createCategoryId) {
      setError("Выбери категорию");
      return;
    }
    const cents = toIntCents(amountCents);
    if (cents === null) {
      setError("amount_cents должен быть целым числом >= 0");
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
        note: note ? note : "",
      });

      // очистим только некоторые поля (как удобнее)
      setMerchant("");
      setNote("");

      // обновим список с текущими фильтрами
      await loadTransactions();
    } catch (e) {
      setError(e.message || "Ошибка создания транзакции");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div style={{ padding: 24, display: "grid", gap: 16 }}>
      <h2 style={{ margin: 0 }}>Транзакции</h2>

      {error && (
        <div style={{ color: "crimson", padding: 12, border: "1px solid crimson", borderRadius: 8 }}>
          {error}
        </div>
      )}

      {/* Filters */}
      <div style={{ padding: 16, border: "1px solid #444", borderRadius: 12 }}>
        <div style={{ fontWeight: 600, marginBottom: 10 }}>Фильтры</div>

        <form onSubmit={handleApplyFilters} style={{ display: "grid", gap: 10 }}>
          <div style={{ display: "flex", gap: 10, flexWrap: "wrap" }}>
            <label style={{ display: "grid", gap: 6 }}>
              <span>from</span>
              <input type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
            </label>

            <label style={{ display: "grid", gap: 6 }}>
              <span>to</span>
              <input type="date" value={to} onChange={(e) => setTo(e.target.value)} />
            </label>

            <label style={{ display: "grid", gap: 6, minWidth: 220 }}>
              <span>поиск (q)</span>
              <input value={q} onChange={(e) => setQ(e.target.value)} placeholder="Starbucks" />
            </label>

            <label style={{ display: "grid", gap: 6, minWidth: 260 }}>
              <span>категория</span>
              <select value={categoryId} onChange={(e) => setCategoryId(e.target.value)} disabled={loadingCats}>
                <option value="">(все)</option>
                {categories.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.type === "expense" ? "расход" : "доход"} · {c.name}
                  </option>
                ))}
              </select>
            </label>
          </div>

          <div style={{ display: "flex", gap: 10 }}>
            <button type="submit" disabled={loading}>
              Применить
            </button>
            <button type="button" onClick={handleClearFilters} disabled={loading}>
              Сбросить
            </button>
          </div>
        </form>
      </div>

      {/* Create */}
      <div style={{ padding: 16, border: "1px solid #444", borderRadius: 12 }}>
        <div style={{ fontWeight: 600, marginBottom: 10 }}>Добавить транзакцию</div>

        <form onSubmit={handleCreate} style={{ display: "grid", gap: 10 }}>
          <div style={{ display: "flex", gap: 10, flexWrap: "wrap" }}>
            <label style={{ display: "grid", gap: 6 }}>
              <span>date</span>
              <input type="date" value={occurredAt} onChange={(e) => setOccurredAt(e.target.value)} />
            </label>

            <label style={{ display: "grid", gap: 6 }}>
              <span>type</span>
              <select value={type} onChange={(e) => setType(e.target.value)}>
                <option value="expense">expense</option>
                <option value="income">income</option>
              </select>
            </label>

            <label style={{ display: "grid", gap: 6 }}>
              <span>amount_cents</span>
              <input value={amountCents} onChange={(e) => setAmountCents(e.target.value)} />
            </label>

            <label style={{ display: "grid", gap: 6, minWidth: 260 }}>
              <span>category</span>
              <select value={createCategoryId} onChange={(e) => setCreateCategoryId(e.target.value)} disabled={loadingCats}>
                {categoriesByType.length === 0 && <option value="">(нет категорий)</option>}
                {categoriesByType.map((c) => (
                  <option key={c.id} value={c.id}>
                    {c.name}
                  </option>
                ))}
              </select>
            </label>

            <label style={{ display: "grid", gap: 6, minWidth: 240 }}>
              <span>merchant</span>
              <input value={merchant} onChange={(e) => setMerchant(e.target.value)} placeholder="Starbucks" />
            </label>

            <label style={{ display: "grid", gap: 6, minWidth: 240 }}>
              <span>note</span>
              <input value={note} onChange={(e) => setNote(e.target.value)} placeholder="опционально" />
            </label>
          </div>

          <div>
            <button type="submit" disabled={loading || loadingCats || !createCategoryId}>
              Добавить
            </button>
          </div>
        </form>
      </div>

      {/* Table */}
      <div style={{ padding: 16, border: "1px solid #444", borderRadius: 12 }}>
        <div style={{ display: "flex", alignItems: "baseline", gap: 10 }}>
          <div style={{ fontWeight: 600 }}>Список</div>
          <div style={{ opacity: 0.8 }}>total: {total}</div>
          <div style={{ flex: 1 }} />
          <button onClick={loadTransactions} disabled={loading}>
            Обновить
          </button>
        </div>

        <div style={{ marginTop: 12, overflowX: "auto" }}>
          <table style={{ width: "100%", borderCollapse: "collapse" }}>
            <thead>
              <tr>
                <th style={th}>date</th>
                <th style={th}>type</th>
                <th style={th}>amount_cents</th>
                <th style={th}>merchant</th>
                <th style={th}>note</th>
                <th style={th}>category_id</th>
              </tr>
            </thead>
            <tbody>
              {items.map((t) => (
                <tr key={t.id}>
                  <td style={td}>{t.occurred_at}</td>
                  <td style={td}>{t.type}</td>
                  <td style={td}>{t.amount_cents}</td>
                  <td style={td}>{t.merchant}</td>
                  <td style={td}>{t.note || ""}</td>
                  <td style={tdMono}>{t.category_id}</td>
                </tr>
              ))}
              {items.length === 0 && (
                <tr>
                  <td style={td} colSpan={6}>
                    {loading ? "Загрузка..." : "Пока нет транзакций"}
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div style={{ marginTop: 10, opacity: 0.8, fontSize: 12 }}>
          Подсказка: можно фильтровать по <b>q</b>, <b>from/to</b>, <b>category</b>.
        </div>
      </div>
    </div>
  );
}

const th = {
  textAlign: "left",
  padding: "8px 10px",
  borderBottom: "1px solid #555",
  fontWeight: 600,
};

const td = {
  padding: "8px 10px",
  borderBottom: "1px solid #333",
  verticalAlign: "top",
};

const tdMono = {
  ...td,
  fontFamily: "ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace",
  fontSize: 12,
};