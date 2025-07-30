import { render, screen, waitFor } from "~~/test-utils";
import userEvent from "@testing-library/user-event";
import { AddRoleDialog } from ".";

describe("AddRoleDialog", () => {
  const defaultProps = {
    open: true,
    onClose: jest.fn(),
    onSubmit: jest.fn(),
  };

  it("renders dialog title and fields", async () => {
    render(<AddRoleDialog {...defaultProps} />);
    await waitFor(() => {
      screen.findByRole("dialog");
    });
    expect(screen.getByText("Add Role")).toBeInTheDocument();
    expect(screen.getByRole("textbox", { name: /role/i })).toBeInTheDocument();
    expect(
      screen.getByRole("textbox", { name: /policies/i })
    ).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /add/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /close/i })).toBeInTheDocument();
  });

  it("calls onClose when Close button is clicked", async () => {
    render(<AddRoleDialog {...defaultProps} />);
    await waitFor(() => {
      screen.findByRole("dialog");
    });
    userEvent.click(screen.getByRole("button", { name: /close/i }));
    expect(defaultProps.onClose).toHaveBeenCalled();
  });

  it("calls onSubmit with correct values", async () => {
    render(<AddRoleDialog {...defaultProps} />);

    const textInput = screen.getByRole("textbox", { name: /role/i });
    const policiesInput = screen.getByRole("textbox", { name: /policies/i });

    userEvent.type(textInput, "dev-team");
    userEvent.type(policiesInput, "resources=application;actions=get");

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /add/i })).toBeEnabled();
    });

    userEvent.click(screen.getByRole("button", { name: /add/i }));

    await waitFor(() => {
      expect(defaultProps.onSubmit).toHaveBeenCalledWith({
        name: "dev-team",
        policies: "resources=application;actions=get",
      });
    });
  });
});
