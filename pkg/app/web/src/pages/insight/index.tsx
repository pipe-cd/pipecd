import { Box } from "@material-ui/core";
import { FC, memo, useEffect } from "react";
import { ChangeFailureRateChart } from "./components/change-failure-rate-chart";
import { DeploymentFrequencyChart } from "./components/deployment-frequency-chart";
import { LeadTimeForChangesChart } from "./components/lead-time-for-changes-chart";
import { MeanTimeToRestoreChart } from "./components/mean-time-to-restore-chart";
import { InsightHeader } from "~/components/insight-header";
import { useAppSelector, useAppDispatch } from "~/hooks/redux";
import { fetchApplications, selectById } from "~/modules/applications";
import { InsightDataPoint } from "~/modules/insight";

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const dispatch = useAppDispatch();

  const deploymentFrequency = useAppSelector<InsightDataPoint.AsObject[]>(
    (state) => state.deploymentFrequency.data
  );
  const selectedAppName = useAppSelector<string | undefined>((state) =>
    state.insight.applicationId
      ? selectById(state.applications, state.insight.applicationId)?.name
      : undefined
  );

  const data: { name: string; points: InsightDataPoint.AsObject[] }[] = [];

  if (deploymentFrequency.length > 0) {
    data.push({ name: selectedAppName || "All", points: deploymentFrequency });
  }

  useEffect(() => {
    dispatch(fetchApplications());
  }, [dispatch]);

  return (
    <Box flex={1} p={2} overflow="auto">
      <InsightHeader />
      <Box
        display="grid"
        gridGap="24px"
        gridTemplateColumns="repeat(2, 1fr)"
        mt={2}
      >
        <DeploymentFrequencyChart data={data} />
        <ChangeFailureRateChart data={[]} />
        <LeadTimeForChangesChart data={[]} />
        <MeanTimeToRestoreChart data={[]} />
      </Box>
    </Box>
  );
});
