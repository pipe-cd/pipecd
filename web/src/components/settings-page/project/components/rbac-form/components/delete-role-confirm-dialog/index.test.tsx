import { render, screen } from "@testing-library/react";
import { DeleteRoleConfirmDialog } from ".";

describe("DeleteRoleConfirmDialog", () => {
  const mockOnDelete = jest.fn();
  const mockOnClose = jest.fn();

  beforeEach(() => {
    render(
      <DeleteRoleConfirmDialog
        roleName={"Developers"}
        onDelete={mockOnDelete}
        onClose={mockOnClose}
      />
    );
  });

  afterAll(() => {
    jest.clearAllMocks();
  });

  it("renders the correct title", () => {
    expect(screen.getByText(/delete role/i)).toBeInTheDocument();
  });

  it("renders the correct description", () => {
    expect(
      screen.getByText(/are you sure you want to delete the role/i)
    ).toBeInTheDocument();
    expect(screen.getByText(/Developers/)).toBeInTheDocument();
  });

  it("renders the delete button", () => {
    expect(screen.getByRole("button", { name: /delete/i })).toBeInTheDocument();
  });

  it("calls onDelete with the correct role when delete button is clicked", () => {
    screen.getByRole("button", { name: /delete/i }).click();
    expect(mockOnDelete).toHaveBeenCalledWith("Developers");
  });

  it("calls onClose when the dialog is closed", () => {
    screen.getByRole("button", { name: /cancel/i }).click();
    expect(mockOnClose).toHaveBeenCalled();
  });
});
