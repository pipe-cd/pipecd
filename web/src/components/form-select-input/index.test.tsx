import { render, screen, fireEvent } from "~~/test-utils";
import FormSelectInput from ".";

const options = [
  { value: "opt1", label: "Option 1" },
  { value: "opt2", label: "Option 2" },
  { value: "opt3", label: "Option 3" },
];

describe("FormSelectInput", () => {
  it("renders the label", () => {
    render(
      <FormSelectInput id="test-select" label="My Label" options={options} />
    );
    expect(screen.getByText("My Label")).toBeInTheDocument();
  });

  it("renders all options when opened", () => {
    render(
      <FormSelectInput id="test-select" label="Kind" options={options} />
    );
    fireEvent.mouseDown(screen.getByRole("combobox"));
    expect(screen.getByText("Option 1")).toBeInTheDocument();
    expect(screen.getByText("Option 2")).toBeInTheDocument();
    expect(screen.getByText("Option 3")).toBeInTheDocument();
  });

  it("calls onChange with correct value and option when an item is selected", () => {
    const onChange = jest.fn();
    render(
      <FormSelectInput
        id="test-select"
        label="Kind"
        options={options}
        onChange={onChange}
      />
    );
    fireEvent.mouseDown(screen.getByRole("combobox"));
    fireEvent.click(screen.getByText("Option 2"));
    expect(onChange).toHaveBeenCalledWith("opt2", options[1]);
  });

  it("uses defaultValue as initial selection", () => {
    render(
      <FormSelectInput
        id="test-select"
        options={options}
        defaultValue="opt3"
      />
    );
    expect(screen.getByRole("combobox")).toHaveTextContent("Option 3");
  });

  it("is disabled when disabled prop is true", () => {
    render(
      <FormSelectInput id="test-select" options={options} disabled={true} />
    );
    expect(screen.getByRole("combobox")).toHaveAttribute("aria-disabled", "true");
  });
});
