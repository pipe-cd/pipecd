import { Box } from "@material-ui/core";
import { FC, memo, useEffect, useCallback } from "react";
import { useHistory } from "react-router-dom";
import { useAppDispatch } from "~/hooks/redux";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import {
  ApplicationKind,
  fetchApplications,
  // selectById,
} from "~/modules/applications";
import { fetchApplicationCount } from "~/modules/application-counts";
import { InsightDataPoint } from "~/modules/insight";
import { ApplicationCounts } from "./application-counts";
import { ChangeFailureRateChart } from "./change-failure-rate-chart";
import { DeploymentFrequencyChart } from "./deployment-frequency-chart";
// import { InsightHeader } from "./insight-header";
import { LeadTimeForChangesChart } from "./lead-time-for-changes-chart";
import { MeanTimeToRestoreChart } from "./mean-time-to-restore-chart";

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const dispatch = useAppDispatch();
  const history = useHistory();

  // const deploymentFrequency = useAppSelector<InsightDataPoint.AsObject[]>(
  //   (state) => state.deploymentFrequency.data
  // );
  // const selectedAppName = useAppSelector<string | undefined>((state) =>
  //   state.insight.applicationId
  //     ? selectById(state.applications, state.insight.applicationId)?.name
  //     : undefined
  // );

  const data: { name: string; points: InsightDataPoint.AsObject[] }[] = [];

  // if (deploymentFrequency.length > 0) {
  //   data.push({ name: selectedAppName || "All", points: deploymentFrequency });
  // }

  useEffect(() => {
    dispatch(fetchApplications());
    dispatch(fetchApplicationCount());
  }, [dispatch]);

  const updateURL = useCallback(
    (kind: ApplicationKind) => {
      history.replace(`${PAGE_PATH_APPLICATIONS}?kind=${kind}`);
    },
    [history]
  );

  const handleApplicationCountClick = useCallback(
    (kind: ApplicationKind) => {
      updateURL(kind);
    },
    [updateURL]
  );

  return (
    <Box flex={1} p={2} overflow="auto">
      <Box display="flex" flexDirection="column" flex={1} p={2}>
        <ApplicationCounts onClick={handleApplicationCountClick} />
      </Box>
      {/* <InsightHeader /> */}
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
