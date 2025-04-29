import makeStyles from "@mui/styles/makeStyles";

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
