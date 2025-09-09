import { render, screen, fireEvent } from "@testing-library/react";
import LabelAutoComplete from "./index";
import userEvent from "@testing-library/user-event";

describe("LabelAutoComplete", () => {
  const labels = ["env:bug", "env:production", "label:urgent", "abc:acb"];

  it("renders input and suggestions", () => {
    render(<LabelAutoComplete options={labels} />);
    expect(screen.getByRole("combobox")).toBeInTheDocument();
  });

  it("show value ", () => {
    render(<LabelAutoComplete options={labels} value={["env:bug"]} />);
    expect(screen.getByText("env:bug")).toBeInTheDocument();
  });

  it("shows suggestions when typing", () => {
    render(<LabelAutoComplete options={labels} />);
    const input = screen.getByRole("combobox");
    fireEvent.change(input, { target: { value: "b" } });
    expect(screen.getByText("env:bug")).toBeInTheDocument();
  });

  it("calls onChange when a suggestion is clicked", () => {
    const handleChange = jest.fn();
    render(<LabelAutoComplete options={labels} onChange={handleChange} />);
    const input = screen.getByRole("combobox");
    fireEvent.change(input, { target: { value: "produ" } });
    fireEvent.click(screen.getByText("env:production"));
    expect(handleChange).toHaveBeenCalledWith(["env:production"]);
  });

  it("filters options based on input", () => {
    render(<LabelAutoComplete options={labels} />);
    const input = screen.getByRole("combobox");
    fireEvent.change(input, { target: { value: "u" } });
    expect(screen.getByText("label:urgent")).toBeInTheDocument();
    expect(screen.queryByText("abc:acb")).not.toBeInTheDocument();
  });

  it("does not show options if input does not match", () => {
    render(<LabelAutoComplete options={labels} />);
    const input = screen.getByRole("combobox");
    fireEvent.change(input, { target: { value: "xyz" } });
    expect(screen.queryByText("bug")).not.toBeInTheDocument();
    expect(screen.queryByText("production")).not.toBeInTheDocument();
    expect(screen.queryByText("urgent")).not.toBeInTheDocument();
  });

  it("accept new label that match 'key:value' event if inputValue not in options", () => {
    const handleChange = jest.fn();
    render(<LabelAutoComplete options={labels} onChange={handleChange} />);
    const input = screen.getByRole("combobox");
    userEvent.type(input, "new:label{enter}");
    expect(handleChange).toHaveBeenCalledWith(["new:label"]);
  });
  it("not accept new label that does not match 'key:value'", () => {
    const handleChange = jest.fn();
    render(<LabelAutoComplete options={labels} onChange={handleChange} />);
    const input = screen.getByRole("combobox");
    userEvent.type(input, "new-label{enter}");
    expect(handleChange).toHaveBeenCalledWith([]);
  });
});
