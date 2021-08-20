import { render, screen } from "~~/test-utils";
import { ApplicationCounts } from "./index";

test("displaying application counts", () => {
  render(<ApplicationCounts onClick={jest.fn()} />, {
    initialState: {
      applicationCounts: {
        counts: {
          CLOUDRUN: {
            DISABLED: 0,
            ENABLED: 0,
          },
          CROSSPLANE: {
            DISABLED: 0,
            ENABLED: 0,
          },
          ECS: {
            DISABLED: 0,
            ENABLED: 0,
          },
          KUBERNETES: {
            DISABLED: 8,
            ENABLED: 123,
          },
          LAMBDA: {
            DISABLED: 0,
            ENABLED: 0,
          },
          TERRAFORM: {
            DISABLED: 2,
            ENABLED: 75,
          },
        },
        updatedAt: 0,
      },
    },
  });

  expect(screen.queryByText("KUBERNETES")).toBeInTheDocument();
  expect(screen.queryByText("123")).toBeInTheDocument();
  expect(screen.queryByText("/8")).toBeInTheDocument();
  expect(screen.queryByText("TERRAFORM")).toBeInTheDocument();
  expect(screen.queryByText("75")).toBeInTheDocument();
  expect(screen.queryByText("/2")).toBeInTheDocument();
  expect(screen.queryByText("CROSSPLANE")).toBeInTheDocument();
  expect(screen.queryByText("LAMBDA")).toBeInTheDocument();
  expect(screen.queryByText("CLOUDRUN")).toBeInTheDocument();
  expect(screen.queryByText("ECS")).toBeInTheDocument();
});
