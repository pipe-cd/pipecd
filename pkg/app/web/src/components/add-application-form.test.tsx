import React from "react";
import { AddApplicationForm } from "./add-application-form";
import { render, fireEvent } from "../../test-utils";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../constants/ui-text";
import { dummyEnv } from "../__fixtures__/dummy-environment";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { ApplicationKind } from "../modules/applications";
import { dummyPiped } from "../__fixtures__/dummy-piped";

describe("AddApplicationForm", () => {
  it("calls onSubmit if clicked SAVE button", async () => {
    const onSubmit = jest.fn();
    const { getByRole } = render(
      <AddApplicationForm
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

    fireEvent.change(getByRole("textbox", { name: "Name" }), {
      target: {
        value: "App",
      },
    });
    fireEvent.mouseDown(getByRole("button", { name: /Kind/ }));
    fireEvent.click(
      getByRole("option", {
        name: APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM],
      })
    );
    fireEvent.mouseDown(getByRole("button", { name: /Environment/ }));
    fireEvent.click(getByRole("option", { name: dummyEnv.name }));

    fireEvent.mouseDown(getByRole("button", { name: /Piped/ }));
    fireEvent.click(getByRole("option", { name: /dummy piped/ }));

    fireEvent.mouseDown(getByRole("button", { name: /Repository/ }));
    fireEvent.click(getByRole("option", { name: /debug-repo/ }));

    fireEvent.change(getByRole("textbox", { name: "Path" }), {
      target: {
        value: "path",
      },
    });

    fireEvent.mouseDown(getByRole("button", { name: /Cloud Provider/ }));
    fireEvent.click(getByRole("option", { name: /terraform-default/ }));

    fireEvent.click(getByRole("button", { name: UI_TEXT_SAVE }));

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
      <AddApplicationForm
        projectName="pipecd"
        onSubmit={jest.fn()}
        onClose={onClose}
        isAdding={false}
      />,
      {}
    );

    fireEvent.click(getByRole("button", { name: UI_TEXT_CANCEL }));

    expect(onClose).toHaveBeenCalled();
  });
});
