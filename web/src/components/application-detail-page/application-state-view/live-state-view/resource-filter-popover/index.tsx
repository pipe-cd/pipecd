import {
  Box,
  Button,
  Checkbox,
  FormControlLabel,
  Popover,
} from "@material-ui/core";
import FilterListIcon from "@material-ui/icons/FilterList";
import { FC, useEffect, useRef, useState } from "react";
import { UI_TEXT_FILTER, UI_TEXT_FILTERED } from "~/constants/ui-text";

type Props = {
  filterState: Record<string, boolean>;
  onChange: (state: Record<string, boolean>) => void;
};

export const ResourceFilterPopover: FC<Props> = ({ filterState, onChange }) => {
  const buttonRef = useRef<HTMLButtonElement | null>(null);
  const [state, setState] = useState<Record<string, boolean>>(filterState);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    setState(filterState);
  }, [filterState]);

  const handleApply = (): void => {
    setOpen(false);
    onChange(state);
  };

  const handleClose = (): void => {
    setOpen(false);
  };

  const isFiltered = Object.keys(filterState).some(
    (key) => filterState[key] === false
  );

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
          {Object.keys(filterState).map((resourceType) => (
            <Box key={resourceType} display="flex" alignItems="center">
              <FormControlLabel
                control={
                  <Checkbox
                    color="primary"
                    checked={state[resourceType]}
                    onChange={() => {
                      setState({
                        ...state,
                        [resourceType]: !state[resourceType],
                      });
                    }}
                  />
                }
                label={resourceType}
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
