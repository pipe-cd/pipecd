import userEvent from "@testing-library/user-event";
import { UI_TEXT_CLEAR } from "~/constants/ui-text";
import { ApplicationKind } from "~/types/applications";
import { DeploymentStatus } from "~/types/deployment";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen, waitFor } from "~~/test-utils";
import { DeploymentFilter } from ".";
import { setupServer } from "msw/node";
import { listApplicationsHandler } from "~/mocks/services/application";

const server = setupServer(listApplicationsHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("Change filter values", async () => {
  const onChange = jest.fn();
  render(
    <DeploymentFilter options={{}} onChange={onChange} onClear={() => null} />
  );

  userEvent.type(
    screen.getByRole("combobox", { name: /application id/i }),
    dummyApplication.id
  );

  await waitFor(() => {
    userEvent.click(
      screen.getByRole("option", {
        name: `${dummyApplication.name} (${dummyApplication.id})`,
      })
    );
  });

  expect(onChange).toHaveBeenCalledWith({ applicationId: dummyApplication.id });
  onChange.mockClear();

  userEvent.click(screen.getByRole("combobox", { name: /application kind/i }));
  userEvent.click(screen.getByRole("option", { name: /kubernetes/i }));

  expect(onChange).toHaveBeenCalledWith({
    kind: `${ApplicationKind.KUBERNETES}`,
  });
  onChange.mockClear();

  userEvent.click(screen.getByRole("combobox", { name: /deployment status/i }));
  userEvent.click(screen.getByRole("option", { name: /success/i }));

  expect(onChange).toHaveBeenCalledWith({
    status: `${DeploymentStatus.DEPLOYMENT_SUCCESS}`,
  });
});

test("Click clear filter", () => {
  const onClear = jest.fn();
  render(
    <DeploymentFilter onChange={() => null} onClear={onClear} options={{}} />
  );

  userEvent.click(screen.getByRole("button", { name: UI_TEXT_CLEAR }));

  expect(onClear).toHaveBeenCalled();
});
