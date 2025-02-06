import DialogConfirm from ".";
import { render, screen } from "~~/test-utils";
import userEvent from "@testing-library/user-event";

const onCancelMock = jest.fn();
const onConfirmMock = jest.fn();

const props = {
  open: true,
  onCancel: onCancelMock,
  onConfirm: onConfirmMock,
  title: "title",
  description: "description",
  cancelText: "cancel",
  confirmText: "confirm",
};

describe("DialogConfirm component", () => {
  it("should have correct title and description", () => {
    render(<DialogConfirm {...props} />, {});
    expect(screen.getByText("title")).toBeInTheDocument();
    expect(screen.getByText("description")).toBeInTheDocument();
  });

  it("calls onCancel when cancel button is clicked", () => {
    render(<DialogConfirm {...props} />, {});
    const button = screen.getByRole("button", { name: "cancel" });
    userEvent.click(button);
    expect(button).toBeInTheDocument();
    expect(onCancelMock).toHaveBeenCalledTimes(1);
  });

  it("calls onConfirm when confirm button is clicked", () => {
    render(<DialogConfirm {...props} />, {});
    const button = screen.getByRole("button", { name: "confirm" });
    userEvent.click(button);
    expect(button).toBeInTheDocument();
    expect(onConfirmMock).toHaveBeenCalledTimes(1);
  });
});
