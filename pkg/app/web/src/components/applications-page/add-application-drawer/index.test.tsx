import { server } from "~/mocks/server";

jest.setTimeout(50_000);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("AddApplicationDrawer", () => {
  it("should pass");

  // TODO: Move all these tests into application-form component
  //
  // it("should calls onSubmit if clicked SAVE button", async () => {
  //   const store = createStore({
  //     pipeds: {
  //       entities: { [dummyPiped.id]: dummyPiped },
  //       ids: [dummyPiped.id],
  //     },
  //     environments: {
  //       entities: { [dummyEnv.id]: dummyEnv },
  //       ids: [dummyEnv.id],
  //     },
  //   });
  //   render(<AddApplicationDrawer open onClose={jest.fn()} />, {
  //     store,
  //   });

  //   userEvent.type(screen.getByRole("textbox", { name: "Name" }), "App");
  //   userEvent.click(screen.getByRole("button", { name: /Kind/ }));
  //   userEvent.click(
  //     screen.getByRole("option", {
  //       name: APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
  //     })
  //   );

  //   userEvent.click(screen.getByRole("button", { name: /Environment/ }));
  //   userEvent.click(screen.getByRole("option", { name: dummyEnv.name }));

  //   userEvent.click(screen.getByRole("button", { name: /Piped/ }));
  //   userEvent.click(screen.getByRole("option", { name: /dummy-piped/ }));

  //   userEvent.click(screen.getByRole("button", { name: /Repository/ }));
  //   userEvent.click(screen.getByRole("option", { name: /debug-repo/ }));

  //   userEvent.type(screen.getByRole("textbox", { name: "Path" }), "path");

  //   userEvent.click(screen.getByRole("button", { name: /Cloud Provider/ }));
  //   userEvent.click(screen.getByRole("option", { name: /terraform-default/ }));

  //   userEvent.click(screen.getByRole("button", { name: UI_TEXT_SAVE }));

  //   await waitFor(() =>
  //     expect(store.getActions()).toEqual(
  //       expect.arrayContaining([
  //         {
  //           type: addApplication.pending.type,
  //           meta: expect.objectContaining({
  //             arg: {
  //               cloudProvider: "terraform-default",
  //               configFilename: "app.pipecd.yaml",
  //               env: dummyEnv.id,
  //               kind: ApplicationKind.TERRAFORM,
  //               name: "App",
  //               pipedId: dummyPiped.id,
  //               repo: {
  //                 branch: "master",
  //                 id: "debug-repo",
  //                 remote: "git@github.com:pipe-cd/debug.git",
  //               },
  //               repoPath: "path",
  //               labels: [],
  //             },
  //           }),
  //         },
  //       ])
  //     )
  //   );

  //   await waitFor(() => {
  //     expect(screen.queryByDisplayValue("App")).not.toBeInTheDocument();
  //   });
  // });

  // it("should clear depended fields if change environment", async () => {
  //   const altEnv = { ...dummyEnv, id: "env-2", name: "env-2" };
  //   render(<AddApplicationDrawer open />, {
  //     initialState: {
  //       pipeds: {
  //         entities: { [dummyPiped.id]: dummyPiped },
  //         ids: [dummyPiped.id],
  //       },
  //       environments: {
  //         entities: { [dummyEnv.id]: dummyEnv, [altEnv.id]: altEnv },
  //         ids: [dummyEnv.id, altEnv.id],
  //       },
  //     },
  //   });

  //   await waitFor(() =>
  //     expect(screen.getByRole("button", { name: UI_TEXT_SAVE })).toBeDisabled()
  //   );

  //   userEvent.type(screen.getByRole("textbox", { name: "Name" }), "App");

  //   userEvent.click(screen.getByRole("button", { name: /kind/i }));
  //   userEvent.click(
  //     screen.getByRole("option", {
  //       name: APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
  //     })
  //   );

  //   userEvent.click(screen.getByRole("button", { name: /Environment/ }));
  //   userEvent.click(screen.getByRole("option", { name: dummyEnv.name }));

  //   userEvent.click(screen.getByRole("button", { name: /Piped/ }));
  //   userEvent.click(screen.getByRole("option", { name: /dummy-piped/i }));

  //   userEvent.click(screen.getByRole("button", { name: /Repository/ }));
  //   userEvent.click(screen.getByRole("option", { name: /debug-repo/i }));

  //   userEvent.type(screen.getByRole("textbox", { name: /Path/ }), "path");

  //   expect(screen.getByRole("button", { name: /Piped/ })).toHaveTextContent(
  //     /dummy-piped/
  //   );
  //   expect(
  //     screen.getByRole("button", { name: /Repository/ })
  //   ).toHaveTextContent(/debug-repo/);
  //   expect(screen.getByRole("textbox", { name: /Path/ })).toHaveDisplayValue(
  //     "path"
  //   );

  //   userEvent.click(screen.getByRole("button", { name: /Environment/ }));
  //   userEvent.click(screen.getByRole("option", { name: altEnv.name }));

  //   expect(screen.getByRole("button", { name: /Piped/ })).not.toHaveTextContent(
  //     /dummy-piped/
  //   );
  //   expect(screen.getByRole("textbox", { name: /Path/ })).toHaveDisplayValue(
  //     ""
  //   );
  //   expect(
  //     screen.getByRole("button", { name: /Repository/ })
  //   ).not.toHaveTextContent(/debug-repo/);
  // });

  // it("should calls onClose handler if clicked CANCEL button", () => {
  //   const onClose = jest.fn();
  //   render(<AddApplicationDrawer open onClose={onClose} />, {});

  //   userEvent.click(screen.getByRole("button", { name: UI_TEXT_CANCEL }));

  //   expect(onClose).toHaveBeenCalled();
  // });
});
