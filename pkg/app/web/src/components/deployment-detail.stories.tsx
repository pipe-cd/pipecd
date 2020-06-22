import React from "react";
import { DeploymentDetail } from "./deployment-detail";
import { DeploymentStatus } from "pipe/pkg/app/web/model/deployment_pb";

export default {
  title: "DeploymentDetail",
  component: DeploymentDetail,
};

export const overview: React.FC = () => (
  <DeploymentDetail
    name="deployment-1"
    env="production"
    pipedId="piped-1"
    status={DeploymentStatus.DEPLOYMENT_SUCCESS}
    description="This deployment is debug"
    commit={{
      message: "Add a description field to deployment model",
      author: "awesome-user",
      branch: "fix-bug",
      createdAt: 0,
      hash: "1234abcd",
      pullRequest: 0,
    }}
  />
);
