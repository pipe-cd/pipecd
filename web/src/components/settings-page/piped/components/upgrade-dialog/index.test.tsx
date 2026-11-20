import userEvent from "@testing-library/user-event";
import { render, screen, waitFor, within } from "~~/test-utils";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import * as pipedsApi from "~/api/piped";
import { UpgradePipedDialog } from "./index";

afterEach(() => jest.restoreAllMocks());

it("requires a fresh Piped selection after canceling an upgrade", async () => {
  const upgrade = jest
    .spyOn(pipedsApi, "updatePipedDesiredVersion")
    .mockResolvedValue({});

  const enabledPiped = { ...dummyPiped, id: "piped-a", name: "Piped A" };
  const disabledPiped = {
    ...dummyPiped,
    id: "piped-b",
    name: "Piped B",
    disabled: true,
  };
  const props = { releasedVersions: ["v0.59.0"], onClose: jest.fn() };
  const { rerender } = render(
    <UpgradePipedDialog {...props} open pipeds={[enabledPiped]} />
  );

  const initialDialog = within(screen.getByRole("dialog"));
  await userEvent.click(initialDialog.getByRole("checkbox"));
  await userEvent.click(initialDialog.getByRole("button", { name: "Cancel" }));
  expect(props.onClose).toHaveBeenCalledTimes(1);
  rerender(
    <UpgradePipedDialog {...props} open={false} pipeds={[enabledPiped]} />
  );
  await waitFor(() =>
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument()
  );

  rerender(<UpgradePipedDialog {...props} open pipeds={[disabledPiped]} />);
  const dialog = within(await screen.findByRole("dialog"));
  const checkbox = dialog.getByRole("checkbox");
  expect(dialog.getByText("Piped B")).toBeInTheDocument();
  expect(checkbox).not.toBeChecked();
  await userEvent.type(dialog.getByRole("combobox"), "v0.59.0");
  await userEvent.click(await screen.findByRole("option", { name: "v0.59.0" }));
  const submit = dialog.getByRole("button", { name: /^upgrade$/i });
  expect(submit).toBeDisabled();

  await userEvent.click(checkbox);
  await userEvent.click(submit);
  await waitFor(() => expect(upgrade).toHaveBeenCalledTimes(1));
  expect(upgrade).toHaveBeenCalledWith({
    version: "v0.59.0",
    pipedIdsList: ["piped-b"],
  });
});
