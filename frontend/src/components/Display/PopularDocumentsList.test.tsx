import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { PopularDocumentsList } from "./PopularDocumentsList";
import type { PopularDocument } from "../../types/domain";

describe("PopularDocumentsList", () => {
  const mockItems: PopularDocument[] = [
    {
      document_id: "doc-1",
      total_views: 1234,
      unique_viewers: 567,
      last_viewed_at: "2024-01-15T10:00:00Z",
    },
    {
      document_id: "doc-2",
      total_views: 987,
      unique_viewers: 432,
      last_viewed_at: "2024-01-15T09:00:00Z",
    },
  ];

  it("displays popular documents with rankings", () => {
    render(<PopularDocumentsList items={mockItems} />);
    expect(screen.getByText("Popular Documents")).toBeInTheDocument();
    expect(screen.getByText("1")).toBeInTheDocument();
    expect(screen.getByText("2")).toBeInTheDocument();
  });

  it("shows loading spinner when isLoading is true", () => {
    const { container } = render(<PopularDocumentsList items={[]} isLoading={true} />);
    expect(container.querySelector(".animate-spin")).toBeInTheDocument();
  });

  it("shows empty state when no items", () => {
    render(<PopularDocumentsList items={[]} />);
    expect(screen.getByText("No popular documents found.")).toBeInTheDocument();
  });

  it("displays view counts and viewer counts", () => {
    render(<PopularDocumentsList items={mockItems} />);
    expect(screen.getByText("1,234")).toBeInTheDocument();
    expect(screen.getByText("567")).toBeInTheDocument();
    expect(screen.getAllByText("views")).toHaveLength(2);
    expect(screen.getAllByText("viewers")).toHaveLength(2);
  });
});
