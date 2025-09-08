import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Toolbar,
  Typography,
} from "@mui/material";
import { FC, useCallback, useEffect, useMemo, useRef, useState } from "react";
import CloseIcon from "@mui/icons-material/Close";
import FilterIcon from "@mui/icons-material/FilterList";
import RefreshIcon from "@mui/icons-material/Refresh";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_MORE,
  UI_TEXT_REFRESH,
} from "~/constants/ui-text";
import { SpinnerIcon } from "~/styles/button";
import DeploymentTraceFilter from "./deployment-trace-filter";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENT_TRACE } from "~/constants/path";
import {
  arrayFormat,
  stringifySearchParams,
  useSearchParams,
} from "~/utils/search-params";
import useGroupedDeploymentTrace from "./useGroupedDeploymentTrace";
import DeploymentTraceItem from "./deployment-trace-item";
import { useInView } from "react-intersection-observer";
import { useGetDeploymentTracesInfinite } from "~/queries/deployment-traces/use-deployment-traces-infinite";

const DeploymentTracePage: FC = () => {
  const [openFilter, setOpenFilter] = useState(true);
  const navigate = useNavigate();
  const filterValues = useSearchParams();

  const listRef = useRef(null);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });

  const {
    data: deploymentTracesData,
    isFetching,
    fetchNextPage: fetchMoreDeploymentTraces,
    refetch: refetchDeploymentTraces,
    isSuccess,
  } = useGetDeploymentTracesInfinite(filterValues);

  const deploymentTracesList = useMemo(() => {
    return deploymentTracesData?.pages.flatMap((item) => item.tracesList) || [];
  }, [deploymentTracesData]);

  const hasMore = useMemo(() => {
    if (!deploymentTracesData || deploymentTracesData.pages.length === 0) {
      return false;
    }
    const lastIndex = deploymentTracesData?.pages.length - 1;
    return deploymentTracesData?.pages?.[lastIndex]?.hasMore || false;
  }, [deploymentTracesData]);

  const { dates, deploymentTracesMap } = useGroupedDeploymentTrace(
    deploymentTracesList || []
  );

  useEffect(() => {
    if (inView && hasMore && isFetching === false) {
      fetchMoreDeploymentTraces;
    }
  }, [fetchMoreDeploymentTraces, hasMore, inView, isFetching]);

  const handleRefreshClick = (): void => {
    refetchDeploymentTraces();
  };

  const handleMoreClick = useCallback(() => {
    // dispatch(fetchMoreDeploymentTraces(filterValues || {}));
    fetchMoreDeploymentTraces();
  }, [fetchMoreDeploymentTraces]);

  const handleFilterChange = (options: { commitHash?: string }): void => {
    navigate(
      `${PAGE_PATH_DEPLOYMENT_TRACE}?${stringifySearchParams(
        { ...options },
        { arrayFormat: arrayFormat }
      )}`,
      { replace: true }
    );
  };

  const handleFilterClear = useCallback(() => {
    navigate(PAGE_PATH_DEPLOYMENT_TRACE, { replace: true });
  }, [navigate]);

  return (
    <Box
      sx={{
        display: "flex",
        overflow: "hidden",
        flex: 1,
        flexDirection: "column",
      }}
    >
      <Toolbar variant="dense">
        <Box
          sx={{
            flexGrow: 1,
          }}
        />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={handleRefreshClick}
          disabled={isFetching}
        >
          {UI_TEXT_REFRESH}
          {isFetching && <SpinnerIcon />}
        </Button>
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />
      <Box
        sx={{
          display: "flex",
          overflow: "hidden",
          flex: 1,
        }}
      >
        <Box
          component={"ol"}
          sx={{
            listStyle: "none",
            padding: 3,
            paddingTop: 0,
            margin: 0,
            flex: 1,
            overflowY: "scroll",
          }}
          ref={listRef}
        >
          {dates.length === 0 && isFetching && (
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                mt: 3,
              }}
            >
              <CircularProgress />
            </Box>
          )}
          {dates.length === 0 && !isFetching && (
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                mt: 3,
              }}
            >
              <Typography>No deployments</Typography>
            </Box>
          )}
          {dates.map((date) => (
            <Box
              key={date}
              sx={{
                mb: 1,
              }}
            >
              <Typography variant="subtitle1" sx={{ mt: 2, mb: 2 }}>
                {date}
              </Typography>
              <Box
                sx={{
                  bgcolor: "background.paper",
                }}
              >
                {deploymentTracesMap[date].map(({ trace, deploymentsList }) => (
                  <DeploymentTraceItem
                    key={trace?.id}
                    trace={trace}
                    deploymentList={deploymentsList}
                  />
                ))}
              </Box>
            </Box>
          ))}
          {isSuccess && <div ref={ref} />}
          {!hasMore && (
            <Button
              color="primary"
              variant="outlined"
              size="large"
              fullWidth
              onClick={handleMoreClick}
              disabled={isFetching}
            >
              {UI_TEXT_MORE}
              {isFetching && <SpinnerIcon />}
            </Button>
          )}
        </Box>

        {openFilter && (
          <DeploymentTraceFilter
            filterValues={filterValues}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </Box>
    </Box>
  );
};

export default DeploymentTracePage;
