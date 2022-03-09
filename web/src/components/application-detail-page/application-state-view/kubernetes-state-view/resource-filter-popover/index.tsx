import {
  Box,
  Button,
  Checkbox,
  FormControlLabel,
  Popover,
} from "@material-ui/core";
import FilterListIcon from "@material-ui/icons/FilterList";
import { FC, useRef, useState } from "react";
import { UI_TEXT_FILTER, UI_TEXT_FILTERED } from "~/constants/ui-text";

export interface ResourceFilterPopoverProps {
  enables: Record<string, boolean>;
  onChange: (state: Record<string, boolean>) => void;
}

export const ResourceFilterPopover: FC<ResourceFilterPopoverProps> = ({
  enables,
  onChange,
}) => {
  const buttonRef = useRef<HTMLButtonElement | null>(null);
  const [state, setState] = useState<Record<string, boolean>>(enables);
  const [open, setOpen] = useState(false);

  const handleApply = (): void => {
    setOpen(false);
    onChange(state);
  };

  const handleClose = (): void => {
    setOpen(false);
  };

  const isFiltered = Object.keys(enables).some((key) => enables[key] === false);

  return (
    <>
      <Box p={2}>
        <Button
          ref={buttonRef}
          startIcon={<FilterListIcon />}
          color={isFiltered ? "primary" : "default"}
          onClick={() => setOpen(!open)}
        >
          {isFiltered ? UI_TEXT_FILTERED : UI_TEXT_FILTER}
        </Button>
      </Box>
      <Popover open={open} anchorEl={buttonRef.current} onClose={handleClose}>
        <Box p={2} minWidth={250}>
          {Object.keys(state).map((key) => (
            <Box key={key} display="flex" alignItems="center">
              <FormControlLabel
                control={
                  <Checkbox
                    color="primary"
                    checked={state[key]}
                    onChange={() => {
                      setState({ ...state, [key]: !state[key] });
                    }}
                  />
                }
                label={key}
              />
            </Box>
          ))}
          <Box textAlign="right">
            <Button onClick={handleClose}>CANCEL</Button>
            <Button color="primary" onClick={handleApply}>
              APPLY
            </Button>
          </Box>
        </Box>
      </Popover>
    </>
  );
};
