import { render, screen, fireEvent } from "@testing-library/react";
import DeploymentTraceFilter from "./index";
import { MemoryRouter } from "~~/test-utils";

jest.useFakeTimers();

describe("DeploymentTraceFilter", () => {
  const mockOnChange = jest.fn();
  const mockOnClear = jest.fn();
  const filterValues = { commitHash: "12345" };

  beforeEach(() => {
    render(
      <MemoryRouter>
        <DeploymentTraceFilter
          filterValues={filterValues}
          onChange={mockOnChange}
          onClear={mockOnClear}
        />
      </MemoryRouter>
    );
  });

  it("should render filter inputs", () => {
    expect(
      screen.getByRole("textbox", { name: /commit hash/i })
    ).toBeInTheDocument();
  });

  it("should call onChange when filter value changes", () => {
    const input = screen.getByRole("textbox", { name: /commit hash/i });
    fireEvent.change(input, { target: { value: "67890" } });

    jest.runAllTimers();
    expect(mockOnChange).toHaveBeenCalledWith({ commitHash: "67890" });
  });

  it("should call onClear when clear button is clicked", () => {
    const clearButton = screen.getByRole("button", { name: /clear/i });
    fireEvent.click(clearButton);
    expect(mockOnClear).toHaveBeenCalled();
  });
});
