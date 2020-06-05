import React, { FC } from "react";
import { makeStyles, Typography, Box } from "@material-ui/core";

const useStyles = makeStyles(theme => ({
  text: {
    marginLeft: theme.spacing(1)
  }
}));

interface Props {
  label: string;
  text: string;
}

export const LabeledText: FC<Props> = ({ label, text }) => {
  const classes = useStyles();
  return (
    <Box display="flex">
      <Typography variant="subtitle2" color="textSecondary">
        {label}:
      </Typography>
      <Typography variant="body2" className={classes.text}>
        {text}
      </Typography>
    </Box>
  );
};
