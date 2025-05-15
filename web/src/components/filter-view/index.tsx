import { FC } from "react";
import * as React from "react";
import { Box, Button, Paper, Typography } from "@mui/material";
import { UI_TEXT_CLEAR } from "~/constants/ui-text";

export interface FilterViewProps {
  onClear: () => void;
  children: React.ReactNode;
}

export const FilterView: FC<FilterViewProps> = ({ onClear, children }) => {
  return (
    <Paper
      square
      sx={{
        width: 360,
        minWidth: 280,
        padding: 3,
        height: "100%",
      }}
    >
      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
        }}
      >
        <Typography variant="h6" component="span">
          Filters
        </Typography>
        <Button color="primary" onClick={onClear}>
          {UI_TEXT_CLEAR}
        </Button>
      </Box>
      {children}
    </Paper>
  );
};
