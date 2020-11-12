import React from "react";
import { AddApplicationDrawer } from "./add-application-drawer";
import { render, screen, waitFor } from "../../test-utils";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../constants/ui-text";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { ApplicationKind } from "../modules/applications";
import { dummyPiped } from "../__fixtures__/dummy-piped";
import userEvent from "@testing-library/user-event";

describe("AddApplicationDrawer", () => {
  it("calls onSubmit if clicked SAVE button", async () => {
    const onSubmit = jest.fn();
    render(
      <AddApplicationDrawer
        open
        projectName="pipecd"
        onSubmit={onSubmit}
        onClose={jest.fn()}
        isAdding={false}
      />,
      {
        initialState: {
          pipeds: {
            entities: { [dummyPiped.id]: dummyPiped },
            ids: [dummyPiped.id],
          },
          environments: {
            entities: { [dummyEnv.id]: dummyEnv },
            ids: [dummyEnv.id],
          },
        },
      }
    );

    userEvent.type(screen.getByRole("textbox", { name: "Name" }), "App");

    userEvent.click(screen.getByRole("button", { name: /Kind/ }));
    userEvent.click(
      screen.getByRole("option", {
        name: APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
      })
    );

    userEvent.click(screen.getByRole("button", { name: /Environment/ }));
    userEvent.click(screen.getByRole("option", { name: dummyEnv.name }));

    userEvent.click(screen.getByRole("button", { name: /Piped/ }));
    userEvent.click(screen.getByRole("option", { name: /dummy piped/ }));

    userEvent.click(screen.getByRole("button", { name: /Repository/ }));
    userEvent.click(screen.getByRole("option", { name: /debug-repo/ }));

    userEvent.type(screen.getByRole("textbox", { name: "Path" }), "path");

    userEvent.click(screen.getByRole("button", { name: /Cloud Provider/ }));
    userEvent.click(screen.getByRole("option", { name: /terraform-default/ }));

    userEvent.click(screen.getByRole("button", { name: UI_TEXT_SAVE }));

    await waitFor(() =>
      expect(onSubmit).toHaveBeenCalledWith({
        cloudProvider: "terraform-default",
        configFilename: "",
        env: "env-1",
        kind: ApplicationKind.TERRAFORM,
        name: "App",
        pipedId: "piped-1",
        repo: {
          branch: "master",
          id: "debug-repo",
          remote: "git@github.com:pipe-cd/debug.git",
        },
        repoPath: "path",
      })
    );
  });

  it("calls onClose handler if clicked CANCEL button", () => {
    const onClose = jest.fn();
    render(
      <AddApplicationDrawer
        open
        projectName="pipecd"
        onSubmit={jest.fn()}
        onClose={onClose}
        isAdding={false}
      />,
      {}
    );

    userEvent.click(screen.getByRole("button", { name: UI_TEXT_CANCEL }));

    expect(onClose).toHaveBeenCalled();
  });
});
