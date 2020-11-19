import userEvent from "@testing-library/user-event";
import React from "react";
import { render, screen } from "../../test-utils";
import { dummyApplication } from "../__fixtures__/dummy-application";
import { SealedSecretDialog } from "./sealed-secret-dialog";

test("render", () => {
  render(
    <SealedSecretDialog
      applicationId={dummyApplication.id}
      onClose={() => null}
      open
    />,
    {
      initialState: {
        applications: {
          entities: {
            [dummyApplication.id]: dummyApplication,
          },
          ids: [dummyApplication.id],
        },
      },
    }
  );

  expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
});

test("cancel", () => {
  const onClose = jest.fn();
  render(
    <SealedSecretDialog
      applicationId={dummyApplication.id}
      onClose={onClose}
      open
    />,
    {
      initialState: {
        applications: {
          entities: {
            [dummyApplication.id]: dummyApplication,
          },
          ids: [dummyApplication.id],
        },
      },
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Cancel" }));
  expect(onClose).toHaveBeenCalled();
});
