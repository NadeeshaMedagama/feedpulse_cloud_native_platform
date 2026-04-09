export type APIResponse<T> = {
  success: boolean;
  data?: T;
  error?: string;
  message?: string;
};

export type Feedback = {
  id: string;
  title: string;
  description: string;
  category: string;
  status: string;
  ai_sentiment?: "Positive" | "Neutral" | "Negative";
  ai_priority?: number;
  ai_summary?: string;
  ai_tags?: string[];
  createdAt: string;
};

export type FeedbackListData = {
  items: Feedback[];
  total: number;
  page: number;
  limit: number;
  stats: {
    totalFeedback: number;
    openItems: number;
    averagePriority: number;
    mostCommonTag: string;
  };
};

export type LoginData = {
  token: string;
};

