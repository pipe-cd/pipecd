import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Link,
  Toolbar,
  Typography,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import RefreshIcon from "@mui/icons-material/Refresh";
import { FC } from "react";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { useStyles as useButtonStyles } from "~/styles/button";

const COMING_SOON_MESSAGE = "This UI is under development";
const FEATURE_STATUS_INTRO = "PipeCD feature status";

const useStyles = makeStyles(() => ({
  container: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
}));

export const DeploymentChainsIndexPage: FC = () => {
  const classes = useStyles();
  const buttonClasses = useButtonStyles();
  const isLoading = false;

  return (
    <Box display="flex" overflow="hidden" flex={1} flexDirection="column">
      <Toolbar variant="dense">
        <Box flexGrow={1} />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
          {isLoading && (
            <CircularProgress size={24} className={buttonClasses.progress} />
          )}
        </Button>
      </Toolbar>
      <Divider />
      <Box flexDirection="column" className={classes.container}>
        <Typography variant="body1">{COMING_SOON_MESSAGE}</Typography>
        <Link
          href="https://pipecd.dev/docs/feature-status/"
          target="_blank"
          rel="noreferrer"
        >
          {FEATURE_STATUS_INTRO}
        </Link>
      </Box>
    </Box>
  );
};
