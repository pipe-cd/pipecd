import { FC } from "react";
import { Typography } from "@mui/material";

import makeStyles from "@mui/styles/makeStyles";

const useStyles = makeStyles((theme) => ({
  label: {
    color: theme.palette.text.secondary,
    marginRight: theme.spacing(2),
    minWidth: 120,
  },
  root: {
    display: "flex",
    alignItems: "center",
  },
}));

export interface ProjectSettingLabeledTextProps {
  label: string;
  value: string;
}

export const ProjectSettingLabeledText: FC<ProjectSettingLabeledTextProps> = ({
  label,
  value,
}) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Typography variant="subtitle1" className={classes.label}>
        {label}
      </Typography>
      <Typography variant="body1">{value}</Typography>
    </div>
  );
};
