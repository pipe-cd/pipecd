import React, { memo, FC, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  Deployment,
  selectAll,
  fetchDeployments,
} from "../../modules/deployments";
import { AppState } from "../../modules";
import dayjs from "dayjs";
import { Link, makeStyles, Typography } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../../constants";
import { DeploymentItem } from "../../components/deployment-item";
import { fetchApplications } from "../../modules/applications";

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

const sortComp = (a: string | number, b: string | number) => {
  return dayjs(b).valueOf() - dayjs(a).valueOf();
};

const useGroupedDeployments = () => {
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

export const DeploymentIndexPage: FC = memo(() => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const groupedDeployments = useGroupedDeployments();

  useEffect(() => {
    dispatch(fetchDeployments());
    dispatch(fetchApplications());
  }, []);

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
              <ol className={classes.deployments}>
                {groupedDeployments[date]
                  .sort((a, b) => sortComp(a.createdAt, b.createdAt))
                  .map((deployment) => (
                    <li key={deployment.id}>
                      <Link
                        component={RouterLink}
                        to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
                      >
                        <DeploymentItem id={deployment.id} />
                      </Link>
                    </li>
                  ))}
              </ol>
            </li>
          ))}
      </ol>
    </div>
  );
});
