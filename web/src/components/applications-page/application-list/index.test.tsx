import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import * as applicationApi from "~/api/applications";
import * as pipedApi from "~/api/piped";
import { render, screen, waitFor, MemoryRouter } from "~~/test-utils";
import { ApplicationList } from ".";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("delete", async () => {
  jest.spyOn(applicationApi, "deleteApplication");

  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} applications={[dummyApplication]} />
    </MemoryRouter>
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

  await waitFor(() =>
    expect(
      screen.getByRole("heading", { name: "Delete Application" })
    ).toBeInTheDocument()
  );

  userEvent.click(screen.getByRole("button", { name: "Delete" }));

  await waitFor(() =>
    expect(
      screen.queryByRole("heading", { name: "Delete Application" })
    ).not.toBeInTheDocument()
  );
  expect(applicationApi.deleteApplication).toHaveBeenCalledWith({
    applicationId: dummyApplication.id,
  });
});

test("show specific page", async () => {
  const apps = [...new Array(30)].map((_, i) => ({
    ...dummyApplication,
    id: `${dummyApplication.id}${i}`,
  }));
  render(
    <MemoryRouter>
      <ApplicationList currentPage={2} applications={apps} />
    </MemoryRouter>
  );

  const items = await screen.findAllByText(dummyApplication.name);
  expect(items).toHaveLength(10);
});

test("edit", async () => {
  render(
    <MemoryRouter>
      <ApplicationList
        currentPage={1}
        applications={[
          {
            ...dummyApplication,
            pipedId: "",
            platformProvider: "",
            gitPath: {
              configFilename: "",
              path: "dir/dir1",
              url: "",
              repo: { id: "", branch: "", remote: "" },
            },
          },
        ]}
      />
    </MemoryRouter>
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Edit" }));

  await waitFor(() =>
    expect(screen.getByRole("textbox", { name: "Name" })).toHaveValue(
      dummyApplication.name
    )
  );
});

test("disable", async () => {
  jest.spyOn(applicationApi, "disableApplication");
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} applications={[dummyApplication]} />
    </MemoryRouter>
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Disable" }));

  userEvent.click(screen.getByRole("button", { name: "Disable" }));

  await waitFor(() => {
    expect(applicationApi.disableApplication).toHaveBeenCalledWith({
      applicationId: dummyApplication.id,
    });
  });

  await waitFor(() => {
    expect(
      screen.queryByRole("button", { name: "Disable" })
    ).not.toBeInTheDocument();
  });
});

test("enable", async () => {
  jest.spyOn(applicationApi, "enableApplication");

  render(
    <MemoryRouter>
      <ApplicationList
        currentPage={1}
        applications={[{ ...dummyApplication, disabled: true }]}
      />
    </MemoryRouter>
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Enable" }));

  await waitFor(() => {
    expect(applicationApi.enableApplication).toHaveBeenCalledWith({
      applicationId: dummyApplication.id,
    });
  });
});

test("Encrypt Secret", async () => {
  jest.spyOn(pipedApi, "generateApplicationSealedSecret");
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} applications={[dummyApplication]} />
    </MemoryRouter>
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Encrypt Secret" }));

  userEvent.type(
    screen.getByRole("textbox", { name: "Secret Data" }),
    "secret data"
  );

  userEvent.click(screen.getByRole("button", { name: "Encrypt" }));

  await waitFor(() => {
    expect(pipedApi.generateApplicationSealedSecret).toHaveBeenCalledWith({
      base64Encoding: false,
      data: "secret data",
      pipedId: dummyApplication.pipedId,
    });
  });

  await waitFor(() => {
    expect(screen.getByText("Encrypted secret data")).toBeInTheDocument();
  });
});
