import React, { memo, FC, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  Deployment,
  selectAll,
  fetchDeployments,
} from "../../modules/deployments";
import { AppState } from "../../modules";
import dayjs from "dayjs";
import { Link } from "@material-ui/core";
import { Link as RouterLink } from "react-router-dom";
import { PAGE_PATH_DEPLOYMENTS } from "../../constants";

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
  const dispatch = useDispatch();
  const groupedDeployments = useGroupedDeployments();

  useEffect(() => {
    dispatch(fetchDeployments());
  }, []);

  console.log(groupedDeployments);

  return (
    <div>
      <ol>
        {Object.keys(groupedDeployments).map((date) => (
          <li key={date}>
            {date}
            <ol>
              {groupedDeployments[date].map((deployment) => (
                <li key={deployment.id}>
                  <Link
                    component={RouterLink}
                    to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
                  >
                    {deployment.id}
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
