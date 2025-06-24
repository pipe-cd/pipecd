import { useEffect, useRef, useState } from "react";
import * as echarts from "echarts/core";

type Props = {
  extensions: Parameters<typeof echarts.use>[0];
};

type EChartState = ({
  extensions,
}: Props) => {
  chartElm: React.MutableRefObject<HTMLDivElement | null>;
  chart: echarts.ECharts | null;
};

const useEChartState: EChartState = ({ extensions }) => {
  const chartElm = useRef<HTMLDivElement | null>(null);
  const [chart, setChart] = useState<echarts.ECharts | null>(null);
  const [isExtensionsAdded, setIsExtensionsAdded] = useState(false);

  useEffect(() => {
    echarts.use(extensions);
    setIsExtensionsAdded(true);
    // Only trigger use one, should not add extensions to dependencies to retrigger the effect
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (chartElm.current && isExtensionsAdded) {
      setChart(echarts.init(chartElm.current));
    }
  }, [chartElm, isExtensionsAdded]);

  return {
    chartElm,
    chart,
  };
};

export default useEChartState;
