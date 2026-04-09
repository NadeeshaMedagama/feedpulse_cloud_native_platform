import { APIResponse, Feedback, FeedbackListData, LoginData } from "./types";

async function parseResponse<T>(res: Response): Promise<APIResponse<T>> {
  const payload = (await res.json()) as APIResponse<T>;
  if (!res.ok || !payload.success) {
    throw new Error(payload.error || payload.message || "Request failed");
  }
  return payload;
}

export async function login(email: string, password: string): Promise<string> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  const payload = await parseResponse<LoginData>(res);
  return payload.data?.token || "";
}

export async function submitFeedback(input: {
  title: string;
  description: string;
  category: string;
  name?: string;
  email?: string;
}) {
  const res = await fetch("/api/feedback", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(input),
  });
  return parseResponse<Feedback>(res);
}

export async function getFeedbackList(token: string, params: URLSearchParams) {
  const res = await fetch(`/api/feedback?${params.toString()}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return parseResponse<FeedbackListData>(res);
}

export async function updateStatus(token: string, id: string, status: string) {
  const res = await fetch(`/api/feedback/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ status }),
  });
  return parseResponse<Record<string, never>>(res);
}

export async function deleteFeedback(token: string, id: string) {
  const res = await fetch(`/api/feedback/${id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
  });
  return parseResponse<Record<string, never>>(res);
}

export async function reanalyzeFeedback(token: string, id: string) {
  const res = await fetch(`/api/feedback/${id}/reanalyze`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
  });
  return parseResponse<Record<string, never>>(res);
}

export async function getSummary(token: string) {
  const res = await fetch("/api/feedback/summary", {
    headers: { Authorization: `Bearer ${token}` },
  });
  return parseResponse<{ summary: string }>(res);
}

