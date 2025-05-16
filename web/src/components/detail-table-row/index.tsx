import { FC } from "react";
import * as React from "react";
import { Box, Typography } from "@mui/material";

export interface DetailTableRowProps {
  label: string;
  value: React.ReactChild;
}

export const DetailTableRow: FC<DetailTableRowProps> = ({ label, value }) => {
  return (
    <tr>
      <Box
        component="th"
        sx={{
          textAlign: "left",
          whiteSpace: "nowrap",
          display: "block",
        }}
      >
        <Typography variant="subtitle2">{label}</Typography>
      </Box>
      <td>
        <Typography variant="body2" ml={1}>
          {value}
        </Typography>
      </td>
    </tr>
  );
};
