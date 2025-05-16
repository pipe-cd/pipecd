import { Box, Button, Divider, Link, Toolbar, Typography } from "@mui/material";
import RefreshIcon from "@mui/icons-material/Refresh";
import { FC } from "react";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { SpinnerIcon } from "~/styles/button";

const COMING_SOON_MESSAGE = "This UI is under development";
const FEATURE_STATUS_INTRO = "PipeCD feature status";

export const DeploymentChainsIndexPage: FC = () => {
  const isLoading = false;

  return (
    <Box
      sx={{
        display: "flex",
        overflow: "hidden",
        flex: 1,
        flexDirection: "column",
      }}
    >
      <Toolbar variant="dense">
        <Box
          sx={{
            flexGrow: 1,
          }}
        />
        <Button
          color="primary"
          startIcon={<RefreshIcon />}
          disabled={isLoading}
        >
          {UI_TEXT_REFRESH}
          {isLoading && <SpinnerIcon />}
        </Button>
      </Toolbar>
      <Divider />
      <Box
        sx={{
          flex: 1,
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          flexDirection: "column",
        }}
      >
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
