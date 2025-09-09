import { TextField } from "@mui/material";
import Autocomplete from "@mui/material/Autocomplete";
import { FC, useMemo } from "react";
import { useGetApplications } from "~/queries/applications/use-get-applications";
import { sortedSet } from "~/utils/sorted-set";

interface Props {
  value: string | null;
  onChange: (value: string) => void;
}

export const ApplicationAutocomplete: FC<Props> = ({ value, onChange }) => {
  const { data: applications = [] } = useGetApplications();
  const options = useMemo(() => {
    return sortedSet(applications.map((app) => app.name));
  }, [applications]);

  return (
    <Autocomplete
      id="name"
      options={options}
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
