import { waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { render, screen } from "~~/test-utils";
import { CopyIconButton } from "./";

test("copy text", async () => {
  jest.spyOn(navigator.clipboard, "writeText");
  render(<CopyIconButton name="ID" value="id" />);
  userEvent.click(screen.getByRole("button"));

  expect(navigator.clipboard.writeText).toHaveBeenCalledWith("id");
  await waitFor(() =>
    expect(screen.getByText("ID copied to clipboard.")).toBeInTheDocument()
  );
});
