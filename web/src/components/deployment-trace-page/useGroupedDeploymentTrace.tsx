import dayjs from "dayjs";
import { useMemo } from "react";
import { useShallowEqualSelector } from "~/hooks/redux";
import { selectIds, selectById } from "~/modules/deploymentTrace";
import { sortDateFunc } from "~/utils/common";
import { ListDeploymentTracesResponse } from "~~/api_client/service_pb";

type GroupedDeploymentTrace = {
  dates: string[];
  deploymentTracesMap: Record<
    string,
    ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[]
  >;
};

const useGroupedDeploymentTrace = (): GroupedDeploymentTrace => {
  const traceList = useShallowEqualSelector((state) => {
    const list = selectIds(state.deploymentTrace)
      .map((id) => selectById(state.deploymentTrace, id))
      .filter((trace) => trace !== undefined);
    return list;
  }) as ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[];

  const deploymentTracesMap = useMemo(() => {
    const listMap: Record<
      string,
      ListDeploymentTracesResponse.DeploymentTraceRes.AsObject[]
    > = {};

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
    () => Object.keys(deploymentTracesMap).sort(sortDateFunc),
    [deploymentTracesMap]
  );

  return { dates, deploymentTracesMap };
};

export default useGroupedDeploymentTrace;
