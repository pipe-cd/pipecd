import { render, screen } from "@testing-library/react";
import { DeleteUserGroupConfirmDialog } from ".";

describe("DeleteUserGroupConfirmDialog", () => {
  const mockOnDelete = jest.fn();
  const mockOnClose = jest.fn();

  beforeEach(() => {
    jest.clearAllMocks();
    render(
      <DeleteUserGroupConfirmDialog
        ssoGroup={"Developers"}
        onDelete={mockOnDelete}
        onCancel={mockOnClose}
      />
    );
  });

  it("renders the correct title", () => {
    expect(screen.getByText(/delete user group/i)).toBeInTheDocument();
  });

  it("renders the correct description", () => {
    expect(
      screen.getByText(/are you sure you want to delete the User Group/i)
    ).toBeInTheDocument();
    expect(screen.getByText(/Developers/)).toBeInTheDocument();
  });

  it("renders the delete button", () => {
    expect(screen.getByRole("button", { name: /delete/i })).toBeInTheDocument();
  });

  it("calls onDelete with the correct user group when delete button is clicked", () => {
    screen.getByRole("button", { name: /delete/i }).click();
    expect(mockOnDelete).toHaveBeenCalledWith("Developers");
  });

  it("calls onCancel when the dialog is closed", () => {
    screen.getByRole("button", { name: /cancel/i }).click();
    expect(mockOnClose).toHaveBeenCalled();
  });
});
