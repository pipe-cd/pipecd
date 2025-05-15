import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";
import { FC } from "react";
import { FilterView } from "~/components/filter-view";

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
  return (
    <FilterView onClear={() => onChange({ enabled: true })}>
      <FormControl sx={{ marginTop: 4 }} fullWidth variant="outlined">
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
          fullWidth
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
