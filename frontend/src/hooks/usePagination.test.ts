import { describe, it, expect } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { usePagination } from "./usePagination";

describe("usePagination", () => {
  it("initializes with correct default values", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10 })
    );

    expect(result.current.currentPage).toBe(1);
    expect(result.current.pageSize).toBe(10);
    expect(result.current.totalItems).toBe(100);
    expect(result.current.totalPages).toBe(10);
    expect(result.current.hasPrevPage).toBe(false);
    expect(result.current.hasNextPage).toBe(true);
  });

  it("navigates to next page", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10 })
    );

    act(() => {
      result.current.nextPage();
    });

    expect(result.current.currentPage).toBe(2);
    expect(result.current.hasPrevPage).toBe(true);
  });

  it("navigates to previous page", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10, initialPage: 3 })
    );

    act(() => {
      result.current.prevPage();
    });

    expect(result.current.currentPage).toBe(2);
  });

  it("navigates to specific page", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10 })
    );

    act(() => {
      result.current.goToPage(5);
    });

    expect(result.current.currentPage).toBe(5);
  });

  it("clamps page to valid range", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10 })
    );

    act(() => {
      result.current.goToPage(100);
    });

    expect(result.current.currentPage).toBe(10);

    act(() => {
      result.current.goToPage(0);
    });

    expect(result.current.currentPage).toBe(1);
  });

  it("navigates to first and last page", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10, initialPage: 5 })
    );

    act(() => {
      result.current.firstPage();
    });

    expect(result.current.currentPage).toBe(1);

    act(() => {
      result.current.lastPage();
    });

    expect(result.current.currentPage).toBe(10);
  });

  it("changes page size and resets to first page", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10, initialPage: 5 })
    );

    act(() => {
      result.current.setPageSize(20);
    });

    expect(result.current.pageSize).toBe(20);
    expect(result.current.currentPage).toBe(1);
    expect(result.current.totalPages).toBe(5);
  });

  it("slices items correctly", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 30, pageSize: 10, initialPage: 2 })
    );

    const items = Array.from({ length: 30 }, (_, i) => i);
    const pageItems = result.current.getPageItems(items);

    expect(pageItems).toEqual([10, 11, 12, 13, 14, 15, 16, 17, 18, 19]);
  });

  it("generates page numbers with ellipsis", () => {
    const { result } = renderHook(() =>
      usePagination({ totalItems: 200, pageSize: 10, initialPage: 10 })
    );

    const pageNumbers = result.current.getPageNumbers(7);
    expect(pageNumbers).toContain(1);
    expect(pageNumbers).toContain("...");
    expect(pageNumbers).toContain(20);
  });

  it("calls onPageChange callback", () => {
    const onPageChange = vi.fn();
    const { result } = renderHook(() =>
      usePagination({ totalItems: 100, pageSize: 10, onPageChange })
    );

    act(() => {
      result.current.nextPage();
    });

    expect(onPageChange).toHaveBeenCalledWith(2);
  });
});

import { vi } from "vitest";
