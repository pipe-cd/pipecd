import { Box, Card, CardActionArea, Popover, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import clsx from "clsx";
import { FC, memo, useState } from "react";

const useStyles = makeStyles((theme) => ({
  root: {
    minWidth: 200,
    display: "inline-block",
  },
  actionArea: {
    padding: theme.spacing(2),
  },
  textSpace: {
    marginLeft: theme.spacing(1),
  },
  popover: {
    pointerEvents: "none",
  },
  popoverPaper: {
    padding: theme.spacing(1),
  },
}));

export interface ApplicationCountProps {
  enabledCount: number;
  disabledCount: number;
  kindName: string;
  onClick: () => void;
  className?: string;
}

export const ApplicationCount: FC<ApplicationCountProps> = memo(
  function ApplicationCount({
    enabledCount,
    disabledCount,
    kindName,
    onClick,
    className,
  }) {
    const classes = useStyles();

    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const open = Boolean(anchorEl);

    const handlePopoverClose = (): void => {
      setAnchorEl(null);
    };

    return (
      <Card raised className={clsx(classes.root, className)}>
        <CardActionArea
          className={classes.actionArea}
          onClick={onClick}
          onMouseEnter={(event) => {
            setAnchorEl(event.currentTarget);
          }}
          onMouseLeave={handlePopoverClose}
        >
          <Typography variant="h6" component="div" color="textSecondary">
            {kindName}
          </Typography>
          <Box display="flex" justifyContent="center" alignItems="baseline">
            <Typography variant="h4" component="span">
              {enabledCount}
            </Typography>
            {disabledCount > 0 ? (
              <Typography
                variant="h6"
                color="textSecondary"
                component="span"
                className={classes.textSpace}
              >
                {`/${disabledCount}`}
              </Typography>
            ) : null}
            <Typography
              variant="h6"
              component="span"
              className={classes.textSpace}
            >
              apps
            </Typography>
          </Box>
        </CardActionArea>

        <Popover
          id="mouse-over-popover"
          className={classes.popover}
          classes={{
            paper: classes.popoverPaper,
          }}
          open={open}
          anchorEl={anchorEl}
          anchorOrigin={{
            vertical: "bottom",
            horizontal: "center",
          }}
          transformOrigin={{
            vertical: "top",
            horizontal: "left",
          }}
          onClose={handlePopoverClose}
          disableRestoreFocus
        >
          <div>
            <b>{enabledCount}</b>
            {" enabled applications"}
          </div>
          <div>
            <b>{disabledCount}</b>
            {" disabled applications"}
          </div>
        </Popover>
      </Card>
    );
  }
);
