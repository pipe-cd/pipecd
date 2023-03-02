import {
  TextField,
  Select,
  InputLabel,
  FormControl,
  MenuItem,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FC, useState } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAll as selectAllApplications } from "~/modules/applications";
import { selectAllPipeds } from "~/modules/pipeds";
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

export const PipedAutocomplete: FC<Props> = ({ value, onChange }) => {
  const ps = useAppSelector((state) => selectAllPipeds(state));
  const pipeds = ps
    .filter((piped) => !piped.disabled)
    .sort((a, b) => sortComp(a.name, b.name));
  const [selectedPipedId, setSelectedPipedId] = useState(
    pipeds.length === 1 ? pipeds[0].id : ""
  );
  // const selectedPiped = useAppSelector(selectPipedById(selectedPipedId));
  // const handleUpdateFilterValue = (
  //   optionPart: Partial<PlatformProviderFilterOptions>
  // ): void => {
  //   onChange({ ...options, ...optionPart });
  // };
  // const selectedPiped = useAppSelector(selectPipedById(values.pipedId));
  // const piped = useAppSelector<string[]>(
  //   (state) =>
  //     sortedSet(
  //       selectAllApplications(state.applications).map((app) => app.name)
  //     ),
  //   (left, right) => JSON.stringify(left) === JSON.stringify(right)
  // );
  return (
    <FormControl className="piped" variant="outlined">
      <InputLabel id="filter-piped">Piped</InputLabel>
      <Select
        labelId="filter-piped"
        id="filter-piped"
        label="Piped"
        value={selectedPipedId}
        // className={classes.select}
        onChange={(e) => {
          setSelectedPipedId(e.target.value as string);
          onChange((e.target.value as string) || "");
        }}
      >
        {pipeds.map((e) => (
          <MenuItem value={e.id} key={`piped-${e.id}`}>
            {e.name} ({e.id})
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

// function FormSelectInput<T extends { name: string; value: string }>({
//   id,
//   label,
//   value,
//   items,
//   required = true,
//   onChange,
//   disabled = false,
// }: {
//   id: string;
//   label: string;
//   value: string;
//   items: T[];
//   required?: boolean;
//   onChange: (value: T) => void;
//   disabled?: boolean;
// }): ReactElement {
//   return (
//     <TextField
//       id={id}
//       name={id}
//       label={label}
//       fullWidth
//       required={required}
//       select
//       disabled={disabled}
//       variant="outlined"
//       margin="dense"
//       onChange={(e) => {
//         const nextItem = items.find((item) => item.value === e.target.value);
//         if (nextItem) {
//           onChange(nextItem);
//         }
//       }}
//       value={value}
//       style={{ flex: 1 }}
//     >
//       {items.map((item) => (
//         <MenuItem key={item.name} value={item.value}>
//           {item.name}
//         </MenuItem>
//       ))}
//     </TextField>
//   );
// }

const sortComp = (a: string, b: string): number => {
  return a > b ? 1 : -1;
};
