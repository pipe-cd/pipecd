import { Box, styled } from "@mui/material";

export const StyledForm = styled("form")(({ theme }) => ({
  padding: theme.spacing(2),
  display: "grid",
  gap: theme.spacing(2),
}));

export const GroupTwoCol = styled(Box)(({ theme }) => ({
  display: "grid",
  gridTemplateColumns: "1fr 1fr",
  gap: theme.spacing(2),
}));
