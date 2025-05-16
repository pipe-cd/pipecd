import { Box } from "@mui/material";
import { FC, memo, useEffect } from "react";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchApplications, selectById } from "~/modules/applications";
import { fetchApplicationCount } from "~/modules/application-counts";
import {
  InsightDataPoint,
  InsightResolution,
  InsightRange,
} from "~/modules/insight";
import { ChangeFailureRateChart } from "./change-failure-rate-chart";
import { DeploymentFrequencyChart } from "./deployment-frequency-chart";
import { InsightHeader } from "./insight-header";
import {
  fetchDeploymentChangeFailureRate,
  fetchDeploymentChangeFailureRate24h,
} from "~/modules/deployment-change-failure-rate";
import {
  fetchDeployment24h,
  fetchDeploymentFrequency,
} from "~/modules/deployment-frequency";
import { StatisticInformation } from "./statistic-information";

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const dispatch = useAppDispatch();

  const [applicationId, labels, range, resolution] = useAppSelector<
    [string, Array<string>, InsightRange, InsightResolution]
  >((state) => [
    state.insight.applicationId,
    state.insight.labels,
    state.insight.range,
    state.insight.resolution,
  ]);

  const selectedAppName = useAppSelector<string | undefined>((state) =>
    state.insight.applicationId
      ? selectById(state.applications, state.insight.applicationId)?.name
      : undefined
  );

  const selectedLabels = useAppSelector<string>((state) =>
    state.insight.labels.length !== 0
      ? "{" + state.insight.labels.join(", ") + "}"
      : ""
  );

  const deploymentFrequency = useAppSelector<InsightDataPoint.AsObject[]>(
    (state) => state.deploymentFrequency.data
  );
  const deploymentFrequencyDataPoints: {
    name: string;
    points: InsightDataPoint.AsObject[];
  }[] = [];
  if (deploymentFrequency.length > 0) {
    deploymentFrequencyDataPoints.push({
      name: (selectedAppName || "All") + " " + selectedLabels,
      points: deploymentFrequency,
    });
  }

  const deploymentChangeFailureRate = useAppSelector<
    InsightDataPoint.AsObject[]
  >((state) => state.deploymentChangeFailureRate.data);
  const deploymentChangeFailureRateDataPoints: {
    name: string;
    points: InsightDataPoint.AsObject[];
  }[] = [];
  if (deploymentChangeFailureRate.length > 0) {
    deploymentChangeFailureRateDataPoints.push({
      name: (selectedAppName || "All") + " " + selectedLabels,
      points: deploymentChangeFailureRate,
    });
  }

  useEffect(() => {
    dispatch(fetchApplications());
    dispatch(fetchApplicationCount());
  }, [dispatch]);

  useEffect(() => {
    dispatch(fetchDeploymentFrequency());
    dispatch(fetchDeployment24h());
    dispatch(fetchDeploymentChangeFailureRate());
    dispatch(fetchDeploymentChangeFailureRate24h());
  }, [dispatch, applicationId, labels, range, resolution]);

  return (
    <Box flex={1} p={2} overflow="auto">
      <Box display="flex" flexDirection="column" flex={1} p={2}>
        <StatisticInformation />
      </Box>
      <InsightHeader />
      <Box
        display="grid"
        gap="24px"
        gridTemplateColumns="repeat(2, 1fr)"
        mt={2}
      >
        <DeploymentFrequencyChart
          resolution={resolution}
          data={deploymentFrequencyDataPoints}
        />
        <ChangeFailureRateChart
          resolution={resolution}
          data={deploymentChangeFailureRateDataPoints}
        />
      </Box>
    </Box>
  );
});
