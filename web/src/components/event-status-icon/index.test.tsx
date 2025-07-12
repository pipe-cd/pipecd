import { render } from "~~/test-utils";
import { EventStatusIcon } from ".";
import { EventStatus } from "~~/model/event_pb";

describe("EventStatusIcon", () => {
  it("renders the correct status icon", () => {
    const { getByTestId } = render(
      <EventStatusIcon status={EventStatus.EVENT_NOT_HANDLED} />
    );
    expect(getByTestId("event-not-handled-icon")).toBeInTheDocument();
  });
  it("renders the success icon", () => {
    const { getByTestId } = render(
      <EventStatusIcon status={EventStatus.EVENT_SUCCESS} />
    );
    expect(getByTestId("event-success-icon")).toBeInTheDocument();
  });
  it("renders the failure icon", () => {
    const { getByTestId } = render(
      <EventStatusIcon status={EventStatus.EVENT_FAILURE} />
    );
    expect(getByTestId("event-failure-icon")).toBeInTheDocument();
  });
  it("renders the outdated icon", () => {
    const { getByTestId } = render(
      <EventStatusIcon status={EventStatus.EVENT_OUTDATED} />
    );
    expect(getByTestId("event-outdated-icon")).toBeInTheDocument();
  });
});
