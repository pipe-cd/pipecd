import { Story } from "@storybook/react";
import { METADATA_APPROVED_BY } from "~/constants/metadata-keys";
import { Deployment, SyncStrategy } from "~/modules/deployments";
import { createPipelineStage } from "~/__fixtures__/dummy-pipeline";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { Pipeline, PipelineProps } from ".";

const DEPLOYMENT_ID = "debug-deployment-id-01";
const fakeDeployment: Deployment.AsObject = {
  id: DEPLOYMENT_ID,
  applicationId: "debug-project/development/debug-app",
  applicationName: "demo-app",
  envId: "development",
  pipedId: "debug-piped",
  projectId: "debug-project",
  kind: 0,
  gitPath: {
    configPath: "",
    repo: {
      id: "pipe-debug",
      branch: "master",
      remote: "xxx",
    },
    path: "k8s",
    configFilename: "",
    url: "",
  },
  version: "0.0.1",
  cloudProvider: "",
  labelsMap: [],
  trigger: {
    commit: {
      hash: "3808585b46f1e90196d7ffe8dd04c807a251febc",
      message: "Add web page routing (#133)",
      author: "cakecatz",
      branch: "master",
      pullRequest: 0,
      createdAt: 1592201366,
      url: "",
    },
    commander: "cakecatz",
    timestamp: 1592201366,
    syncStrategy: SyncStrategy.AUTO,
    strategySummary: "",
  },
  runningCommitHash: "3808585b46f1e90196d7ffe8dd04c807a251febc",
  summary: "This deployment is debug",
  status: 2,
  statusReason: "",
  stagesList: [
    createPipelineStage({ id: "fake-stage-id-0-0" }),
    createPipelineStage({
      index: 1,
      id: "fake-stage-id-1-0",
      requiresList: ["fake-stage-id-0-0"],
      status: 1,
    }),
    createPipelineStage({
      id: "fake-stage-id-1-1",
      name: "WAIT_APPROVAL",
      index: 2,
      requiresList: ["fake-stage-id-0-0"],
      status: 2,
      metadataMap: [[METADATA_APPROVED_BY, "User"]],
    }),
    createPipelineStage({
      id: "fake-stage-id-1-2",
      name: "K8S_CANARY_ROLLOUT",
      index: 2,
      requiresList: ["fake-stage-id-0-0"],
      status: 3,
    }),
    createPipelineStage({
      id: "fake-stage-id-1-3",
      name: "WAIT_APPROVAL",
      index: 2,
      requiresList: ["fake-stage-id-0-0"],
      status: 1,
    }),
    createPipelineStage({
      id: "fake-stage-id-2-0",
      name: "K8S_TRAFFIC_ROUTING",
      desc: "waiting approval",
      index: 0,
      requiresList: [
        "fake-stage-id-1-0",
        "fake-stage-id-1-1",
        "fake-stage-id-1-2",
      ],
      status: 0,
      metadataMap: [
        ["baseline-percentage", "0"],
        ["canary-percentage", "50"],
        ["primary-percentage", "50"],
      ],
    }),
    createPipelineStage({
      id: "fake-stage-id-2-1",
      name: "K8S_CANARY_CLEAN",
      desc: "approved by cakecatz",
      index: 1,
      requiresList: [
        "fake-stage-id-1-0",
        "fake-stage-id-1-1",
        "fake-stage-id-1-2",
      ],
      status: 0,
    }),
    createPipelineStage({
      id: "fake-stage-id-3-0",
      name: "K8S_CANARY_ROLLOUT",
      index: 0,
      requiresList: ["fake-stage-id-2-0", "fake-stage-id-2-1"],
      status: 0,
    }),
  ],
  metadataMap: [],
  completedAt: 0,
  createdAt: 1592203166,
  updatedAt: 1592203166,
  deploymentChainId: "",
  deploymentChainBlockIndex: 0,
};

export default {
  title: "DEPLOYMENT/Pipeline/Pipeline",
  component: Pipeline,
  decorators: [
    createDecoratorRedux({
      deployments: {
        ids: [DEPLOYMENT_ID],
        entities: {
          [DEPLOYMENT_ID]: fakeDeployment,
        },
      },
    }),
  ],
};

const Template: Story<PipelineProps> = (args) => <Pipeline {...args} />;
export const Overview = Template.bind({});
Overview.args = { deploymentId: DEPLOYMENT_ID };
