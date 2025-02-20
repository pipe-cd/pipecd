import { fireEvent, render, screen } from "@testing-library/react";
import DeploymentTraceItem from "./index";
import { dummyDeploymentTrace } from "~/__fixtures__/dummy-deployment-trace";
import { MemoryRouter } from "~~/test-utils";

describe("DeploymentTraceItem", () => {
  it("should render trace information", () => {
    render(
      <MemoryRouter>
        <DeploymentTraceItem
          trace={dummyDeploymentTrace.trace}
          deploymentList={dummyDeploymentTrace.deploymentsList}
        />
      </MemoryRouter>
    );

    const expectedValues = {
      title: "title",
      author: "user",
      commitMessage: "commit-message",
      commitHash: "commit-hash",
      commitUrl: "/commit-url",
    };

    expect(screen.getByText(expectedValues.title)).toBeInTheDocument();
    expect(
      screen.getByText(expectedValues.author + " authored")
    ).toBeInTheDocument();
    expect(screen.getByText(expectedValues.commitMessage)).toBeInTheDocument();
    expect(screen.getByText(expectedValues.commitHash)).toBeInTheDocument();
    expect(screen.getByRole("link")).toHaveAttribute(
      "href",
      expectedValues.commitUrl
    );
  });

  it("should render deployment items", () => {
    render(
      <MemoryRouter>
        <DeploymentTraceItem
          trace={dummyDeploymentTrace.trace}
          deploymentList={dummyDeploymentTrace.deploymentsList}
        />
      </MemoryRouter>
    );
    const buttonExpand = screen.getByRole("button", { name: /expand/i });
    fireEvent.click(buttonExpand);
    expect(screen.getByText("DemoApp")).toBeInTheDocument();
  });
});
