import makeStyles from "@mui/styles/makeStyles";

export const useProjectSettingStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
  },
  description: {
    paddingRight: theme.spacing(6),
  },
  titleWithIcon: {
    display: "flex",
    alignItems: "center",
  },
  valuesWrapper: {
    padding: theme.spacing(1),
    display: "flex",
    justifyContent: "space-between",
  },
  values: {
    padding: theme.spacing(2),
  },
}));
