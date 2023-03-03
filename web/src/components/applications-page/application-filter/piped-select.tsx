import {
  Select,
  InputLabel,
  FormControl,
  MenuItem,
} from "@material-ui/core";
import { FC, useState } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAllPipeds } from "~/modules/pipeds";

interface Props {
  value: string | null;
  onChange: (value: string) => void;
  className: string | undefined;
}
const ALL_VALUE = "ALL"

export const PipedSelect: FC<Props> = ({ value, onChange, className }) => {
  const ps = useAppSelector((state) => selectAllPipeds(state));
  const pipeds = ps
    .filter((piped) => !piped.disabled)
    .sort((a, b) => sortComp(a.name, b.name));
  const [selectedPipedId, setSelectedPipedId] = useState(
    pipeds.length === 1 ? pipeds[0].id : ""
  );
  return (
    <FormControl className="piped" variant="outlined">
      <InputLabel id="filter-piped">Piped</InputLabel>
      <Select
        labelId="filter-piped"
        id="filter-piped"
        label="Piped"
        value={selectedPipedId ?? ALL_VALUE}
        className={className}
        onChange={(e) => {
          setSelectedPipedId(e.target.value as string);
          onChange((
            e.target.value === ALL_VALUE
              ? ""
              : e.target.value as string
          ));
        }}
      >
        <MenuItem value={ALL_VALUE}>
          <em>All</em>
        </MenuItem>
        {pipeds.map((e) => (
          <MenuItem value={e.id} key={`piped-${e.id}`}>
            {e.name} ({e.id})
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};

const sortComp = (a: string, b: string): number => {
  return a > b ? 1 : -1;
};