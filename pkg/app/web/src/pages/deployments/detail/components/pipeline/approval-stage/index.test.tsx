import userEvent from "@testing-library/user-event";
import { render, screen } from "~~/test-utils";
import { ApprovalStage } from ".";

it("shows stage name", () => {
  const { container } = render(
    <ApprovalStage
      id="stageId"
      name="APPROVAL_STAGE"
      onClick={jest.fn()}
      active
    />,
    {}
  );

  expect(container).toHaveTextContent("APPROVAL_STAGE");
});

it("calls onClick handler if clicked component", () => {
  const onClick = jest.fn();
  render(
    <ApprovalStage
      id="stageId"
      name="APPROVAL_STAGE"
      onClick={onClick}
      active
    />,
    {}
  );

  userEvent.click(screen.getByText("APPROVAL_STAGE"));
  expect(onClick).toHaveBeenCalledWith("stageId", "APPROVAL_STAGE");
});
