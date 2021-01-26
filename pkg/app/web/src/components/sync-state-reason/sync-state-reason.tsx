import React, { FC, useState } from "react";
import { makeStyles, Paper, Typography, Button } from "@material-ui/core";
import { DiffView } from "../diff-view";

const useStyles = makeStyles((theme) => ({
  summary: {
    display: "flex",
    alignItems: "center",
  },
  detail: {
    padding: theme.spacing(2),
    fontFamily: "Roboto Mono",
    marginTop: theme.spacing(1),
    wordBreak: "break-all",
    overflow: "auto",
    maxHeight: 400,
  },
  showButton: {
    color: theme.palette.primary.light,
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  summary: string;
  detail: string;
}

export const SyncStateReason: FC<Props> = ({ summary, detail }) => {
  const classes = useStyles();
  const [showReason, setShowReason] = useState(false);
  return (
    <>
      <div className={classes.summary}>
        <Typography variant="body2" color="textSecondary">
          {summary}
        </Typography>
        {detail && (
          <Button
            variant="text"
            size="small"
            className={classes.showButton}
            onClick={() => setShowReason(!showReason)}
          >
            {showReason ? "HIDE DETAILS" : "SHOW DETAILS"}
          </Button>
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
