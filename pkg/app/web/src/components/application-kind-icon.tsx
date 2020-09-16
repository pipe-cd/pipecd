import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";
import { ApplicationKind } from "../modules/applications";
import kubernetesIcon from "../../assets/kubernetes.svg";

const useStyles = makeStyles((theme) => ({
  main: {
    width: 32,
    height: 32,
    alignSelf: "center",
    marginRight: theme.spacing(1),
  },
}));

interface Props {
  kind: ApplicationKind;
}

export const ApplicationKindIcon: FC<Props> = ({ kind }) => {
  const classes = useStyles();

  switch (kind) {
    case ApplicationKind.KUBERNETES:
      return <img src={kubernetesIcon} className={classes.main} />;
    default:
      return null;
  }
};
