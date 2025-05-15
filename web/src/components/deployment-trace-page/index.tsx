import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Toolbar,
  Typography,
} from "@mui/material";
import { FC, useCallback, useEffect, useRef, useState } from "react";
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
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  fetchDeploymentTraces,
  fetchMoreDeploymentTraces,
} from "~/modules/deploymentTrace";
import useGroupedDeploymentTrace from "./useGroupedDeploymentTrace";
import DeploymentTraceItem from "./deployment-trace-item";
import { useInView } from "react-intersection-observer";

const DeploymentTracePage: FC = () => {
  const [openFilter, setOpenFilter] = useState(true);
  const dispatch = useAppDispatch();
  const status = useAppSelector((state) => state.deploymentTrace.status);
  const hasMore = useAppSelector((state) => state.deploymentTrace.hasMore);
  const navigate = useNavigate();
  const filterValues = useSearchParams();
  const { dates, deploymentTracesMap } = useGroupedDeploymentTrace();
  const isLoading = status === "loading";

  const listRef = useRef(null);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });

  useEffect(() => {
    dispatch(fetchDeploymentTraces(filterValues));
  }, [dispatch, filterValues]);

  useEffect(() => {
    if (inView && hasMore && isLoading === false) {
      dispatch(fetchMoreDeploymentTraces(filterValues || {}));
    }
  }, [dispatch, inView, hasMore, isLoading, filterValues]);

  const handleRefreshClick = (): void => {
    dispatch(fetchDeploymentTraces(filterValues));
  };

  const handleMoreClick = useCallback(() => {
    dispatch(fetchMoreDeploymentTraces(filterValues || {}));
  }, [dispatch, filterValues]);

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
    <Box display="flex" overflow="hidden" flex={1} flexDirection="column">
      <Toolbar variant="dense">
        <Box flexGrow={1} />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          onClick={handleRefreshClick}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
          {isLoading && <SpinnerIcon />}
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
      <Box display="flex" overflow="hidden" flex={1}>
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
          {dates.length === 0 && isLoading && (
            <Box display="flex" justifyContent="center" mt={3}>
              <CircularProgress />
            </Box>
          )}
          {dates.length === 0 && !isLoading && (
            <Box display="flex" justifyContent="center" mt={3}>
              <Typography>No deployments</Typography>
            </Box>
          )}
          {dates.map((date) => (
            <Box key={date} mb={1}>
              <Typography variant="subtitle1" sx={{ mt: 2, mb: 2 }}>
                {date}
              </Typography>
              <Box bgcolor={"background.paper"}>
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
          {status === "succeeded" && <div ref={ref} />}
          {!hasMore && (
            <Button
              color="primary"
              variant="outlined"
              size="large"
              fullWidth
              onClick={handleMoreClick}
              disabled={isLoading}
            >
              {UI_TEXT_MORE}
              {isLoading && <SpinnerIcon />}
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
