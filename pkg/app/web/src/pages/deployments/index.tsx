import { List, makeStyles, Typography } from "@material-ui/core";
import dayjs from "dayjs";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { DeploymentItem } from "../../components/deployment-item";
import { AppState } from "../../modules";
import { fetchApplications } from "../../modules/applications";
import {
  Deployment,
  fetchDeployments,
  selectAll,
} from "../../modules/deployments";

const useStyles = makeStyles((theme) => ({
  deploymentLists: {
    listStyle: "none",
    padding: theme.spacing(3),
    paddingTop: 0,
    margin: 0,
  },
  deployments: {
    listStyle: "none",
    padding: 0,
  },

  date: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2),
  },
}));

const sortComp = (a: string | number, b: string | number): number => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

const useGroupedDeployments = (): Record<string, Deployment[]> => {
  const deployments = useSelector<AppState, Deployment[]>((state) =>
    selectAll(state.deployments)
  );

  const result: Record<string, Deployment[]> = {};

  deployments.forEach((deployment) => {
    const dateStr = dayjs(deployment.createdAt * 1000).format("YYYY/MM/DD");
    if (!result[dateStr]) {
      result[dateStr] = [];
    }
    result[dateStr].push(deployment);
  });

  return result;
};

export const DeploymentIndexPage: FC = memo(function DeploymentIndexPage() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const groupedDeployments = useGroupedDeployments();

  useEffect(() => {
    dispatch(fetchDeployments());
    dispatch(fetchApplications());
  }, [dispatch]);

  return (
    <div>
      <ol className={classes.deploymentLists}>
        {Object.keys(groupedDeployments)
          .sort(sortComp)
          .map((date) => (
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
      </ol>
    </div>
  );
});
