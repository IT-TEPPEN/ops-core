import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { DocumentStatistics } from "./DocumentStatistics";
import type { DocumentStatistics as DocumentStatisticsType } from "../../types/domain";

describe("DocumentStatistics", () => {
  const mockStatistics: DocumentStatisticsType = {
    document_id: "doc-123",
    total_views: 1234,
    unique_viewers: 567,
    last_viewed_at: "2024-01-15T10:00:00Z",
    average_view_duration: 120,
  };

  it("displays all statistics", () => {
    render(<DocumentStatistics statistics={mockStatistics} />);
    expect(screen.getByText("Total Views")).toBeInTheDocument();
    expect(screen.getByText("1,234")).toBeInTheDocument();
    expect(screen.getByText("Unique Viewers")).toBeInTheDocument();
    expect(screen.getByText("567")).toBeInTheDocument();
    expect(screen.getByText("Avg. Duration")).toBeInTheDocument();
    expect(screen.getByText("120s")).toBeInTheDocument();
  });

  it("shows loading spinner when isLoading is true", () => {
    const { container } = render(
      <DocumentStatistics statistics={mockStatistics} isLoading={true} />
    );
    expect(container.querySelector(".animate-spin")).toBeInTheDocument();
  });

  it("displays correct icons", () => {
    const { container } = render(<DocumentStatistics statistics={mockStatistics} />);
    expect(container.textContent).toContain("📊");
    expect(container.textContent).toContain("👥");
    expect(container.textContent).toContain("🕒");
    expect(container.textContent).toContain("⏱️");
  });
});
