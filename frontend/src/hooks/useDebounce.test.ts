import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useDebounce, useDebouncedCallback, useDebouncedState } from "./useDebounce";

describe("useDebounce", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("returns initial value immediately", () => {
    const { result } = renderHook(() => useDebounce("initial", 500));
    expect(result.current).toBe("initial");
  });

  it("debounces value updates", () => {
    const { result, rerender } = renderHook(
      ({ value }) => useDebounce(value, 500),
      { initialProps: { value: "initial" } }
    );

    expect(result.current).toBe("initial");

    rerender({ value: "updated" });
    expect(result.current).toBe("initial");

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current).toBe("updated");
  });
});

describe("useDebouncedCallback", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("debounces callback execution", () => {
    const callback = vi.fn();
    const { result } = renderHook(() => useDebouncedCallback(callback, 500));

    act(() => {
      result.current("arg1");
      result.current("arg2");
      result.current("arg3");
    });

    expect(callback).not.toHaveBeenCalled();

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(callback).toHaveBeenCalledTimes(1);
    expect(callback).toHaveBeenCalledWith("arg3");
  });
});

describe("useDebouncedState", () => {
  beforeEach(() => {
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("returns current value and debounced value", () => {
    const { result } = renderHook(() => useDebouncedState("initial", 500));

    const [value, debouncedValue] = result.current;
    expect(value).toBe("initial");
    expect(debouncedValue).toBe("initial");
  });

  it("updates current value immediately but debounces debounced value", () => {
    const { result } = renderHook(() => useDebouncedState("initial", 500));

    act(() => {
      result.current[2]("updated");
    });

    expect(result.current[0]).toBe("updated");
    expect(result.current[1]).toBe("initial");

    act(() => {
      vi.advanceTimersByTime(500);
    });

    expect(result.current[1]).toBe("updated");
  });
});
