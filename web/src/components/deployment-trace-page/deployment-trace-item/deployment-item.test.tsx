import DeploymentItem from "./deployment-item";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { render } from "~~/test-utils";

describe("DeploymentItem", () => {
  it("should render deployment item with correct data", () => {
    const { getByText } = render(
      <DeploymentItem deployment={dummyDeployment} />,
      {}
    );
    const expectedValues = {
      status: "SUCCESS",
      applicationName: "DemoApp",
      kind: "KUBERNETES",
      description:
        "Quick sync by deploying the new version and configuring all traffic to it because no pipeline was configured",
    };

    Object.values(expectedValues).forEach((value) => {
      expect(getByText(value)).toBeInTheDocument();
    });
  });

  it("should display 'No description.' when summary is empty", () => {
    const deploymentWithoutSummary = { ...dummyDeployment, summary: "" };
    const { getByText } = render(
      <DeploymentItem deployment={deploymentWithoutSummary} />,
      {}
    );

    expect(getByText("No description.")).toBeInTheDocument();
  });
});
