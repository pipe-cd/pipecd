import { screen, render, waitFor } from "~~/test-utils";
import BreakingChangeNotes from "./index"; // adjust import if needed
import userEvent from "@testing-library/user-event";

describe("BreakingChange component", () => {
  it("renders empty if no breaking changes", async () => {
    render(<BreakingChangeNotes notes={""} />);

    expect(
      screen.queryByText(/breaking change notes/i)
    ).not.toBeInTheDocument();
  });

  it("renders breaking changes note after fetching", async () => {
    render(<BreakingChangeNotes notes={"warning notes"} />);

    expect(screen.getByText(/warning notes/)).toBeInTheDocument();
    expect(
      screen.getByRole("button", { name: /view detail/i })
    ).toBeInTheDocument();

    userEvent.click(screen.getByRole("button", { name: /view detail/i }));
    expect(screen.getByText(/Breaking Changes/i)).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /close/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /ignore/i })).toBeInTheDocument();
  });

  it("close dialog when close button is clicked", async () => {
    render(<BreakingChangeNotes notes={"warning notes"} />);

    await userEvent.click(screen.getByRole("button", { name: /view detail/i }));
    await userEvent.click(screen.getByRole("button", { name: /close/i }));

    // Wait for the dialog to be removed from the DOM
    await waitFor(() => {
      expect(screen.queryByText(/Breaking Changes/i)).not.toBeInTheDocument();
    });
  });

  it("close dialog when ignore button is clicked", async () => {
    render(<BreakingChangeNotes notes={"warning notes"} />);

    await userEvent.click(screen.getByRole("button", { name: /view detail/i }));
    await userEvent.click(screen.getByRole("button", { name: /ignore/i }));

    // Wait for the dialog to be removed from the DOM
    await waitFor(() => {
      expect(screen.queryByText(/Breaking Changes/i)).not.toBeInTheDocument();
    });
  });
});
