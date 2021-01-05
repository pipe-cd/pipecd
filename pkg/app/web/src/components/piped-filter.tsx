import {
  Button,
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Paper,
  Select,
  Typography,
} from "@material-ui/core";
import React, { FC } from "react";

const FILTER_PAPER_WIDTH = 360;

const useStyles = makeStyles((theme) => ({
  filterPaper: {
    width: FILTER_PAPER_WIDTH,
    padding: theme.spacing(3),
    height: "100%",
  },
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
  header: {
    display: "flex",
    justifyContent: "space-between",
  },
  select: {
    width: "100%",
  },
}));

const ALL_VALUE = "ALL";
const getActiveStatusText = (v: boolean): string =>
  v ? "enabled" : "disabled";
const textValueMap: Record<string, undefined | boolean> = {
  [ALL_VALUE]: undefined,
  enabled: true,
  disabled: false,
};

export interface FilterValues {
  enabled: boolean | undefined;
}

interface Props {
  values: FilterValues;
  onChange: (values: FilterValues) => void;
}

export const PipedFilter: FC<Props> = ({ values, onChange }) => {
  const classes = useStyles();

  return (
    <Paper square className={classes.filterPaper}>
      <div className={classes.header}>
        <Typography variant="h6" component="span">
          Filters
        </Typography>
        <Button
          color="primary"
          onClick={() => {
            onChange({ enabled: true });
          }}
        >
          Clear
        </Button>
      </div>

      <FormControl className={classes.formItem} variant="outlined">
        <InputLabel id="filter-active-status">Active Status</InputLabel>
        <Select
          labelId="filter-active-status"
          id="filter-active-status"
          value={
            values.enabled === undefined
              ? ALL_VALUE
              : getActiveStatusText(values.enabled)
          }
          label="Active Status"
          className={classes.select}
          onChange={(e) => {
            onChange({
              enabled: textValueMap[e.target.value as string],
            });
          }}
        >
          <MenuItem value={ALL_VALUE}>
            <em>All</em>
          </MenuItem>
          <MenuItem value="enabled">Enabled</MenuItem>
          <MenuItem value="disabled">Disabled</MenuItem>
        </Select>
      </FormControl>
    </Paper>
  );
};
