import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import AdminPage from "@/app/admin/page";
import { getFeedbackList, login } from "@/lib/api";

jest.mock("@/lib/api", () => ({
  deleteFeedback: jest.fn(),
  getFeedbackList: jest.fn(),
  getSummary: jest.fn(),
  login: jest.fn(),
  reanalyzeFeedback: jest.fn(),
  updateStatus: jest.fn(),
}));

describe("AdminPage", () => {
  beforeEach(() => {
    localStorage.clear();
    jest.clearAllMocks();
  });

  it("shows login error message when credentials are invalid", async () => {
    (login as jest.Mock).mockRejectedValue(new Error("invalid email or password"));

    render(<AdminPage />);

    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: "admin@example.com" } });
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: "wrong-password" } });
    fireEvent.click(screen.getByRole("button", { name: /^login$/i }));

    expect(await screen.findByText(/invalid email or password/i)).toBeInTheDocument();
  });

  it("loads dashboard data when token exists in localStorage", async () => {
    localStorage.setItem("adminToken", "token-123");
    (getFeedbackList as jest.Mock).mockResolvedValue({
      success: true,
      data: {
        items: [
          {
            id: "abc123",
            title: "Need dark mode",
            description: "Please add dark mode support.",
            category: "Feature Request",
            status: "New",
            ai_sentiment: "Positive",
            ai_priority: 8,
            ai_summary: "User wants dark mode.",
            ai_tags: ["UI"],
            createdAt: new Date().toISOString(),
          },
        ],
        total: 1,
        page: 1,
        limit: 10,
        stats: {
          totalFeedback: 1,
          openItems: 1,
          averagePriority: 8,
          mostCommonTag: "UI",
        },
      },
    });

    render(<AdminPage />);

    await waitFor(() => {
      expect(getFeedbackList).toHaveBeenCalled();
    });
    expect(await screen.findByText(/need dark mode/i)).toBeInTheDocument();
    expect(screen.getByText(/most common tag/i)).toBeInTheDocument();
  });
});

