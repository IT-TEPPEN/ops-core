import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, waitFor, fireEvent } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import RepositoryDetailPage from "./RepositoryDetailPage";

// Mock fetch globally
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock("react-router-dom", async () => {
  const actual = await vi.importActual("react-router-dom");
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe("RepositoryDetailPage", () => {
  beforeEach(() => {
    mockFetch.mockClear();
    mockNavigate.mockClear();
  });

  const renderWithRouter = (repoId = "test-repo") => {
    return render(
      <MemoryRouter initialEntries={[`/repositories/${repoId}`]}>
        <Routes>
          <Route path="/repositories/:repoId" element={<RepositoryDetailPage />} />
        </Routes>
      </MemoryRouter>
    );
  };

  it("renders page title and loading state", async () => {
    mockFetch.mockImplementation(
      () =>
        new Promise((resolve) =>
          setTimeout(
            () =>
              resolve({
                ok: true,
                json: async () => ({}),
              }),
            100
          )
        )
    );

    renderWithRouter();

    expect(screen.getByText("Loading repository information...")).toBeInTheDocument();
  });

  it("displays repository details when fetch succeeds", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    // Mock repository fetch
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    // Mock files fetch
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ files: [] }),
    });

    renderWithRouter();

    await waitFor(() => {
      expect(screen.getByText("test-repository")).toBeInTheDocument();
    });
  });

  it("displays markdown files when available", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    const mockFiles = [
      { path: "README.md", type: "file" },
      { path: "docs/guide.md", type: "file" },
      { path: "docs", type: "dir" },
    ];

    // Mock repository fetch
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    // Mock files fetch
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ files: mockFiles }),
    });

    renderWithRouter();

    await waitFor(() => {
      expect(screen.getByText("README.md")).toBeInTheDocument();
      expect(screen.getByText("docs/guide.md")).toBeInTheDocument();
    });
  });

  it("displays error message when repository fetch fails", async () => {
    mockFetch.mockRejectedValue(new Error("Network error"));

    renderWithRouter();

    await waitFor(() => {
      expect(
        screen.getByText(/Failed to load repository details/)
      ).toBeInTheDocument();
    });
  });

  it("handles file selection", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    const mockFiles = [{ path: "README.md", type: "file" }];

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ files: mockFiles }),
    });

    renderWithRouter();

    await waitFor(() => {
      expect(screen.getByText("README.md")).toBeInTheDocument();
    });

    const checkbox = screen.getByRole("checkbox");
    fireEvent.click(checkbox);

    expect(checkbox).toBeChecked();
  });

  it("displays access token input when needed", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    // Mock repository fetch
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    // Mock files fetch with 400 error and ACCESS_TOKEN_REQUIRED code
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 400,
      json: async () => ({ code: "ACCESS_TOKEN_REQUIRED" }),
    });

    renderWithRouter();

    await waitFor(() => {
      expect(
        screen.getByText(/Access token is required/)
      ).toBeInTheDocument();
    });
  });

  it("renders back to repositories link", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ files: [] }),
    });

    renderWithRouter();

    await waitFor(() => {
      const backLink = screen.getByText(/Back to Repositories/);
      expect(backLink).toBeInTheDocument();
    });
  });

  it("handles file selection submission", async () => {
    const mockRepo = {
      id: "test-repo",
      name: "test-repository",
      url: "https://github.com/test/repo",
      createdAt: "2024-01-01T00:00:00Z",
      updatedAt: "2024-01-01T00:00:00Z",
    };

    const mockFiles = [{ path: "README.md", type: "file" }];

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockRepo,
    });

    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ files: mockFiles }),
    });

    renderWithRouter();

    await waitFor(() => {
      expect(screen.getByText("README.md")).toBeInTheDocument();
    });

    const checkbox = screen.getByRole("checkbox");
    fireEvent.click(checkbox);

    // Mock file selection submission
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({}),
    });

    const submitButton = screen.getByRole("button", { name: /Process Selected Files/i });
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(screen.getByText(/Files selected successfully/)).toBeInTheDocument();
    });
  });
});
