import userEvent from "@testing-library/user-event";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen, MemoryRouter } from "~~/test-utils";
import { ApplicationListItem } from ".";

test("delete", async () => {
  const handleDelete = jest.fn();
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            app={dummyApplication}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={handleDelete}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>
  );

  await userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  await userEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

  expect(handleDelete).toHaveBeenCalled();
});

test("edit", async () => {
  const handleEdit = jest.fn();
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            app={dummyApplication}
            onEdit={handleEdit}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>
  );

  await userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  await userEvent.click(screen.getByRole("menuitem", { name: "Edit" }));

  expect(handleEdit).toHaveBeenCalled();
});

test("disable", async () => {
  const handleDisable = jest.fn();
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            app={dummyApplication}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={handleDisable}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>
  );

  await userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  await userEvent.click(screen.getByRole("menuitem", { name: "Disable" }));

  expect(handleDisable).toHaveBeenCalled();
});

test("enable", async () => {
  const handleEnable = jest.fn();

  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            app={{ ...dummyApplication, disabled: true }}
            onEdit={() => null}
            onEnable={handleEnable}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>
  );

  await userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  await userEvent.click(screen.getByRole("menuitem", { name: "Enable" }));

  expect(handleEnable).toHaveBeenCalled();
});

test("Encrypt Secret", async () => {
  const handleGenerateSecret = jest.fn();
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            app={dummyApplication}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={handleGenerateSecret}
          />
        </tbody>
      </table>
    </MemoryRouter>
  );

  await userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  await userEvent.click(
    screen.getByRole("menuitem", { name: "Encrypt Secret" })
  );

  expect(handleGenerateSecret).toHaveBeenCalled();
});
