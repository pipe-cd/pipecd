import { Box, styled, Typography } from "@mui/material";

export const ProjectTitleWrap = styled("div")(() => ({
  display: "flex",
  alignItems: "center",
  justifyContent: "space-between",
}));

export const ProjectTitle = styled(Typography)(() => ({
  display: "flex",
  alignItems: "center",
}));

export const ProjectDescription = styled(Typography)(({ theme }) => ({
  paddingRight: theme.spacing(6),
}));

export const ProjectValuesWrapper = styled(Box)(({ theme }) => ({
  padding: theme.spacing(1),
  display: "flex",
  justifyContent: "space-between",
}));

export const ProjectValues = styled(Box)(({ theme }) => ({
  padding: theme.spacing(2),
}));
