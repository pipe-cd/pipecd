import {
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  ListItem,
  makeStyles,
  Toolbar,
  Typography,
} from "@material-ui/core";
import { FC, useCallback, useEffect, useRef, useState } from "react";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import RefreshIcon from "@material-ui/icons/Refresh";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_MORE,
  UI_TEXT_REFRESH,
} from "~/constants/ui-text";
import { useStyles as useButtonStyles } from "~/styles/button";
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

const useStyles = makeStyles((theme) => ({
  list: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
    flex: 1,
    overflowY: "scroll",
  },
  listItem: {
    backgroundColor: theme.palette.background.paper,
  },
  listDeployment: {
    backgroundColor: theme.palette.background.paper,
  },
  date: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
}));
const DeploymentTracePage: FC = () => {
  const classes = useStyles();
  const [openFilter, setOpenFilter] = useState(true);
  const dispatch = useAppDispatch();
  const buttonClasses = useButtonStyles();
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
          {isLoading && (
            <CircularProgress size={24} className={buttonClasses.progress} />
          )}
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
        <ol className={classes.list} ref={listRef}>
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
            <Box key={date}>
              <Typography variant="subtitle1" className={classes.date}>
                {date}
              </Typography>
              <List className={classes.listDeployment}>
                {deploymentTracesMap[date].map(({ trace, deploymentsList }) => (
                  <ListItem
                    key={trace?.id}
                    dense
                    divider
                    color="primary"
                    className={classes.listItem}
                  >
                    <DeploymentTraceItem
                      key={trace?.id}
                      trace={trace}
                      deploymentList={deploymentsList}
                    />
                  </ListItem>
                ))}
              </List>
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
              {isLoading && (
                <CircularProgress
                  size={24}
                  className={buttonClasses.progress}
                />
              )}
            </Button>
          )}
        </ol>

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
