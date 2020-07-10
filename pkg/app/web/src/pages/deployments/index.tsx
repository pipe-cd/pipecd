import {
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
import dayjs from "dayjs";
import React, { FC, memo, useCallback, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { DeploymentFilter } from "../../components/deployment-filter";
import { DeploymentItem } from "../../components/deployment-item";
import { AppState } from "../../modules";
import { fetchApplications } from "../../modules/applications";
import {
  Deployment,
  fetchDeployments,
  selectAll,
} from "../../modules/deployments";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    overflow: "hidden",
    flex: 1,
    flexDirection: "column",
  },
  main: {
    display: "flex",
    overflow: "hidden",
    flex: 1,
  },
  toolbarSpacer: {
    flexGrow: 1,
  },
  deploymentLists: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
    flex: 1,
    overflowY: "scroll",
  },
  loadingContainer: {
    display: "flex",
    justifyContent: "center",
    marginTop: theme.spacing(3),
  },
  deployments: {
    listStyle: "none",
    padding: 0,
  },
  date: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
  empty: {
    display: "flex",
    justifyContent: "center",
    marginTop: theme.spacing(3),
  },
}));

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

const useGroupedDeployments = (): [boolean, Record<string, Deployment[]>] => {
  const [isLoading, deployments] = useSelector<
    AppState,
    [boolean, Deployment[]]
  >((state) => [state.deployments.loadingList, selectAll(state.deployments)]);

  const result: Record<string, Deployment[]> = {};

  deployments.forEach((deployment) => {
    const dateStr = dayjs(deployment.createdAt * 1000).format("YYYY/MM/DD");
    if (!result[dateStr]) {
      result[dateStr] = [];
    }
    result[dateStr].push(deployment);
  });

  return [isLoading, result];
};

export const DeploymentIndexPage: FC = memo(function DeploymentIndexPage() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [isLoading, groupedDeployments] = useGroupedDeployments();
  const [isOpenFilter, setIsOpenFilter] = useState(false);

  useEffect(() => {
    dispatch(fetchApplications());
  }, [dispatch]);

  const handleChangeFilter = useCallback(
    (options) => {
      dispatch(fetchDeployments(options));
    },
    [dispatch]
  );

  const dates = Object.keys(groupedDeployments).sort(sortComp);

  return (
    <div className={classes.root}>
      <Toolbar variant="dense">
        <div className={classes.toolbarSpacer} />
        <Button
          color="primary"
          startIcon={isOpenFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setIsOpenFilter(!isOpenFilter)}
        >
          {isOpenFilter ? "HIDE FILTER" : "FILTER"}
        </Button>
      </Toolbar>

      <Divider />
      <div className={classes.main}>
        <ol className={classes.deploymentLists}>
          {isLoading ? (
            <div className={classes.loadingContainer}>
              <CircularProgress />
            </div>
          ) : dates.length === 0 ? (
            <div className={classes.empty}>
              <Typography>No deployments</Typography>
            </div>
          ) : (
            dates.map((date) => (
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
            ))
          )}
        </ol>
        <DeploymentFilter open={isOpenFilter} onChange={handleChangeFilter} />
      </div>
    </div>
  );
});
