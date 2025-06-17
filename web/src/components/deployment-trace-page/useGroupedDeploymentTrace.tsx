import dayjs from "dayjs";
import { useMemo } from "react";
import { sortDateFunc } from "~/utils/common";
import { ListDeploymentTracesResponse } from "~~/api_client/service_pb";

type DeploymentTracesMapByDate = Record<
  string,
  ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[]
>;

type GroupedDeploymentTrace = {
  dates: string[];
  deploymentTracesMap: DeploymentTracesMapByDate;
};

const useGroupedDeploymentTrace = (
  traceList: ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[]
): GroupedDeploymentTrace => {
  const deploymentTracesMap = useMemo(() => {
    const listMap: DeploymentTracesMapByDate = {};

    traceList.forEach((item) => {
      if (!item.trace?.commitTimestamp) return;

      const dateStr = dayjs(item.trace?.commitTimestamp * 1000).format(
        "YYYY/MM/DD"
      );
      if (!listMap[dateStr]) listMap[dateStr] = [];
      listMap[dateStr].push(item);
    });

    return listMap;
  }, [traceList]);

  const dates = useMemo(
    () =>
      Object.keys(deploymentTracesMap).sort((a, b) =>
        sortDateFunc(a, b, "DESC")
      ),
    [deploymentTracesMap]
  );

  return { dates, deploymentTracesMap };
};

export default useGroupedDeploymentTrace;
