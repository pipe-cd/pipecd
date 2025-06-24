import { FC } from "react";
import { Box, Typography } from "@mui/material";

export interface ProjectSettingLabeledTextProps {
  label: string;
  value: string;
}

export const ProjectSettingLabeledText: FC<ProjectSettingLabeledTextProps> = ({
  label,
  value,
}) => {
  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
      }}
    >
      <Typography
        variant="subtitle1"
        sx={{
          color: "text.secondary",
          marginRight: 2,
          minWidth: "120px",
        }}
      >
        {label}
      </Typography>
      <Typography variant="body1">{value}</Typography>
    </Box>
  );
};
