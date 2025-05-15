import {
  Button,
  ButtonGroup,
  ButtonGroupOwnProps,
  ClickAwayListener,
  Grow,
  MenuItem,
  MenuList,
  Paper,
  Popper,
} from "@mui/material";
import ArrowDropDownIcon from "@mui/icons-material/ArrowDropDown";
import { FC, useRef, useState } from "react";
import * as React from "react";
import { SpinnerIcon } from "~/styles/button";

export interface SplitButtonProps {
  options: string[];
  label: string;
  onClick: (index: number) => void;
  startIcon?: React.ReactNode;
  disabled: boolean;
  loading: boolean;
  color?: ButtonGroupOwnProps["color"];
  className?: string;
}

export const SplitButton: FC<SplitButtonProps> = ({
  onClick,
  options,
  disabled,
  loading,
  startIcon,
  className,
  color,
  label,
}) => {
  const anchorRef = useRef<HTMLDivElement | null>(null);
  const [open, setOpen] = useState(false);
  const [selectedCancelOption, setSelectedCancelOption] = useState(0);

  const handleClose = (event: MouseEvent | TouchEvent): void => {
    if (anchorRef.current && anchorRef.current.contains(event.target as Node)) {
      return;
    }

    setOpen(false);
  };

  return (
    <div className={className}>
      <ButtonGroup
        variant="outlined"
        color={color || "inherit"}
        ref={anchorRef}
        disabled={disabled}
      >
        <Button
          startIcon={startIcon}
          onClick={() => onClick(selectedCancelOption)}
          disabled={disabled}
        >
          {options[selectedCancelOption]}
          {loading && <SpinnerIcon />}
        </Button>
        <Button
          size="small"
          aria-controls={open ? "split-button-menu" : undefined}
          aria-expanded={open ? "true" : undefined}
          aria-label={label}
          aria-haspopup="menu"
          onClick={() => setOpen(!open)}
        >
          <ArrowDropDownIcon />
        </Button>
      </ButtonGroup>
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        role={undefined}
        transition
        disablePortal
      >
        {({ TransitionProps, placement }) => (
          <Grow
            {...TransitionProps}
            style={{
              transformOrigin:
                placement === "bottom" ? "center top" : "center bottom",
            }}
          >
            <Paper>
              <ClickAwayListener onClickAway={handleClose}>
                <MenuList id="split-button-menu">
                  {options.map((option, index) => (
                    <MenuItem
                      key={option}
                      selected={index === selectedCancelOption}
                      onClick={() => {
                        setSelectedCancelOption(index);
                        setOpen(false);
                      }}
                    >
                      {option}
                    </MenuItem>
                  ))}
                </MenuList>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </div>
  );
};
