import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC } from "react";
import { FilterView } from "~/components/filter-view";

const useStyles = makeStyles((theme) => ({
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
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

export interface PipedFilterProps {
  values: FilterValues;
  onChange: (values: FilterValues) => void;
}

export const PipedFilter: FC<PipedFilterProps> = ({ values, onChange }) => {
  const classes = useStyles();

  return (
    <FilterView onClear={() => onChange({ enabled: true })}>
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
    </FilterView>
  );
};
