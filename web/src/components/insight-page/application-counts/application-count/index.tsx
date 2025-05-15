import { Box, Card, CardActionArea, Popover, Typography } from "@mui/material";
import { FC, memo, useState } from "react";

export interface ApplicationCountProps {
  enabledCount: number;
  disabledCount: number;
  kindName: string;
  onClick: () => void;
}

export const ApplicationCount: FC<ApplicationCountProps> = memo(
  function ApplicationCount({
    enabledCount,
    disabledCount,
    kindName,
    onClick,
  }) {
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const open = Boolean(anchorEl);

    const handlePopoverClose = (): void => {
      setAnchorEl(null);
    };

    return (
      <Card
        raised
        sx={{
          minWidth: 200,
          display: "inline-block",
        }}
      >
        <CardActionArea
          sx={{ p: 2 }}
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
                // className={classes.textSpace}
                sx={{ ml: 1 }}
              >
                {`/${disabledCount}`}
              </Typography>
            ) : null}
            <Typography
              variant="h6"
              component="span"
              sx={{ ml: 1 }}
              // className={classes.textSpace}
            >
              apps
            </Typography>
          </Box>
        </CardActionArea>

        <Popover
          id="mouse-over-popover"
          sx={{
            pointerEvents: "none",
          }}
          slotProps={{
            paper: { sx: { p: 1 } }, // Add padding to the popover
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
