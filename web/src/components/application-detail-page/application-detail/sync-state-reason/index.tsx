import { Box, Button, Paper, Typography } from "@mui/material";
import { FC, useState } from "react";
import { CopyIconButton } from "~/components/copy-icon-button";
import { DiffView } from "./diff-view";

export interface SyncStateReasonProps {
  summary: string;
  detail: string;
}

export const OutOfSyncReason: FC<SyncStateReasonProps> = ({
  summary,
  detail,
}) => {
  const [showReason, setShowReason] = useState(false);
  return (
    <>
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
        }}
      >
        <Typography variant="body2" color="textSecondary">
          {summary}
        </Typography>
        {detail && (
          <>
            <Button
              variant="text"
              size="small"
              sx={{
                color: "primary.light",
                marginLeft: 1,
                marginRight: 1,
              }}
              onClick={() => setShowReason(!showReason)}
            >
              {showReason ? "HIDE DETAILS" : "SHOW DETAILS"}
            </Button>
            {showReason && (
              <CopyIconButton name="Diff" value={detail} size="small" />
            )}
          </>
        )}
      </Box>
      {showReason && (
        <Paper
          elevation={0}
          variant="outlined"
          sx={{
            padding: 2,
            fontFamily: "fontFamilyMono",
            marginTop: 1,
            wordBreak: "break-all",
            overflow: "auto",
            maxHeight: 400,
            fontSize: 14,
          }}
        >
          <DiffView content={detail} />
        </Paper>
      )}
    </>
  );
};

export const InvalidConfigReason: FC<SyncStateReasonProps> = ({ detail }) => {
  const [showReason, setShowReason] = useState(false);

  const msgHeader = "Failed to load application config: ";
  const MAX_DISPLAY_LENGTH = 200;
  if (detail.length < MAX_DISPLAY_LENGTH) {
    return (
      <>
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
          }}
        >
          <Typography variant="body2" color="error">
            {msgHeader}
            <strong>{detail}</strong>
          </Typography>
        </Box>
      </>
    );
  }

  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
      }}
    >
      {showReason ? (
        <Typography variant="body2" color="error">
          {msgHeader}
          <strong>{detail}</strong>
        </Typography>
      ) : (
        <Typography variant="body2" color="error">
          {msgHeader}
          <strong>{detail.slice(0, MAX_DISPLAY_LENGTH) + "..."}</strong>
        </Typography>
      )}
      {detail && (
        <Typography
          variant="body2"
          onClick={() => setShowReason(!showReason)}
          sx={{
            color: "grey.500",
            ml: 0.5,
            cursor: "pointer",
          }}
        >
          {showReason ? "show less" : "show more"}
        </Typography>
      )}
    </Box>
  );
};
