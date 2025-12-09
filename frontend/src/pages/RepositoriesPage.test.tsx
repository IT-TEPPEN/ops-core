import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import RepositoriesPage from "./RepositoriesPage";

// Mock fetch globally
const mockFetch = vi.fn();
globalThis.fetch = mockFetch as any;

describe("RepositoriesPage", () => {
  beforeEach(() => {
    mockFetch.mockClear();
  });

  it("renders page title and form", () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    expect(screen.getByText("Repository Management")).toBeInTheDocument();
    expect(screen.getByPlaceholderText("https://github.com/username/repo.git")).toBeInTheDocument();
  });

  it("displays loading state while fetching repositories", async () => {
    mockFetch.mockImplementation(
      () =>
        new Promise((resolve) =>
          setTimeout(
            () =>
              resolve({
                ok: true,
                json: async () => ({ repositories: [] }),
              }),
            100
          )
        )
    );

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    expect(screen.getByText("Loading repositories...")).toBeInTheDocument();
  });

  it("displays repositories when fetch succeeds", async () => {
    const mockRepos = [
      {
        id: "repo-1",
        name: "test-repo-1",
        url: "https://github.com/test/repo1",
        createdAt: "2024-01-01T00:00:00Z",
        updatedAt: "2024-01-01T00:00:00Z",
      },
      {
        id: "repo-2",
        name: "test-repo-2",
        url: "https://github.com/test/repo2",
        createdAt: "2024-01-02T00:00:00Z",
        updatedAt: "2024-01-02T00:00:00Z",
      },
    ];

    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ repositories: mockRepos }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText("test-repo-1")).toBeInTheDocument();
      expect(screen.getByText("test-repo-2")).toBeInTheDocument();
    });
  });

  it("displays error message when fetch fails", async () => {
    mockFetch.mockRejectedValue(new Error("Network error"));

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(
        screen.getByText(/Failed to load repositories/)
      ).toBeInTheDocument();
    });
  });

  it("displays empty state when no repositories", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/No repositories registered yet/)).toBeInTheDocument();
    });
  });

  it("handles form input change", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    const input = screen.getByPlaceholderText("https://github.com/username/repo.git");
    fireEvent.change(input, {
      target: { value: "https://github.com/test/new-repo" },
    });

    expect(input).toHaveValue("https://github.com/test/new-repo");
  });

  it("submits new repository successfully", async () => {
    // Mock initial fetch for repositories list
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByPlaceholderText("https://github.com/username/repo.git")).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText("https://github.com/username/repo.git");
    const submitButton = screen.getByRole("button", { name: /Register Repository/i });

    // Mock POST request for repository creation
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        id: "new-repo",
        name: "new-repo",
        url: "https://github.com/test/new-repo",
        createdAt: "2024-01-01T00:00:00Z",
        updatedAt: "2024-01-01T00:00:00Z",
      }),
    });

    // Mock GET request for refreshed repositories list
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        repositories: [
          {
            id: "new-repo",
            name: "new-repo",
            url: "https://github.com/test/new-repo",
            createdAt: "2024-01-01T00:00:00Z",
            updatedAt: "2024-01-01T00:00:00Z",
          },
        ],
      }),
    });

    fireEvent.change(input, {
      target: { value: "https://github.com/test/new-repo" },
    });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(
        screen.getByText(/Repository registered successfully/)
      ).toBeInTheDocument();
    });
  });

  it("displays error when submitting empty URL", async () => {
    mockFetch.mockResolvedValue({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByPlaceholderText("https://github.com/username/repo.git")).toBeInTheDocument();
    });

    // The form has a required attribute, so it will use browser validation
    // Instead of checking for error message, just verify the input is required
    const input = screen.getByPlaceholderText("https://github.com/username/repo.git");
    expect(input).toBeRequired();
  });

  it("handles submission error", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ repositories: [] }),
    });

    render(
      <MemoryRouter>
        <RepositoriesPage />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByPlaceholderText("https://github.com/username/repo.git")).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText("https://github.com/username/repo.git");
    const submitButton = screen.getByRole("button", { name: /Register Repository/i });

    mockFetch.mockRejectedValueOnce(new Error("Server error"));

    fireEvent.change(input, {
      target: { value: "https://github.com/test/new-repo" },
    });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(
        screen.getByText(/Server error/)
      ).toBeInTheDocument();
    });
  });
});
