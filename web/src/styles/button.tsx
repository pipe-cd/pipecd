import { CircularProgress, CircularProgressProps, styled } from "@mui/material";
import { makeStyles } from "@mui/styles";

export const useStyles = makeStyles((theme) => ({
  progress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

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
