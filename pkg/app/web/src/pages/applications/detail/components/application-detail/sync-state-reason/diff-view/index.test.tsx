import { render, screen } from "~~/test-utils";
import { DiffView } from ".";

it("should render normal text", () => {
  render(<DiffView content="normal text" />, {});

  expect(screen.queryByTestId("added-line")).not.toBeInTheDocument();
  expect(screen.queryByTestId("deleted-line")).not.toBeInTheDocument();
});

it("should render line as added line if the line start with'+'", () => {
  render(<DiffView content="+ added-line" />, {});

  expect(screen.getByTestId("added-line")).toBeInTheDocument();
  expect(screen.queryByTestId("deleted-line")).not.toBeInTheDocument();
});

it("should render line as deleted line if the line start with '-'", () => {
  render(<DiffView content="- deleted line" />, {});

  expect(screen.queryByTestId("added-line")).not.toBeInTheDocument();
  expect(screen.getByTestId("deleted-line")).toBeInTheDocument();
});
