import { render, screen, waitFor } from "~~/test-utils";
import { AddUserGroupDialog, AddUserGroupDialogProps } from "./index";
import userEvent from "@testing-library/user-event";

// Mock useGetProject hook
jest.mock("~/queries/project/use-get-project", () => ({
  useGetProject: () => ({
    data: {
      rbacRoles: [{ name: "Admin" }, { name: "Viewer" }],
    },
  }),
}));

describe("AddUserGroupDialog", () => {
  const defaultProps: AddUserGroupDialogProps = {
    open: true,
    onClose: jest.fn(),
    onSubmit: jest.fn(),
  };

  it("renders dialog title and fields", async () => {
    render(<AddUserGroupDialog {...defaultProps} />);
    await screen.findByRole("dialog");
    expect(screen.getByText("Add User Group")).toBeInTheDocument();
    expect(
      screen.getByRole("textbox", { name: /team\/group/i })
    ).toBeInTheDocument();
    expect(screen.getByRole("combobox", { name: /role/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /add/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /close/i })).toBeInTheDocument();
  });

  it("calls onClose when Close button is clicked", async () => {
    render(<AddUserGroupDialog {...defaultProps} />);
    await screen.findByRole("dialog");
    userEvent.click(screen.getByRole("button", { name: /close/i }));
    expect(defaultProps.onClose).toHaveBeenCalled();
  });

  it("calls onSubmit with correct values", async () => {
    render(<AddUserGroupDialog {...defaultProps} />);
    await screen.findByRole("dialog");

    const textInput = screen.getByRole("textbox", { name: /team\/group/i });
    userEvent.type(textInput, "dev-team");

    userEvent.click(screen.getByRole("combobox"));
    const adminOption = await screen.findByRole("option", { name: "Admin" });
    userEvent.click(adminOption);

    await waitFor(() => {
      expect(screen.getByText("Admin")).toBeInTheDocument();
    });

    userEvent.click(screen.getByRole("button", { name: /add/i }));

    await waitFor(() => {
      expect(defaultProps.onSubmit).toHaveBeenCalledWith({
        ssoGroup: "dev-team",
        role: "Admin",
      });
    });
  });

  it("disables Add button if form is invalid or pristine", async () => {
    render(<AddUserGroupDialog {...defaultProps} />);
    const addButton = screen.getByRole("button", { name: /add/i });
    expect(addButton).toBeDisabled();

    const textInput = screen.getByRole("textbox", { name: /team\/group/i });
    userEvent.type(textInput, "dev-team");

    userEvent.click(screen.getByRole("combobox"));
    const adminOption = await screen.findByRole("option", { name: "Admin" });
    userEvent.click(adminOption);

    await waitFor(() => {
      expect(screen.getByText("Admin")).toBeInTheDocument();
    });

    // Now the form is valid and dirty
    expect(addButton).not.toBeDisabled();
  });
});
