import { FC } from "react";
import * as React from "react";
import { Box, Button, Paper, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { UI_TEXT_CLEAR } from "~/constants/ui-text";

const useStyles = makeStyles((theme) => ({
  filterPaper: {
    width: 360,
    minWidth: 280,
    padding: theme.spacing(3),
    height: "100%",
  },
}));

export interface FilterViewProps {
  onClear: () => void;
  children: React.ReactNode;
}

export const FilterView: FC<FilterViewProps> = ({ onClear, children }) => {
  const classes = useStyles();
  return (
    <Paper square className={classes.filterPaper}>
      <Box display="flex" justifyContent="space-between">
        <Typography variant="h6" component="span">
          Filters
        </Typography>
        <Button color="primary" onClick={onClear}>
          {UI_TEXT_CLEAR}
        </Button>
      </Box>

      {children}
    </Paper>
  );
};
