import { render, screen } from "~~/test-utils";
import { MemoryRouter } from "~~/test-utils";
import { DeploymentItem } from ".";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";

describe("DeploymentItem", () => {
  it("renders null when deployment is undefined", () => {
    const { container } = render(
      <MemoryRouter>
        <DeploymentItem />
      </MemoryRouter>
    );
    expect(container.firstChild).toBeNull();
  });

  it("renders application name and status", () => {
    render(
      <MemoryRouter>
        <DeploymentItem deployment={dummyDeployment} />
      </MemoryRouter>
    );
    expect(screen.getByText(dummyDeployment.applicationName)).toBeInTheDocument();
    expect(screen.getByText("SUCCESS")).toBeInTheDocument();
  });

  it("renders deployment summary", () => {
    render(
      <MemoryRouter>
        <DeploymentItem deployment={dummyDeployment} />
      </MemoryRouter>
    );
    expect(screen.getByText(dummyDeployment.summary)).toBeInTheDocument();
  });

  it("renders 'No description.' when summary is empty", () => {
    render(
      <MemoryRouter>
        <DeploymentItem deployment={{ ...dummyDeployment, summary: "" }} />
      </MemoryRouter>
    );
    expect(screen.getByText("No description.")).toBeInTheDocument();
  });

  it("links to the correct deployment detail page", () => {
    render(
      <MemoryRouter>
        <DeploymentItem deployment={dummyDeployment} />
      </MemoryRouter>
    );
    const link = screen.getByRole("button");
    expect(link).toHaveAttribute(
      "href",
      `${PAGE_PATH_DEPLOYMENTS}/${dummyDeployment.id}`
    );
  });

  it("renders KUBERNETES kind for v0 piped", () => {
    render(
      <MemoryRouter>
        <DeploymentItem
          deployment={{
            ...dummyDeployment,
            platformProvider: "kube-1",
            deployTargetsByPluginMap: [],
          }}
        />
      </MemoryRouter>
    );
    expect(screen.getByText("KUBERNETES")).toBeInTheDocument();
  });

  it("renders APPLICATION for v1 piped", () => {
    render(
      <MemoryRouter>
        <DeploymentItem
          deployment={{
            ...dummyDeployment,
            platformProvider: "",
            deployTargetsByPluginMap: [],
          }}
        />
      </MemoryRouter>
    );
    expect(screen.getByText("APPLICATION")).toBeInTheDocument();
  });
});
