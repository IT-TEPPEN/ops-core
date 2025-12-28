import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatisticsChart } from "./StatisticsChart";

describe("StatisticsChart", () => {
  const mockData = [
    { label: "Week 1", value: 100 },
    { label: "Week 2", value: 150 },
    { label: "Week 3", value: 200 },
  ];

  it("displays chart with title and data", () => {
    render(<StatisticsChart title="Test Chart" data={mockData} />);
    expect(screen.getByText("Test Chart")).toBeInTheDocument();
    expect(screen.getByText("Week 1")).toBeInTheDocument();
    expect(screen.getByText("100")).toBeInTheDocument();
    expect(screen.getByText("Week 2")).toBeInTheDocument();
    expect(screen.getByText("150")).toBeInTheDocument();
    expect(screen.getByText("Week 3")).toBeInTheDocument();
    expect(screen.getByText("200")).toBeInTheDocument();
  });

  it("shows loading spinner when isLoading is true", () => {
    const { container } = render(
      <StatisticsChart title="Test Chart" data={[]} isLoading={true} />
    );
    expect(container.querySelector(".animate-spin")).toBeInTheDocument();
  });

  it("shows empty state when no data", () => {
    render(<StatisticsChart title="Test Chart" data={[]} />);
    expect(screen.getByText("No data available.")).toBeInTheDocument();
  });

  it("renders correct number of bars", () => {
    const { container } = render(
      <StatisticsChart title="Test Chart" data={mockData} />
    );
    const bars = container.querySelectorAll(".bg-blue-600");
    expect(bars).toHaveLength(3);
  });
});
