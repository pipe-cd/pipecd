import React from "react";
import { StageStatusIcon } from "./";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";

export default {
  title: "DEPLOYMENT/StageStatusIcon",
  component: StageStatusIcon,
};

export const overview: React.FC = () => (
  <>
    <StageStatusIcon status={StageStatus.STAGE_CANCELLED} />
    <StageStatusIcon status={StageStatus.STAGE_FAILURE} />
    <StageStatusIcon status={StageStatus.STAGE_NOT_STARTED_YET} />
    <StageStatusIcon status={StageStatus.STAGE_RUNNING} />
    <StageStatusIcon status={StageStatus.STAGE_SUCCESS} />
  </>
);
