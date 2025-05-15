import { Select, MenuItem, IconButton } from "@mui/material";
import ClearIcon from "@mui/icons-material/Clear";
import { FC } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAllPipeds } from "~/modules/pipeds";

interface Props {
  value: string;
  onChange: (value: string) => void;
}

export const PipedSelect: FC<Props> = ({ value, onChange }) => {
  const ps = useAppSelector((state) => selectAllPipeds(state));
  const pipeds = ps
    .filter((piped) => !piped.disabled)
    .sort((a, b) => sortComp(a.name, b.name));

  return (
    <Select
      labelId="filter-piped"
      id="filter-piped"
      label="Piped"
      value={value}
      onChange={(e) => {
        onChange(e.target.value as string);
      }}
      sx={{
        "&:hover .clearIndicator, &:focus-within .clearIndicator": {
          visibility: value && value.length > 0 ? "visible" : "hidden",
        },
      }}
      endAdornment={
        <IconButton
          sx={{
            visibility: "hidden",
            right: 20,
          }}
          size="small"
          className="clearIndicator"
          onClick={() => {
            onChange("");
          }}
        >
          <ClearIcon fontSize="small" />
        </IconButton>
      }
    >
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
