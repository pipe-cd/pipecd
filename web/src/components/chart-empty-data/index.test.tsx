import { render, screen } from "~~/test-utils";
import ChartEmptyData from ".";

describe("ChartEmptyData", () => {
  it("renders with default text when visible", () => {
    render(<ChartEmptyData visible={true} />);
    expect(screen.getByText("No data is available.")).toBeInTheDocument();
  });

  it("renders with custom noDataText", () => {
    render(<ChartEmptyData visible={true} noDataText="Nothing here" />);
    expect(screen.getByText("Nothing here")).toBeInTheDocument();
  });

  it("is hidden when visible is false", () => {
    const { container } = render(<ChartEmptyData visible={false} />);
    const box = container.firstChild as HTMLElement;
    expect(box).toHaveStyle({ display: "none" });
  });

  it("is shown when visible is true", () => {
    const { container } = render(<ChartEmptyData visible={true} />);
    const box = container.firstChild as HTMLElement;
    expect(box).toHaveStyle({ display: "flex" });
  });
});
