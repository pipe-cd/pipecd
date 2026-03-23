import { waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { render, screen } from "~~/test-utils";
import { TextWithCopyButton } from ".";

describe("TextWithCopyButton", () => {
  it("renders the value in a readonly input", () => {
    render(<TextWithCopyButton name="API Key" value="my-secret-key" />);
    const input = screen.getByDisplayValue("my-secret-key");
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute("readonly");
  });

  it("renders the name as a legend", () => {
    render(<TextWithCopyButton name="API Key" value="my-secret-key" />);
    expect(screen.getByText("API Key")).toBeInTheDocument();
  });

  it("copies value to clipboard when copy button is clicked", async () => {
    jest.spyOn(navigator.clipboard, "writeText");
    render(<TextWithCopyButton name="Token" value="abc-123" />);
    userEvent.click(screen.getByRole("button", { name: /copy token/i }));
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith("abc-123");
    await waitFor(() =>
      expect(screen.getByText("Token copied to clipboard.")).toBeInTheDocument()
    );
  });
});
