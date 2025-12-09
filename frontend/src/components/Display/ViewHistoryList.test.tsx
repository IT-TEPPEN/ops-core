import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ViewHistoryList } from "./ViewHistoryList";
import type { ViewHistory } from "../../types/domain";

describe("ViewHistoryList", () => {
  const mockItems: ViewHistory[] = [
    {
      id: "vh-1",
      document_id: "doc-123",
      user_id: "user-456",
      viewed_at: "2024-01-15T10:00:00Z",
      view_duration: 120,
    },
    {
      id: "vh-2",
      document_id: "doc-789",
      user_id: "user-456",
      viewed_at: "2024-01-15T11:00:00Z",
      view_duration: 60,
    },
  ];

  it("displays view history items", () => {
    render(<ViewHistoryList items={mockItems} />);
    expect(screen.getByText("Document ID")).toBeInTheDocument();
    expect(screen.getByText("User ID")).toBeInTheDocument();
    expect(screen.getAllByText("user-456")).toHaveLength(2);
  });

  it("shows loading spinner when isLoading is true", () => {
    const { container } = render(<ViewHistoryList items={[]} isLoading={true} />);
    expect(container.querySelector(".animate-spin")).toBeInTheDocument();
  });

  it("shows empty state when no items", () => {
    render(<ViewHistoryList items={[]} />);
    expect(screen.getByText("No view history records found.")).toBeInTheDocument();
  });

  it("displays correct number of rows", () => {
    const { container } = render(<ViewHistoryList items={mockItems} />);
    const rows = container.querySelectorAll("tbody tr");
    expect(rows).toHaveLength(2);
  });
});
