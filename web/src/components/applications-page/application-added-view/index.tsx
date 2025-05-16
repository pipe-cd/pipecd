import { Box, Button, Divider, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC, memo } from "react";

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  contentBox: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(4),
  },
  note: {
    marginTop: theme.spacing(1),
    color: theme.palette.text.secondary,
  },
}));

const TEXT = {
  TITLE: "Congratulation!",
  MESSAGE: "Your application has been added successfully.",
  NOTE:
    "Please ensure that your application directory in Git is containing an application config file since PipeCD needs it to know how to run the application's deployments.",
};

export interface ApplicationAddedViewProps {
  onClose: () => void;
}

export const ApplicationAddedView: FC<ApplicationAddedViewProps> = memo(
  function ApplicationAddedView({ onClose }) {
    const classes = useStyles();

    return (
      <Box width={600} flex={1} display="flex" flexDirection="column">
        <Typography className={classes.title} variant="h6">
          {TEXT.TITLE}
        </Typography>

        <Divider />

        <Box p={2}>
          <Box className={classes.contentBox}>
            <Typography variant="subtitle1">{TEXT.MESSAGE}</Typography>
            <Typography variant="body2" className={classes.note}>
              {TEXT.NOTE}
            </Typography>
          </Box>

          <Box mt={1} textAlign="right">
            <Button onClick={onClose} variant="outlined">
              CLOSE
            </Button>
          </Box>
        </Box>
      </Box>
    );
  }
);
