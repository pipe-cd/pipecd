import React, { FC } from "react";
import { makeStyles, Typography, Box } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  value: {
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  label: string;
  value: React.ReactChild;
}

export const LabeledText: FC<Props> = ({ label, value }) => {
  const classes = useStyles();
  return (
    <Box display="flex">
      <Typography variant="subtitle2" color="textSecondary">
        {label}:
      </Typography>
      <Typography variant="body2" className={classes.value}>
        {value}
      </Typography>
    </Box>
  );
};
