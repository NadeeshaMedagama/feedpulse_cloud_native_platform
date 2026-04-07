import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import HomePage from "@/app/page";
import { submitFeedback } from "@/lib/api";

jest.mock("@/lib/api", () => ({
  submitFeedback: jest.fn(),
}));

describe("HomePage", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows validation error when title is empty", async () => {
    render(<HomePage />);

    fireEvent.change(screen.getByLabelText(/^title$/i), { target: { value: "   " } });
    fireEvent.change(screen.getByLabelText(/description/i), {
      target: { value: "This description is long enough for validation." },
    });
    fireEvent.click(screen.getByRole("button", { name: /submit feedback/i }));

    expect(await screen.findByText(/title is required and description must be at least 20 characters/i)).toBeInTheDocument();
    expect(submitFeedback).not.toHaveBeenCalled();
  });

  it("submits valid feedback successfully", async () => {
    (submitFeedback as jest.Mock).mockResolvedValue({ success: true, data: {} });

    render(<HomePage />);

    fireEvent.change(screen.getByLabelText(/^title$/i), { target: { value: "Need dark mode" } });
    fireEvent.change(screen.getByLabelText(/description/i), {
      target: { value: "Please add dark mode support to the admin dashboard soon." },
    });

    fireEvent.click(screen.getByRole("button", { name: /submit feedback/i }));

    await waitFor(() => {
      expect(submitFeedback).toHaveBeenCalledTimes(1);
    });
    expect(await screen.findByText(/feedback submitted successfully/i)).toBeInTheDocument();
  });
});


