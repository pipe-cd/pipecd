import React from "react";
import { ApplicationDetail } from "./application-detail";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export default {
  title: "ApplicationDetail",
  component: ApplicationDetail,
};

export const overview: React.FC = () => (
  <ApplicationDetail
    name="DemoApp"
    env="production"
    piped="awesome-piped"
    status={ApplicationSyncStatus.SYNCED}
    deployedAt={1591671036493}
    version="v0.0.1"
    deploymentId="hello12345"
    description="description description description"
  />
);
