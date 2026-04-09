"use client";

import Link from "next/link";
import { ChangeEvent, FormEvent, useEffect, useMemo, useState } from "react";
import { ThemeToggle } from "@/components/theme-toggle";
import { deleteFeedback, getFeedbackList, getSummary, login, reanalyzeFeedback, updateStatus } from "@/lib/api";
import { Feedback } from "@/lib/types";

function sentimentBadge(sentiment?: string) {
  const value = sentiment || "Neutral";
  return <span className={`badge ${value}`}>{value}</span>;
}

export default function AdminPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [loginMessage, setLoginMessage] = useState("");

  const [items, setItems] = useState<Feedback[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [category, setCategory] = useState("");
  const [status, setStatus] = useState("");
  const [sortBy, setSortBy] = useState("date");
  const [search, setSearch] = useState("");
  const [summary, setSummary] = useState("");
  const [stats, setStats] = useState({ totalFeedback: 0, openItems: 0, averagePriority: 0, mostCommonTag: "-" });

  const maxPage = useMemo(() => Math.max(1, Math.ceil(total / 10)), [total]);

  useEffect(() => {
    const saved = localStorage.getItem("adminToken") || "";
    setToken(saved);
  }, []);

  const load = async (currentToken: string) => {
    if (!currentToken) {
      return;
    }

    const params = new URLSearchParams({
      page: String(page),
      limit: "10",
      category,
      status,
      sortBy,
      search,
    });

    try {
      const response = await getFeedbackList(currentToken, params);
      const data = response.data;
      if (!data) {
        return;
      }
      setItems(data.items || []);
      setTotal(data.total || 0);
      setStats(
        data.stats || {
          totalFeedback: 0,
          openItems: 0,
          averagePriority: 0,
          mostCommonTag: "-",
        }
      );
    } catch (error) {
      setLoginMessage(error instanceof Error ? error.message : "Failed to load dashboard.");
    }
  };

  useEffect(() => {
    void load(token);
  }, [token, page, category, status, sortBy, search]);

  const onLogin = async (e: FormEvent) => {
    e.preventDefault();
    setLoginMessage("Signing in...");
    try {
      const jwt = await login(email, password);
      localStorage.setItem("adminToken", jwt);
      setToken(jwt);
      setLoginMessage("Login successful.");
    } catch (error) {
      setLoginMessage(error instanceof Error ? error.message : "Login failed.");
    }
  };

  const onStatusChange = async (id: string, nextStatus: string) => {
    await updateStatus(token, id, nextStatus);
    await load(token);
  };

  const onDelete = async (id: string) => {
    await deleteFeedback(token, id);
    await load(token);
  };

  const onReanalyze = async (id: string) => {
    await reanalyzeFeedback(token, id);
    await load(token);
  };

  const loadWeeklySummary = async () => {
    try {
      const response = await getSummary(token);
      setSummary(response.data?.summary || "No summary available.");
    } catch (error) {
      setSummary(error instanceof Error ? error.message : "Failed to load summary.");
    }
  };

  const authReady = token.length > 0;

  return (
    <>
      <header className="topbar">
        <div>
          <h1>FeedPulse Admin</h1>
          <p>Feedback intelligence dashboard</p>
        </div>
        <div className="top-actions">
          <ThemeToggle />
          <Link className="button secondary" href="/">
            Back to Home
          </Link>
        </div>
      </header>

      <main className="container">
        {!authReady ? (
          <section className="card">
            <h2>Admin Login</h2>
            <form onSubmit={onLogin} className="grid-two">
              <label>
                Email
                <input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} />
              </label>
              <label>
                Password
                <input type="password" required value={password} onChange={(e) => setPassword(e.target.value)} />
              </label>
              <button type="submit">Login</button>
            </form>
            <p className="message">{loginMessage}</p>
          </section>
        ) : (
          <>
            <div className="stats-grid">
              <div className="stat">
                <span>Total Feedback</span>
                <b>{stats.totalFeedback}</b>
              </div>
              <div className="stat">
                <span>Open Items</span>
                <b>{stats.openItems}</b>
              </div>
              <div className="stat">
                <span>Avg Priority</span>
                <b>{Number(stats.averagePriority || 0).toFixed(1)}</b>
              </div>
              <div className="stat">
                <span>Most Common Tag</span>
                <b>{stats.mostCommonTag || "-"}</b>
              </div>
            </div>

            <section className="card">
              <div className="toolbar">
                <select value={category} onChange={(e) => setCategory(e.target.value)}>
                  <option value="">All Categories</option>
                  <option>Bug</option>
                  <option>Feature Request</option>
                  <option>Improvement</option>
                  <option>Other</option>
                </select>
                <select value={status} onChange={(e) => setStatus(e.target.value)}>
                  <option value="">All Statuses</option>
                  <option>New</option>
                  <option>In Review</option>
                  <option>Resolved</option>
                </select>
                <select value={sortBy} onChange={(e) => setSortBy(e.target.value)}>
                  <option value="date">Sort by Date</option>
                  <option value="priority">Sort by Priority</option>
                  <option value="sentiment">Sort by Sentiment</option>
                </select>
                <input value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search title or AI summary" />
                <button className="secondary" onClick={() => void load(token)} type="button">
                  Refresh
                </button>
                <button className="secondary" onClick={loadWeeklySummary} type="button">
                  7-day AI Summary
                </button>
              </div>

              <p className="summary-box">{summary}</p>

              <div className="table-wrapper">
                <table>
                  <thead>
                    <tr>
                      <th>Title</th>
                      <th>Category</th>
                      <th>Sentiment</th>
                      <th>Priority</th>
                      <th>Status</th>
                      <th>Date</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {items.map((item) => (
                      <tr key={item.id}>
                        <td>{item.title}</td>
                        <td>{item.category}</td>
                        <td>{sentimentBadge(item.ai_sentiment)}</td>
                        <td>{item.ai_priority || "-"}</td>
                        <td>
                          <select
                            value={item.status}
                            onChange={(event: ChangeEvent<HTMLSelectElement>) =>
                              void onStatusChange(item.id, event.target.value)
                            }
                          >
                            <option>New</option>
                            <option>In Review</option>
                            <option>Resolved</option>
                          </select>
                        </td>
                        <td>{new Date(item.createdAt).toLocaleDateString()}</td>
                        <td>
                          <div className="toolbar">
                            <button className="secondary" type="button" onClick={() => void onReanalyze(item.id)}>
                              Re-run AI
                            </button>
                            <button className="secondary" type="button" onClick={() => void onDelete(item.id)}>
                              Delete
                            </button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="toolbar">
                <button
                  className="secondary"
                  type="button"
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                  disabled={page <= 1}
                >
                  Previous
                </button>
                <span>
                  Page {page} of {maxPage}
                </span>
                <button
                  className="secondary"
                  type="button"
                  onClick={() => setPage((p) => Math.min(maxPage, p + 1))}
                  disabled={page >= maxPage}
                >
                  Next
                </button>
              </div>
            </section>
          </>
        )}
      </main>
    </>
  );
}

