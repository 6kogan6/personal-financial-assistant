import {
  BarElement,
  CategoryScale,
  Chart as ChartJS,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
} from "chart.js";
import { useEffect, useMemo, useState } from "react";
import { Bar, Line } from "react-chartjs-2";
import { getSummary } from "../api/reports";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Tooltip,
  Legend,
);

function monthNow() {
  const d = new Date();
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  return `${y}-${m}`;
}

function formatMoney(cents) {
  const n = Number(cents || 0) / 100;
  return n.toLocaleString("ru-RU", {
    style: "currency",
    currency: "RUB",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}

function formatType(type) {
  return type === "income" ? "Доход" : "Расход";
}

const commonChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: "top",
      labels: {
        color: "#cbd5e1",
        boxWidth: 14,
        usePointStyle: true,
      },
    },
  },
  scales: {
    x: {
      ticks: { color: "#94a3b8" },
      grid: { color: "rgba(148, 163, 184, 0.08)" },
    },
    y: {
      ticks: { color: "#94a3b8" },
      grid: { color: "rgba(148, 163, 184, 0.08)" },
    },
  },
};

export default function Dashboard() {
  const [month, setMonth] = useState(monthNow());
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [data, setData] = useState(null);

  async function load() {
    setLoading(true);
    setError("");

    try {
      const summary = await getSummary(month);
      setData(summary);
    } catch (e) {
      setError(e.message || "Не удалось загрузить отчёт");
      setData(null);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [month]);

  const totals = data?.totals || {
    income_cents: 0,
    expense_cents: 0,
    balance_cents: 0,
  };

  const dailyChart = useMemo(() => {
    const daily = data?.daily || [];
    const labels = daily.map((x) => x.date);
    const income = daily.map((x) => (x.income_cents || 0) / 100);
    const expense = daily.map((x) => (x.expense_cents || 0) / 100);

    return {
      labels,
      datasets: [
        {
          label: "Доход",
          data: income,
          borderColor: "#10b981",
          backgroundColor: "rgba(16, 185, 129, 0.22)",
          tension: 0.35,
          pointRadius: 3,
          pointHoverRadius: 5,
          fill: false,
        },
        {
          label: "Расход",
          data: expense,
          borderColor: "#f43f5e",
          backgroundColor: "rgba(244, 63, 94, 0.22)",
          tension: 0.35,
          pointRadius: 3,
          pointHoverRadius: 5,
          fill: false,
        },
      ],
    };
  }, [data]);

  const categoriesChart = useMemo(() => {
    const by = data?.by_category || [];

    const expense = by.filter((x) => x.type === "expense");
    const income = by.filter((x) => x.type === "income");

    return {
      expense: {
        labels: expense.map((x) => x.category_name),
        datasets: [
          {
            label: "Расходы по категориям",
            data: expense.map((x) => (x.amount_cents || 0) / 100),
            backgroundColor: [
              "#6366f1",
              "#8b5cf6",
              "#ec4899",
              "#f97316",
              "#eab308",
              "#06b6d4",
              "#10b981",
            ],
            borderRadius: 12,
          },
        ],
      },
      income: {
        labels: income.map((x) => x.category_name),
        datasets: [
          {
            label: "Доходы по категориям",
            data: income.map((x) => (x.amount_cents || 0) / 100),
            backgroundColor: [
              "#10b981",
              "#22c55e",
              "#14b8a6",
              "#06b6d4",
              "#3b82f6",
            ],
            borderRadius: 12,
          },
        ],
      },
      hasIncome: income.length > 0,
      hasExpense: expense.length > 0,
    };
  }, [data]);

  return (
    <div className="page">
      <div className="page-header">
        <div className="page-title-block">
          <h1 className="page-title">Дашборд</h1>
          <div className="page-subtitle">
            Смотри ключевые показатели, динамику по дням и разбивку по
            категориям.
          </div>
        </div>

        <div className="page-actions">
          <label className="field" style={{ minWidth: 180 }}>
            <span className="field-label">Месяц</span>
            <input
              type="month"
              value={month}
              onChange={(e) => setMonth(e.target.value)}
            />
          </label>

          <button className="btn btn-primary" onClick={load} disabled={loading}>
            {loading ? "Обновляем..." : "Обновить"}
          </button>
        </div>
      </div>

      {error && <div className="alert">{error}</div>}

      <div className="stat-grid">
        <StatCard
          variant="income"
          title="Доходы"
          value={formatMoney(totals.income_cents)}
          hint="Все поступления за выбранный месяц"
        />
        <StatCard
          variant="expense"
          title="Расходы"
          value={formatMoney(totals.expense_cents)}
          hint="Все списания за выбранный месяц"
        />
        <StatCard
          variant="balance"
          title="Баланс"
          value={formatMoney(totals.balance_cents)}
          hint="Разница между доходами и расходами"
        />
      </div>

      <div className="grid-2">
        <section className="surface-card section-card chart-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Динамика по дням</h2>
              <div className="section-subtitle">
                Доходы и расходы внутри месяца
              </div>
            </div>
          </div>

          <div className="chart-wrap">
            <Line data={dailyChart} options={commonChartOptions} />
          </div>
        </section>

        <section className="surface-card section-card chart-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Расходы по категориям</h2>
              <div className="section-subtitle">
                Куда уходит больше всего денег
              </div>
            </div>
          </div>

          {categoriesChart.hasExpense ? (
            <div className="chart-wrap">
              <Bar
                data={categoriesChart.expense}
                options={commonChartOptions}
              />
            </div>
          ) : (
            <div className="helper-text">За этот месяц пока нет расходов.</div>
          )}
        </section>
      </div>

      {categoriesChart.hasIncome && (
        <section className="surface-card section-card chart-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Доходы по категориям</h2>
              <div className="section-subtitle">
                Из каких источников приходят деньги
              </div>
            </div>
          </div>

          <div className="chart-wrap">
            <Bar data={categoriesChart.income} options={commonChartOptions} />
          </div>
        </section>
      )}

      <div className="grid-2">
        <section className="surface-card section-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Разбивка по категориям</h2>
              <div className="section-subtitle">
                Суммы по каждому типу категории
              </div>
            </div>
          </div>

          <SimpleTable
            headers={["Тип", "Категория", "Сумма"]}
            rows={(data?.by_category || []).map((x) => [
              <span
                className={`type-badge ${x.type === "income" ? "income" : "expense"}`}
                key={`type-${x.category_name}`}
              >
                {formatType(x.type)}
              </span>,
              x.category_name,
              formatMoney(x.amount_cents),
            ])}
          />
        </section>

        <section className="surface-card section-card">
          <div className="section-head">
            <div>
              <h2 className="section-title">Ежедневные значения</h2>
              <div className="section-subtitle">
                Доходы и расходы по каждому дню
              </div>
            </div>
          </div>

          <SimpleTable
            headers={["Дата", "Доход", "Расход"]}
            rows={(data?.daily || []).map((x) => [
              x.date,
              formatMoney(x.income_cents),
              formatMoney(x.expense_cents),
            ])}
          />
        </section>
      </div>
    </div>
  );
}

function StatCard({ title, value, hint, variant }) {
  return (
    <div className={`stat-card ${variant}`}>
      <div className="stat-label">{title}</div>
      <div className="stat-value">{value}</div>
      <div className="stat-hint">{hint}</div>
    </div>
  );
}

function SimpleTable({ headers, rows }) {
  return (
    <div className="table-wrap">
      <table className="table">
        <thead>
          <tr>
            {headers.map((header) => (
              <th key={header}>{header}</th>
            ))}
          </tr>
        </thead>

        <tbody>
          {rows.length > 0 ? (
            rows.map((row, idx) => (
              <tr key={idx}>
                {row.map((cell, cellIndex) => (
                  <td key={cellIndex}>{cell}</td>
                ))}
              </tr>
            ))
          ) : (
            <tr>
              <td className="table-empty" colSpan={headers.length}>
                Нет данных за выбранный период
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}
