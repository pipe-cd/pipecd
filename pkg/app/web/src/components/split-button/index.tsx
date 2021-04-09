import {
  Button,
  ButtonGroup,
  PropTypes,
  CircularProgress,
  ClickAwayListener,
  Grow,
  makeStyles,
  MenuItem,
  MenuList,
  Paper,
  Popper,
} from "@material-ui/core";
import ArrowDropDownIcon from "@material-ui/icons/ArrowDropDown";
import { FC, useRef, useState } from "react";
import * as React from "react";

const useStyles = makeStyles((theme) => ({
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

export interface SplitButtonProps {
  options: string[];
  label: string;
  onClick: (index: number) => void;
  startIcon?: React.ReactNode;
  disabled: boolean;
  loading: boolean;
  color?: PropTypes.Color;
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
  const classes = useStyles();
  const anchorRef = useRef(null);
  const [openCancelMenu, setOpenCancelMenu] = useState(false);
  const [selectedCancelOption, setSelectedCancelOption] = useState(0);

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
          {loading && (
            <CircularProgress size={24} className={classes.buttonProgress} />
          )}
        </Button>
        <Button
          size="small"
          aria-controls={openCancelMenu ? "split-button-menu" : undefined}
          aria-expanded={openCancelMenu ? "true" : undefined}
          aria-label={label}
          aria-haspopup="menu"
          onClick={() => setOpenCancelMenu(!openCancelMenu)}
        >
          <ArrowDropDownIcon />
        </Button>
      </ButtonGroup>
      <Popper
        open={openCancelMenu}
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
              <ClickAwayListener onClickAway={() => setOpenCancelMenu(false)}>
                <MenuList id="split-button-menu">
                  {options.map((option, index) => (
                    <MenuItem
                      key={option}
                      selected={index === selectedCancelOption}
                      onClick={() => {
                        setSelectedCancelOption(index);
                        setOpenCancelMenu(false);
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
