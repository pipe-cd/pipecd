import React from "react";
import { AddApplicationDrawer } from "./add-application-drawer";
import { render } from "../../test-utils";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../constants/ui-text";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { ApplicationKind } from "../modules/applications";
import { dummyPiped } from "../__fixtures__/dummy-piped";
import userEvent from "@testing-library/user-event";

describe("AddApplicationDrawer", () => {
  it("calls onSubmit if clicked SAVE button", async () => {
    const onSubmit = jest.fn();
    const { getByRole } = render(
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

    userEvent.type(getByRole("textbox", { name: "Name" }), "App");
    userEvent.click(getByRole("button", { name: /Kind/ }));
    userEvent.click(
      getByRole("option", {
        name: APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
      })
    );
    userEvent.click(getByRole("button", { name: /Environment/ }));
    userEvent.click(getByRole("option", { name: dummyEnv.name }));

    userEvent.click(getByRole("button", { name: /Piped/ }));
    userEvent.click(getByRole("option", { name: /dummy piped/ }));

    userEvent.click(getByRole("button", { name: /Repository/ }));
    userEvent.click(getByRole("option", { name: /debug-repo/ }));

    userEvent.type(getByRole("textbox", { name: "Path" }), "path");

    userEvent.click(getByRole("button", { name: /Cloud Provider/ }));
    userEvent.click(getByRole("option", { name: /terraform-default/ }));

    userEvent.click(getByRole("button", { name: UI_TEXT_SAVE }));

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
    });
  });

  it("calls onClose handler if clicked CANCEL button", () => {
    const onClose = jest.fn();
    const { getByRole } = render(
      <AddApplicationDrawer
        open
        projectName="pipecd"
        onSubmit={jest.fn()}
        onClose={onClose}
        isAdding={false}
      />,
      {}
    );

    userEvent.click(getByRole("button", { name: UI_TEXT_CANCEL }));

    expect(onClose).toHaveBeenCalled();
  });
});
