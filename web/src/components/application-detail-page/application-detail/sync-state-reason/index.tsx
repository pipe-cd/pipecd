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

  const MAX_DISPLAY_LENGTH = 100;
  if (detail.length < MAX_DISPLAY_LENGTH) {
    return (
      <>
        <div className={classes.summary}>
          <Typography variant="body2" color="textSecondary">
            {detail}
          </Typography>
        </div>
      </>
    );
  }

  return (
    <>
      <div className={classes.summary}>
        {showReason ? (
          <Typography variant="body2" color="textSecondary">
            {detail}
          </Typography>
        ) : (
          <Typography variant="body2" color="textSecondary">
            {detail.slice(0, MAX_DISPLAY_LENGTH) + "..."}
          </Typography>
        )}
        <>
          <Button
            variant="text"
            size="small"
            className={classes.showButton}
            onClick={() => setShowReason(!showReason)}
          >
            {showReason ? "HIDE" : "MORE"}
          </Button>
        </>
      </div>
    </>
  );
};
