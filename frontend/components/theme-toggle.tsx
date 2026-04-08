"use client";

import { useEffect, useState } from "react";

export function ThemeToggle() {
  const [theme, setTheme] = useState("light");

  useEffect(() => {
    const saved = localStorage.getItem("theme") || "light";
    document.documentElement.setAttribute("data-theme", saved);
    setTheme(saved);
  }, []);

  const toggleTheme = () => {
    const next = theme === "light" ? "dark" : "light";
    document.documentElement.setAttribute("data-theme", next);
    localStorage.setItem("theme", next);
    setTheme(next);
  };

  return (
    <button className="icon-btn" type="button" onClick={toggleTheme} aria-label="Toggle theme" title="Toggle theme">
      <svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true">
        <path
          d="M21 12.79A9 9 0 1 1 11.21 3c.5 0 .79.58.5 1A7 7 0 1 0 20 12.29c.42-.29 1 .01 1 .5z"
          fill="currentColor"
        />
      </svg>
    </button>
  );
}

