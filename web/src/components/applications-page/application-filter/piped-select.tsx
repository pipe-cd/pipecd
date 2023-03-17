import { makeStyles } from "@material-ui/styles";
import { Select, MenuItem, IconButton } from "@material-ui/core";
import ClearIcon from "@material-ui/icons/Clear";
import { FC, useState } from "react";
import { useAppSelector } from "~/hooks/redux";
import { selectAllPipeds } from "~/modules/pipeds";
import clsx from "clsx";

interface Props {
  defaultPipedId: string;
  onChange: (value: string) => void;
}

const useStyles = makeStyles(() => ({
  root: {
    "&:hover $clearIndicatorDirty, & .Mui-focused $clearIndicatorDirty": {
      visibility: "visible",
    },
  },
  clearIndicatorDirty: {},
  clearIndicator: {
    visibility: "hidden",
    right: 20,
  },
}));

export const PipedSelect: FC<Props> = ({ defaultPipedId, onChange }) => {
  const classes = useStyles();

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
      className={classes.root}
      label="Piped"
      value={selectedPipedId ? selectedPipedId : defaultPipedId}
      onChange={(e) => {
        setSelectedPipedId(e.target.value as string);
        onChange(e.target.value as string);
      }}
      endAdornment={
        <IconButton
          className={clsx(classes.clearIndicator, {
            [classes.clearIndicatorDirty]: setSelectedPipedId.length > 0,
          })}
          size="small"
          onClick={() => {
            setSelectedPipedId("");
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
