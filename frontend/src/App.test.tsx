import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import App from "./App";

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe("App", () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  it("renders navigation bar with links", () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ message: "API is running" }),
    });

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>
    );

    expect(screen.getByText("Home")).toBeInTheDocument();
    expect(screen.getByText("Repositories")).toBeInTheDocument();
    expect(screen.getByText("Documents")).toBeInTheDocument();
    expect(screen.getByText("Documentation")).toBeInTheDocument();
  });

  it("renders HomePage on root path", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ message: "API is running" }),
    });

    render(
      <MemoryRouter initialEntries={["/"]}>
        <App />
      </MemoryRouter>
    );

    expect(screen.getByText("OpsCore Documentation System")).toBeInTheDocument();
    expect(screen.getByText("Manage Repositories")).toBeInTheDocument();
    expect(screen.getByText("View Documentation")).toBeInTheDocument();
  });

  it("HomePage fetches and displays backend status", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ message: "Backend connected" }),
    });

    render(
      <MemoryRouter initialEntries={["/"]}>
        <App />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/Backend connected/)).toBeInTheDocument();
    });
  });

  it("HomePage handles API fetch error", async () => {
    mockFetch.mockRejectedValue(new Error("Network error"));

    render(
      <MemoryRouter initialEntries={["/"]}>
        <App />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/Failed to load message from backend/)).toBeInTheDocument();
    });
  });

  it("HomePage handles HTTP error", async () => {
    mockFetch.mockResolvedValue({
      ok: false,
      status: 500,
    });

    render(
      <MemoryRouter initialEntries={["/"]}>
        <App />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/Failed to load message from backend/)).toBeInTheDocument();
    });
  });
});
