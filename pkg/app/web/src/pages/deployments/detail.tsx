import React, { memo, FC, useEffect } from "react";
import { useParams } from "react-router";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchDeploymentById,
  Deployment,
  selectById,
} from "../../modules/deployments";
import { DeploymentDetail } from "../../components/deployment-detail";
import { AppState } from "../../modules";
import { Pipeline } from "../../components/pipeline";

export const DeploymentDetailPage: FC = memo(() => {
  const dispatch = useDispatch();
  const { deploymentId } = useParams<{ deploymentId: string }>();
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectById(state.deployments, deploymentId)
  );

  useEffect(() => {
    if (deploymentId) {
      dispatch(fetchDeploymentById(deploymentId));
    }
  }, [deploymentId]);

  if (!deployment) {
    return <div>loading</div>;
  }

  return (
    <div>
      <DeploymentDetail
        name={deployment.id}
        env={deployment.envId}
        pipedId={deployment.pipedId}
        description={deployment.description}
        status={deployment.status}
        commit={deployment.trigger.commit}
      />
      <Pipeline deploymentId={deploymentId} />
    </div>
  );
});
