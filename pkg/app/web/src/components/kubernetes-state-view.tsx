import { makeStyles } from "@material-ui/core";
import React, { FC } from "react";
import { KubernetesResourceState } from "../modules/applications-live-state";
import { KubernetesResource } from "./kubernetes-resource";

const useStyles = makeStyles(() => ({
  container: {},
}));

interface Props {
  resources: KubernetesResourceState[];
}

export const KubernetesStateView: FC<Props> = ({ resources }) => {
  const classes = useStyles();

  return (
    <div className={classes.container}>
      {resources.map((resource) => (
        <KubernetesResource name={resource.name} kind={resource.kind} />
      ))}
    </div>
  );
};
