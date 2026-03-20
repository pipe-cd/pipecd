import { render, screen } from "~~/test-utils";
import { EventItem } from ".";
import { EventStatus } from "pipecd/web/model/event_pb";
import { Event } from "pipecd/web/model/event_pb";

const dummyEvent: Event.AsObject = {
  id: "event-1",
  name: "my-event",
  data: "data",
  projectId: "project-1",
  eventKey: "key-1",
  status: EventStatus.EVENT_SUCCESS,
  statusDescription: "Deployed successfully",
  labelsMap: [],
  createdAt: Math.floor(Date.now() / 1000),
  updatedAt: Math.floor(Date.now() / 1000),
  handledAt: 0,
};

describe("EventItem", () => {
  it("renders event name", () => {
    render(<EventItem event={dummyEvent} />);
    expect(screen.getByText("my-event")).toBeInTheDocument();
  });

  it("renders event status text", () => {
    render(<EventItem event={dummyEvent} />);
    expect(screen.getByText("SUCCESS")).toBeInTheDocument();
  });

  it("renders status description", () => {
    render(<EventItem event={dummyEvent} />);
    expect(screen.getByText("Deployed successfully")).toBeInTheDocument();
  });

  it("renders 'No description.' when statusDescription is empty", () => {
    render(<EventItem event={{ ...dummyEvent, statusDescription: "" }} />);
    expect(screen.getByText("No description.")).toBeInTheDocument();
  });

  it("renders event id", () => {
    render(<EventItem event={dummyEvent} />);
    expect(screen.getByText("event-1")).toBeInTheDocument();
  });

  it("renders label chips", () => {
    render(
      <EventItem
        event={{ ...dummyEvent, labelsMap: [["env", "prod"]] }}
      />
    );
    expect(screen.getByText("env: prod")).toBeInTheDocument();
  });
});
