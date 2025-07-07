import {
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  Toolbar,
  Typography,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import FilterIcon from "@mui/icons-material/FilterList";
import RefreshIcon from "@mui/icons-material/Refresh";
import dayjs from "dayjs";
import { FC, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useInView } from "react-intersection-observer";
import { useNavigate } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
  UI_TEXT_MORE,
} from "~/constants/ui-text";
import { SpinnerIcon } from "~/styles/button";
import {
  stringifySearchParams,
  useSearchParams,
  arrayFormat,
} from "~/utils/search-params";
import { DeploymentFilter } from "./deployment-filter";
import { DeploymentItem } from "./deployment-item";
import {
  DeploymentFilterOptions,
  useGetDeploymentsInfinite,
} from "~/queries/deployment/use-get-deployments-infinite";
import { Deployment } from "~/types/deployment";

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

const useGroupedDeployments = (
  deployments: Deployment.AsObject[]
): Record<string, Deployment.AsObject[]> => {
  return useMemo(() => {
    const result: Record<string, Deployment.AsObject[]> = {};

    deployments.forEach((deployment) => {
      const dateStr = dayjs(deployment.createdAt * 1000).format("YYYY/MM/DD");
      if (!result[dateStr]) {
        result[dateStr] = [];
      }
      result[dateStr].push(deployment);
    });

    return result;
  }, [deployments]);
};

export const DeploymentIndexPage: FC = () => {
  const navigate = useNavigate();
  const listRef = useRef(null);
  const filterOptions = useSearchParams();
  const [openFilter, setOpenFilter] = useState(true);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });

  const {
    data: deploymentsData,
    isFetching,
    fetchNextPage: fetchMoreDeployments,
    refetch: refreshDeployments,
    isSuccess,
  } = useGetDeploymentsInfinite(filterOptions);

  const deploymentsList = useMemo(() => {
    return deploymentsData?.pages.flatMap((item) => item.deploymentsList) || [];
  }, [deploymentsData]);

  const hasMore = useMemo(() => {
    if (!deploymentsData || deploymentsData.pages.length === 0) {
      return false;
    }
    const lastIndex = deploymentsData?.pages.length - 1;
    return deploymentsData?.pages?.[lastIndex]?.hasMore || false;
  }, [deploymentsData]);

  const groupedDeployments = useGroupedDeployments(deploymentsList || []);

  useEffect(() => {
    if (inView && hasMore && isFetching === false) {
      fetchMoreDeployments();
    }
  }, [inView, isFetching, filterOptions, fetchMoreDeployments, hasMore]);

  // filter handlers
  const handleFilterChange = useCallback(
    (options: DeploymentFilterOptions) => {
      navigate(
        `${PAGE_PATH_DEPLOYMENTS}?${stringifySearchParams(
          { ...options },
          { arrayFormat: arrayFormat }
        )}`,
        { replace: true }
      );
    },
    [navigate]
  );
  const handleFilterClear = useCallback(() => {
    navigate(PAGE_PATH_DEPLOYMENTS, { replace: true });
  }, [navigate]);

  const handleRefreshClick = useCallback(() => {
    refreshDeployments();
  }, [refreshDeployments]);

  const handleMoreClick = useCallback(() => {
    fetchMoreDeployments();
  }, [fetchMoreDeployments]);

  const dates = Object.keys(groupedDeployments).sort(sortComp);

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
          sx={(theme) => ({
            listStyle: "none",
            padding: theme.spacing(3),
            paddingTop: 0,
            margin: 0,
            flex: 1,
            overflowY: "scroll",
          })}
          ref={listRef}
        >
          {dates.length === 0 &&
            (isFetching ? (
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  mt: 3,
                }}
              >
                <CircularProgress />
              </Box>
            ) : (
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "center",
                  mt: 3,
                }}
              >
                <Typography>No deployments</Typography>
              </Box>
            ))}
          {dates.map((date) => (
            <li key={date}>
              <Typography variant="subtitle1" sx={{ mt: 2, mb: 2 }}>
                {date}
              </Typography>
              <List>
                {groupedDeployments[date]
                  .sort((a, b) => sortComp(a.createdAt, b.createdAt))
                  .map((deployment) => (
                    <DeploymentItem
                      key={deployment.id}
                      deployment={deployment}
                    />
                  ))}
              </List>
            </li>
          ))}
          {isSuccess && <div ref={ref} />}
          {!deploymentsData?.pages?.[deploymentsData.pages.length - 1]
            ?.hasMore && (
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
          {/* TODO: Show how many days have been read */}
        </Box>
        {openFilter && (
          <DeploymentFilter
            options={filterOptions}
            onChange={handleFilterChange}
            onClear={handleFilterClear}
          />
        )}
      </Box>
    </Box>
  );
};
