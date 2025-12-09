import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import BlogPage from "./BlogPage";

// Mock fetch globally
const mockFetch = vi.fn();
globalThis.fetch = mockFetch as unknown as typeof fetch;

// Mock useSearchParams
const mockSearchParams = new URLSearchParams();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useSearchParams: () => [mockSearchParams],
  };
});

describe("BlogPage", () => {
  beforeEach(() => {
    mockFetch.mockClear();
    mockSearchParams.delete("repoId");
  });

  it("renders without repository ID", () => {
    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    expect(screen.getByText("Documentation")).toBeInTheDocument();
    expect(screen.getByText(/No repository ID provided/)).toBeInTheDocument();
  });

  it("displays loading state while fetching", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockImplementation(
      () =>
        new Promise((resolve) =>
          setTimeout(
            () =>
              resolve({
                ok: true,
                json: async () => ({ content: "# Test Content" }),
              }),
            100
          )
        )
    );

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    expect(screen.getByText("Loading content...")).toBeInTheDocument();
  });

  it("displays markdown content when fetch succeeds", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({
        content: "# Test Heading\n\nTest paragraph",
      }),
    });

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText("Test Heading")).toBeInTheDocument();
      expect(screen.getByText("Test paragraph")).toBeInTheDocument();
    });
  });

  it("handles markdown with frontmatter", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({
        content: `---
title: "Test Document"
author: "Test Author"
---

# Content Heading

Content body`,
      }),
    });

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText("Document Metadata")).toBeInTheDocument();
      expect(screen.getByText("title")).toBeInTheDocument();
      expect(screen.getByText("Test Document")).toBeInTheDocument();
      expect(screen.getByText("author")).toBeInTheDocument();
      expect(screen.getByText("Test Author")).toBeInTheDocument();
      expect(screen.getByText("Content Heading")).toBeInTheDocument();
      expect(screen.getByText("Content body")).toBeInTheDocument();
    });
  });

  it("displays error message when fetch fails", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockRejectedValue(new Error("Network error"));

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(
        screen.getByText(/Failed to load markdown content/)
      ).toBeInTheDocument();
    });
  });

  it("displays error message on HTTP error", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockResolvedValue({
      ok: false,
      status: 404,
    });

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(
        screen.getByText(/Failed to load markdown content/)
      ).toBeInTheDocument();
    });
  });

  it("renders back to repositories link", () => {
    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    const backLink = screen.getByText("Back to Repositories");
    expect(backLink).toBeInTheDocument();
    expect(backLink.closest("a")).toHaveAttribute("href", "/repositories");
  });

  it("displays message when no content available", async () => {
    mockSearchParams.set("repoId", "test-repo");
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ content: "" }),
    });

    render(
      <MemoryRouter>
        <BlogPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(
        screen.getByText(/No markdown content available/)
      ).toBeInTheDocument();
    });
  });
});
