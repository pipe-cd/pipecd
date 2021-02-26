import { Box } from "@material-ui/core";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { DeploymentFrequencyChart } from "../../components/deployment-frequency-chart";
import { InsightHeader } from "../../components/insight-header";
import { AppState } from "../../modules";
import { fetchApplications, selectById } from "../../modules/applications";
import { InsightDataPoint } from "../../modules/insight";

export const InsightIndexPage: FC = memo(function InsightIndexPage() {
  const dispatch = useDispatch();

  const deploymentFrequency = useSelector<
    AppState,
    InsightDataPoint.AsObject[]
  >((state) => state.deploymentFrequency.data);
  const selectedAppName = useSelector<AppState, string | undefined>((state) =>
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
    <Box flex={1} p={2}>
      <InsightHeader />
      <Box
        display="grid"
        gridGap="24px"
        gridTemplateColumns="repeat(2, 1fr)"
        mt={2}
      >
        <DeploymentFrequencyChart data={data} />
      </Box>
    </Box>
  );
});
