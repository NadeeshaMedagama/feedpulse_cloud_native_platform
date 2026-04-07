"use client";

import Link from "next/link";
import { FormEvent, useMemo, useState } from "react";
import { ThemeToggle } from "@/components/theme-toggle";
import { submitFeedback } from "@/lib/api";

export default function HomePage() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [category, setCategory] = useState("Bug");
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  const count = useMemo(() => description.length, [description]);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!title.trim() || description.trim().length < 20) {
      setMessage("Title is required and description must be at least 20 characters.");
      return;
    }

    setLoading(true);
    setMessage("Submitting...");
    try {
      await submitFeedback({
        title: title.trim(),
        description: description.trim(),
        category,
        name: name.trim(),
        email: email.trim(),
      });
      setMessage("Feedback submitted successfully.");
      setTitle("");
      setDescription("");
      setCategory("Bug");
      setName("");
      setEmail("");
    } catch (error) {
      setMessage(error instanceof Error ? error.message : "Submission failed.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <header className="topbar">
        <div>
          <h1>FeedPulse_Cloud_Native_Platform</h1>
          <p>AI-powered product feedback platform</p>
        </div>
        <div className="top-actions">
          <ThemeToggle />
          <Link className="button" href="/admin">
            Admin Login
          </Link>
        </div>
      </header>

      <main className="container">
        <section className="card">
          <h2>Submit Feedback</h2>
          <p>Share product issues, requests, and ideas. AI will auto-analyze each submission.</p>

          <form onSubmit={onSubmit}>
            <label>
              Title
              <input value={title} onChange={(e) => setTitle(e.target.value)} maxLength={120} required />
            </label>

            <label>
              Description
              <textarea value={description} onChange={(e) => setDescription(e.target.value)} rows={6} minLength={20} required />
              <small>{count} characters (minimum 20)</small>
            </label>

            <label>
              Category
              <select value={category} onChange={(e) => setCategory(e.target.value)}>
                <option>Bug</option>
                <option>Feature Request</option>
                <option>Improvement</option>
                <option>Other</option>
              </select>
            </label>

            <div className="grid-two">
              <label>
                Name (optional)
                <input value={name} onChange={(e) => setName(e.target.value)} />
              </label>
              <label>
                Email (optional)
                <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
              </label>
            </div>

            <button type="submit" disabled={loading}>
              {loading ? "Submitting..." : "Submit Feedback"}
            </button>
          </form>
          <p className="message">{message}</p>
        </section>
      </main>
    </>
  );
}

