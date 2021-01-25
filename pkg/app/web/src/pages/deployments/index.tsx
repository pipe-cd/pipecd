import {
  Box,
  Button,
  CircularProgress,
  Divider,
  List,
  makeStyles,
  Toolbar,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import FilterIcon from "@material-ui/icons/FilterList";
import RefreshIcon from "@material-ui/icons/Refresh";
import dayjs from "dayjs";
import React, {
  FC,
  memo,
  useCallback,
  useEffect,
  useState,
  useRef,
} from "react";
import { useDispatch, useSelector } from "react-redux";
import { DeploymentFilter } from "../../components/deployment-filter";
import { DeploymentItem } from "../../components/deployment-item";
import { AppState } from "../../modules";
import { fetchApplications } from "../../modules/applications";
import {
  Deployment,
  fetchDeployments,
  selectIds as selectDeploymentIds,
  selectById as selectDeploymentById,
  fetchMoreDeployments,
} from "../../modules/deployments";
import { useInView } from "react-intersection-observer";
import { LoadingStatus } from "../../types/module";
import {
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_REFRESH,
} from "../../constants/ui-text";
import { useStyles as useButtonStyles } from "../../styles/button";

const useStyles = makeStyles((theme) => ({
  deploymentLists: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
    flex: 1,
    overflowY: "scroll",
  },
  date: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
}));

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

function filterUndefined<TValue>(value: TValue | undefined): value is TValue {
  return value !== undefined;
}

const useGroupedDeployments = (): [
  LoadingStatus,
  boolean,
  Record<string, Deployment.AsObject[]>
] => {
  const [status, hasMore, deployments] = useSelector<
    AppState,
    [LoadingStatus, boolean, Deployment.AsObject[]]
  >((state) => [
    state.deployments.status,
    state.deployments.hasMore,
    selectDeploymentIds(state.deployments)
      .map((id) => selectDeploymentById(state.deployments, id))
      .filter(filterUndefined),
  ]);

  const result: Record<string, Deployment.AsObject[]> = {};

  deployments.forEach((deployment) => {
    const dateStr = dayjs(deployment.createdAt * 1000).format("YYYY/MM/DD");
    if (!result[dateStr]) {
      result[dateStr] = [];
    }
    result[dateStr].push(deployment);
  });

  return [status, hasMore, result];
};

export const DeploymentIndexPage: FC = memo(function DeploymentIndexPage() {
  const classes = useStyles();
  const buttonClasses = useButtonStyles();
  const dispatch = useDispatch();
  const listRef = useRef(null);
  const [status, hasMore, groupedDeployments] = useGroupedDeployments();
  const [openFilter, setOpenFilter] = useState(false);
  const [ref, inView] = useInView({
    rootMargin: "400px",
    root: listRef.current,
  });
  const isLoading = status === "loading";

  useEffect(() => {
    dispatch(fetchApplications());
    dispatch(fetchDeployments());
  }, [dispatch]);

  useEffect(() => {
    if (inView && hasMore && isLoading) {
      dispatch(fetchMoreDeployments());
    }
  }, [dispatch, inView, hasMore, isLoading]);

  const handleFilterChange = useCallback(() => {
    dispatch(fetchDeployments());
  }, [dispatch]);

  const handleRefreshClick = useCallback(() => {
    dispatch(fetchDeployments());
  }, [dispatch]);

  const dates = Object.keys(groupedDeployments).sort(sortComp);

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
        <ol className={classes.deploymentLists} ref={listRef}>
          {dates.length === 0 &&
            (isLoading ? (
              <Box display="flex" justifyContent="center" mt={3}>
                <CircularProgress />
              </Box>
            ) : (
              <Box display="flex" justifyContent="center" mt={3}>
                <Typography>No deployments</Typography>
              </Box>
            ))}
          {dates.map((date) => (
            <li key={date}>
              <Typography variant="subtitle1" className={classes.date}>
                {date}
              </Typography>
              <List>
                {groupedDeployments[date]
                  .sort((a, b) => sortComp(a.createdAt, b.createdAt))
                  .map((deployment) => (
                    <DeploymentItem
                      id={deployment.id}
                      key={`deployment-item-${deployment.id}`}
                    />
                  ))}
              </List>
            </li>
          ))}
          {status === "succeeded" && <div ref={ref} />}
        </ol>
        {openFilter && <DeploymentFilter onChange={handleFilterChange} />}
      </Box>
    </Box>
  );
});
