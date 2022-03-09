import { waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { render, screen, createStore } from "~~/test-utils";
import { addToast } from "~/modules/toasts";
import { CopyIconButton } from "./";

test("copy text", async () => {
  const store = createStore();
  jest.spyOn(navigator.clipboard, "writeText");
  render(<CopyIconButton name="ID" value="id" />, { store });
  userEvent.click(screen.getByRole("button"));

  expect(navigator.clipboard.writeText).toHaveBeenCalledWith("id");
  await waitFor(() =>
    expect(store.getActions()).toEqual([
      {
        payload: {
          message: "ID copied to clipboard.",
        },
        type: addToast.type,
      },
    ])
  );
});
