import userEvent from "@testing-library/user-event";
import { UI_TEXT_CLEAR } from "~/constants/ui-text";
import { ApplicationKind } from "~/modules/applications";
import { DeploymentStatus } from "~/modules/deployments";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen } from "~~/test-utils";
import { DeploymentFilter } from ".";

const initialState = {
  applications: {
    ids: [dummyApplication.id],
    entities: { [dummyApplication.id]: dummyApplication },
  },
};

test("Change filter values", () => {
  const onChange = jest.fn();
  render(
    <DeploymentFilter onChange={onChange} onClear={() => null} options={{}} />,
    {
      initialState,
    }
  );

  userEvent.type(
    screen.getByRole("combobox", { name: /application id/i }),
    dummyApplication.id
  );
  userEvent.click(
    screen.getByRole("option", {
      name: `${dummyApplication.name} (${dummyApplication.id})`,
    })
  );

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
    <DeploymentFilter onChange={() => null} onClear={onClear} options={{}} />,
    {
      initialState,
    }
  );

  userEvent.click(screen.getByRole("button", { name: UI_TEXT_CLEAR }));

  expect(onClear).toHaveBeenCalled();
});
