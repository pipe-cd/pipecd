import { TextField } from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FC } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAll as selectAllApplications } from "~/modules/applications";
import { sortedSet } from "~/utils/sorted-set";

interface Props {
  value: string | null;
  onChange: (value: string) => void;
}

export const ApplicationAutocomplete: FC<Props> = ({ value, onChange }) => {
  const applications = useAppSelector<string[]>(
    (state) =>
      sortedSet(
        selectAllApplications(state.applications).map((app) => app.name)
      ),
    (left, right) => JSON.stringify(left) === JSON.stringify(right)
  );
  return (
    <Autocomplete
      id="name"
      options={applications}
      value={value}
      onChange={(_, value) => {
        onChange(value || "");
      }}
      renderInput={(params) => (
        <TextField {...params} label="Application Name" variant="outlined" />
      )}
    />
  );
};
