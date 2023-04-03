import { Button, makeStyles, Paper, Typography } from "@material-ui/core";
import { FC, useState } from "react";
import { CopyIconButton } from "~/components/copy-icon-button";
import { DiffView } from "./diff-view";

const useStyles = makeStyles((theme) => ({
  summary: {
    display: "flex",
    alignItems: "center",
  },
  detail: {
    padding: theme.spacing(2),
    fontFamily: theme.typography.fontFamilyMono,
    marginTop: theme.spacing(1),
    wordBreak: "break-all",
    overflow: "auto",
    maxHeight: 400,
  },
  showButton: {
    color: theme.palette.primary.light,
    marginLeft: theme.spacing(1),
    marginRight: theme.spacing(1),
  },
  showText: {
    color: theme.palette.grey[500],
    marginLeft: theme.spacing(0.5),
    cursor: "pointer",
  },
}));

export interface SyncStateReasonProps {
  summary: string;
  detail: string;
}

export const OutOfSyncReason: FC<SyncStateReasonProps> = ({
  summary,
  detail,
}) => {
  const classes = useStyles();
  const [showReason, setShowReason] = useState(false);
  return (
    <>
      <div className={classes.summary}>
        <Typography variant="body2" color="textSecondary">
          {summary}
        </Typography>
        {detail && (
          <>
            <Button
              variant="text"
              size="small"
              className={classes.showButton}
              onClick={() => setShowReason(!showReason)}
            >
              {showReason ? "HIDE DETAILS" : "SHOW DETAILS"}
            </Button>
            {showReason && (
              <CopyIconButton name="Diff" value={detail} size="small" />
            )}
          </>
        )}
      </div>

      {showReason && (
        <Paper elevation={0} variant="outlined" className={classes.detail}>
          <DiffView content={detail} />
        </Paper>
      )}
    </>
  );
};

export const InvalidConfigReason: FC<SyncStateReasonProps> = ({ detail }) => {
  const classes = useStyles();
  const [showReason, setShowReason] = useState(false);

  const msgHeader = "Failed to load aplication config: ";
  const MAX_DISPLAY_LENGTH = 200;
  if (detail.length < MAX_DISPLAY_LENGTH) {
    return (
      <>
        <div className={classes.summary}>
          <Typography variant="body2" color="error">
            {msgHeader}
            <strong>{detail}</strong>
          </Typography>
        </div>
      </>
    );
  }

  return (
    <>
      <div className={classes.summary}>
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
          <span
            className={classes.showText}
            onClick={() => setShowReason(!showReason)}
          >
            {showReason ? "show less" : "show more"}
          </span>
        )}
      </div>
    </>
  );
};
