import { useEffect, useMemo, useState } from "react";
import { Line, Bar } from "react-chartjs-2";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Tooltip,
  Legend,
} from "chart.js";
import { getSummary } from "../api/reports";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Tooltip,
  Legend
);

function monthNow() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  return `${y}-${m}`;
}

function fmtCents(cents) {
  const n = Number(cents || 0) / 100;
  return n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

export default function Dashboard() {
  const [month, setMonth] = useState(monthNow());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [data, setData] = useState(null);

  async function load() {
    setLoading(true);
    setError("");
    try {
      const s = await getSummary(month);
      setData(s);
    } catch (e) {
      setError(e.message || "Ошибка загрузки отчёта");
      setData(null);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [month]);

  const totals = data?.totals || { income_cents: 0, expense_cents: 0, balance_cents: 0 };

  const dailyChart = useMemo(() => {
    const daily = data?.daily || [];
    const labels = daily.map((x) => x.date);
    const income = daily.map((x) => (x.income_cents || 0) / 100);
    const expense = daily.map((x) => (x.expense_cents || 0) / 100);

    return {
      labels,
      datasets: [
        {
          label: "Income",
          data: income,
        },
        {
          label: "Expense",
          data: expense,
        },
      ],
    };
  }, [data]);

  const categoriesChart = useMemo(() => {
    const by = data?.by_category || [];

    // Можно рисовать все категории вместе, но удобнее отдельно расходы/доходы.
    const expense = by.filter((x) => x.type === "expense");
    const income = by.filter((x) => x.type === "income");

    const labels = expense.map((x) => x.category_name);
    const values = expense.map((x) => (x.amount_cents || 0) / 100);

    const incomeLabels = income.map((x) => x.category_name);
    const incomeValues = income.map((x) => (x.amount_cents || 0) / 100);

    return {
      expense: {
        labels,
        datasets: [
          {
            label: "Expenses by category",
            data: values,
          },
        ],
      },
      income: {
        labels: incomeLabels,
        datasets: [
          {
            label: "Income by category",
            data: incomeValues,
          },
        ],
      },
      hasIncome: income.length > 0,
      hasExpense: expense.length > 0,
    };
  }, [data]);

  return (
    <div style={{ padding: 24, display: "grid", gap: 16 }}>
      <div style={{ display: "flex", alignItems: "baseline", gap: 12, flexWrap: "wrap" }}>
        <h2 style={{ margin: 0 }}>Дашборд</h2>

        <label style={{ display: "grid", gap: 6 }}>
          <span style={{ opacity: 0.85 }}>Месяц</span>
          <input
            value={month}
            onChange={(e) => setMonth(e.target.value)}
            placeholder="YYYY-MM"
            style={{ padding: 6, minWidth: 120 }}
          />
        </label>

        <button onClick={load} disabled={loading}>
          {loading ? "Загрузка..." : "Обновить"}
        </button>
      </div>

      {error && (
        <div style={{ color: "crimson", padding: 12, border: "1px solid crimson", borderRadius: 8 }}>
          {error}
        </div>
      )}

      {/* Totals */}
      <div style={{ display: "grid", gap: 12, gridTemplateColumns: "repeat(auto-fit, minmax(220px, 1fr))" }}>
        <Card title="Income" value={`${fmtCents(totals.income_cents)}`} />
        <Card title="Expense" value={`${fmtCents(totals.expense_cents)}`} />
        <Card title="Balance" value={`${fmtCents(totals.balance_cents)}`} />
      </div>

      {/* Charts */}
      <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(380px, 1fr))" }}>
        <div style={panel}>
          <div style={panelTitle}>По дням (Income/Expense)</div>
          <div style={{ height: 320 }}>
            <Line
              data={dailyChart}
              options={{
                responsive: true,
                maintainAspectRatio: false,
                plugins: { legend: { position: "top" } },
              }}
            />
          </div>
        </div>

        <div style={panel}>
          <div style={panelTitle}>Расходы по категориям</div>
          {categoriesChart.hasExpense ? (
            <div style={{ height: 320 }}>
              <Bar
                data={categoriesChart.expense}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: { legend: { position: "top" } },
                }}
              />
            </div>
          ) : (
            <div style={{ opacity: 0.8 }}>Нет расходов за этот месяц</div>
          )}
        </div>

        {categoriesChart.hasIncome && (
          <div style={panel}>
            <div style={panelTitle}>Доходы по категориям</div>
            <div style={{ height: 320 }}>
              <Bar
                data={categoriesChart.income}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: { legend: { position: "top" } },
                }}
              />
            </div>
          </div>
        )}
      </div>

      {/* Tables */}
      <div style={{ display: "grid", gap: 16, gridTemplateColumns: "repeat(auto-fit, minmax(380px, 1fr))" }}>
        <div style={panel}>
          <div style={panelTitle}>By category</div>
          <SimpleTable
            headers={["type", "category", "amount"]}
            rows={(data?.by_category || []).map((x) => [
              x.type,
              x.category_name,
              fmtCents(x.amount_cents),
            ])}
          />
        </div>

        <div style={panel}>
          <div style={panelTitle}>Daily</div>
          <SimpleTable
            headers={["date", "income", "expense"]}
            rows={(data?.daily || []).map((x) => [
              x.date,
              fmtCents(x.income_cents),
              fmtCents(x.expense_cents),
            ])}
          />
        </div>
      </div>
    </div>
  );
}

function Card({ title, value }) {
  return (
    <div style={{ padding: 16, border: "1px solid #444", borderRadius: 12 }}>
      <div style={{ opacity: 0.85, marginBottom: 8 }}>{title}</div>
      <div style={{ fontSize: 24, fontWeight: 700 }}>{value}</div>
    </div>
  );
}

function SimpleTable({ headers, rows }) {
  return (
    <div style={{ overflowX: "auto" }}>
      <table style={{ width: "100%", borderCollapse: "collapse" }}>
        <thead>
          <tr>
            {headers.map((h) => (
              <th key={h} style={th}>{h}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((r, idx) => (
            <tr key={idx}>
              {r.map((cell, j) => (
                <td key={j} style={td}>{cell}</td>
              ))}
            </tr>
          ))}
          {rows.length === 0 && (
            <tr>
              <td style={td} colSpan={headers.length}>Нет данных</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

const panel = { padding: 16, border: "1px solid #444", borderRadius: 12 };
const panelTitle = { fontWeight: 600, marginBottom: 10 };

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