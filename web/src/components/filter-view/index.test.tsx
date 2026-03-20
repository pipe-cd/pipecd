import userEvent from "@testing-library/user-event";
import { render, screen } from "~~/test-utils";
import { FilterView } from ".";

describe("FilterView", () => {
  it("renders the Filters heading", () => {
    render(<FilterView onClear={jest.fn()}>child</FilterView>);
    expect(screen.getByText("Filters")).toBeInTheDocument();
  });

  it("renders children", () => {
    render(
      <FilterView onClear={jest.fn()}>
        <span>my-filter-child</span>
      </FilterView>
    );
    expect(screen.getByText("my-filter-child")).toBeInTheDocument();
  });

  it("calls onClear when CLEAR button is clicked", () => {
    const onClear = jest.fn();
    render(<FilterView onClear={onClear}>child</FilterView>);
    userEvent.click(screen.getByRole("button", { name: /clear/i }));
    expect(onClear).toHaveBeenCalledTimes(1);
  });
});
