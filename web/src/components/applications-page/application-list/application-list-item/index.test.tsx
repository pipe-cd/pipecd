import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { createStore, render, screen } from "~~/test-utils";
import { ApplicationListItem } from ".";

const state = {
  applications: {
    entities: {
      [dummyApplication.id]: dummyApplication,
    },
    ids: [dummyApplication.id],
  },
};

test("delete", () => {
  const handleDelete = jest.fn();
  const store = createStore(state);
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            applicationId={dummyApplication.id}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={handleDelete}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

  expect(handleDelete).toHaveBeenCalledWith(dummyApplication.id);
});

test("edit", () => {
  const handleEdit = jest.fn();
  const store = createStore(state);
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            applicationId={dummyApplication.id}
            onEdit={handleEdit}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Edit" }));

  expect(handleEdit).toHaveBeenCalledWith(dummyApplication.id);
});

test("disable", () => {
  const handleDisable = jest.fn();
  const store = createStore(state);
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            applicationId={dummyApplication.id}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={handleDisable}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Disable" }));

  expect(handleDisable).toHaveBeenCalledWith(dummyApplication.id);
});

test("enable", () => {
  const handleEnable = jest.fn();
  const store = createStore({
    applications: {
      entities: {
        [dummyApplication.id]: { ...dummyApplication, disabled: true },
      },
      ids: [dummyApplication.id],
    },
  });
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            applicationId={dummyApplication.id}
            onEdit={() => null}
            onEnable={handleEnable}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={() => null}
          />
        </tbody>
      </table>
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Enable" }));

  expect(handleEnable).toHaveBeenCalledWith(dummyApplication.id);
});

test("Encrypt Secret", () => {
  const handleGenerateSecret = jest.fn();
  const store = createStore(state);
  render(
    <MemoryRouter>
      <table>
        <tbody>
          <ApplicationListItem
            applicationId={dummyApplication.id}
            onEdit={() => null}
            onEnable={() => null}
            onDisable={() => null}
            onDelete={() => null}
            onEncryptSecret={handleGenerateSecret}
          />
        </tbody>
      </table>
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Encrypt Secret" }));

  expect(handleGenerateSecret).toHaveBeenCalledWith(dummyApplication.id);
});
