import { Select, MenuItem } from "@material-ui/core";
import { FC, useState } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAllPipeds } from "~/modules/pipeds";

interface Props {
  onChange: (value: string) => void;
}
const ALL_VALUE = "ALL";

export const PipedSelect: FC<Props> = ({ onChange }) => {
  const ps = useAppSelector((state) => selectAllPipeds(state));
  const pipeds = ps
    .filter((piped) => !piped.disabled)
    .sort((a, b) => sortComp(a.name, b.name));
  const [selectedPipedId, setSelectedPipedId] = useState(
    pipeds.length === 1 ? pipeds[0].id : ""
  );
  return (
      <Select
        labelId="filter-piped"
        id="filter-piped"
        label="Piped"
        value={selectedPipedId ?? ALL_VALUE}
        onChange={(e) => {
          setSelectedPipedId(e.target.value as string);
          onChange(
            e.target.value === ALL_VALUE ? "" : (e.target.value as string)
          );
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
  );
};

const sortComp = (a: string, b: string): number => {
  return a > b ? 1 : -1;
};
