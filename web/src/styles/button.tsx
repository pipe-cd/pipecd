import { CircularProgress, CircularProgressProps, styled } from "@mui/material";

export const SpinnerIcon = styled(
  ({ size = 24, ...props }: CircularProgressProps) => (
    <CircularProgress size={size} {...props} />
  )
)(({ theme, size = 24 }) => ({
  color: theme.palette.primary.main,
  position: "absolute",
  top: "50%",
  left: "50%",
  marginTop: -size / 2,
  marginLeft: -size / 2,
}));
