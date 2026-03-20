import { fireEvent, render, screen } from "~~/test-utils";
import { WarningBanner } from ".";

describe("WarningBanner", () => {
  it("renders the availability message", () => {
    render(<WarningBanner onClose={jest.fn()} />);
    expect(screen.getByText(/is available!/i)).toBeInTheDocument();
  });

  it("renders a link to the release notes", () => {
    const version = "v0.99.0";
    process.env.PIPECD_VERSION = version;
    render(<WarningBanner onClose={jest.fn()} />);
    const link = screen.getByRole("link");
    expect(link).toHaveAttribute(
      "href",
      `https://github.com/pipe-cd/pipecd/releases/tag/${version}`
    );
  });

  it("calls onClose when the close icon is clicked", () => {
    const onClose = jest.fn();
    const { container } = render(<WarningBanner onClose={onClose} />);
    // CloseIcon renders as an svg inside the AppBar
    const closeIcon = container.querySelector("svg");
    expect(closeIcon).not.toBeNull();
    fireEvent.click(closeIcon!);
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
