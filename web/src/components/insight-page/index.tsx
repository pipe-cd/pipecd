import { Box } from "@mui/material";
import { FC, memo, useMemo, useState } from "react";
import { ChangeFailureRateChart } from "./change-failure-rate-chart";
import { DeploymentFrequencyChart } from "./deployment-frequency-chart";
import { InsightHeader } from "./insight-header";
import { StatisticInformation } from "./statistic-information";
import { useInsightDeploymentFrequency } from "~/queries/insight/use-insight-deployment-frequency";
import { useInsightDeploymentChangeFailureRate } from "~/queries/insight/use-insight-deployment-change-failure-rate";
import { useGetApplications } from "~/queries/applications/use-get-applications";
import {
  InsightRange,
  InsightResolution,
} from "~/queries/insight/insight.config";

export type InsightFilterValues = {
  applicationId: string;
  labels: Array<string>;
  range: InsightRange;
  resolution: InsightResolution;
};

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const { data: applications = [] } = useGetApplications();
  const [filterValues, setFilterValues] = useState<InsightFilterValues>({
    applicationId: "",
    labels: [],
    range: InsightRange.LAST_1_MONTH,
    resolution: InsightResolution.DAILY,
  });

  const selectedAppName = useMemo<string | undefined>(
    () =>
      applications.find((item) => item.id === filterValues.applicationId)?.name,
    [applications, filterValues.applicationId]
  );

  const selectedLabels = useMemo<string>(
    () =>
      filterValues.labels.length !== 0
        ? "{" + filterValues.labels.join(", ") + "}"
        : "",
    [filterValues.labels]
  );

  const { data: deploymentFrequency = [] } = useInsightDeploymentFrequency(
    {
      applicationId: filterValues.applicationId ?? "",
      labels: filterValues.labels ?? [],
      range: filterValues.range ?? InsightRange.LAST_1_MONTH,
      resolution: filterValues.resolution ?? InsightResolution.DAILY,
    },
    { keepPreviousData: true, retry: false, refetchOnWindowFocus: false }
  );

  const {
    data: deploymentChangeFailureRate = [],
  } = useInsightDeploymentChangeFailureRate(
    {
      applicationId: filterValues.applicationId ?? "",
      labels: filterValues.labels ?? [],
      range: filterValues.range ?? InsightRange.LAST_1_MONTH,
      resolution: filterValues.resolution ?? InsightResolution.DAILY,
    },
    { keepPreviousData: true, retry: false, refetchOnWindowFocus: false }
  );

  const deploymentFrequencyDataPoints = useMemo(() => {
    const name = (selectedAppName || "All") + " " + selectedLabels;
    return deploymentFrequency.length > 0
      ? [{ name, points: deploymentFrequency }]
      : [];
  }, [deploymentFrequency, selectedAppName, selectedLabels]);

  const deploymentChangeFailureRateDataPoints = useMemo(() => {
    const name = (selectedAppName || "All") + " " + selectedLabels;
    return deploymentChangeFailureRate.length > 0
      ? [{ name, points: deploymentChangeFailureRate }]
      : [];
  }, [deploymentChangeFailureRate, selectedAppName, selectedLabels]);

  const handleChangeFilter = (
    filterValues: Partial<InsightFilterValues>
  ): void => {
    setFilterValues((prev) => ({
      ...prev,
      ...filterValues,
    }));
  };

  return (
    <Box
      sx={{
        flex: 1,
        p: 2,
        overflow: "auto",
      }}
    >
      <Box
        sx={{
          display: "flex",
          flexDirection: "column",
          flex: 1,
        }}
      >
        <StatisticInformation />
      </Box>
      <InsightHeader
        filterValues={filterValues}
        onChangeFilter={handleChangeFilter}
      />
      <Box
        sx={{
          display: "grid",
          gap: "24px",
          gridTemplateColumns: "repeat(2, 1fr)",
          mt: 2,
        }}
      >
        <DeploymentFrequencyChart
          resolution={filterValues.resolution}
          data={deploymentFrequencyDataPoints}
        />
        <ChangeFailureRateChart
          resolution={filterValues.resolution}
          data={deploymentChangeFailureRateDataPoints}
        />
      </Box>
    </Box>
  );
});
