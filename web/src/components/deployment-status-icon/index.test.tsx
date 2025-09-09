import { render, screen } from "~~/test-utils";
import { DeploymentStatus } from "~/types/deployment";
import { DeploymentStatusIcon } from "./";

test("DEPLOYMENT_CANCELLED", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_CANCELLED} />,
    {}
  );

  expect(screen.getByTestId("deployment-cancel-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_FAILURE", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_FAILURE} />,
    {}
  );

  expect(screen.getByTestId("deployment-error-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_PENDING", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_PENDING} />,
    {}
  );

  expect(screen.getByTestId("deployment-pending-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_PLANNED", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_PLANNED} />,
    {}
  );

  expect(screen.getByTestId("deployment-pending-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_ROLLING_BACK", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_ROLLING_BACK} />,
    {}
  );

  expect(screen.getByTestId("deployment-rollback-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_RUNNING", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_RUNNING} />,
    {}
  );

  expect(screen.getByTestId("deployment-running-icon")).toBeInTheDocument();
});

test("DEPLOYMENT_SUCCESS", () => {
  render(
    <DeploymentStatusIcon status={DeploymentStatus.DEPLOYMENT_SUCCESS} />,
    {}
  );

  expect(screen.getByTestId("deployment-success-icon")).toBeInTheDocument();
});
