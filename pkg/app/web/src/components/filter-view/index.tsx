import React, { FC } from "react";
import { Box, Button, makeStyles, Paper, Typography } from "@material-ui/core";
import { FILTER_PAPER_WIDTH } from "../../styles/size";

const useStyles = makeStyles((theme) => ({
  filterPaper: {
    width: FILTER_PAPER_WIDTH,
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
          Clear
        </Button>
      </Box>

      {children}
    </Paper>
  );
};
