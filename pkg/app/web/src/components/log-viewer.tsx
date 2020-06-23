import {
  Divider,
  makeStyles,
  Toolbar,
  IconButton,
  Typography,
} from "@material-ui/core";
import React, { FC, memo } from "react";
import { useSelector, useDispatch } from "react-redux";
import { AppState } from "../modules";
import { StageLog } from "../modules/stage-logs";
import { Log } from "./log";
import { Close } from "@material-ui/icons";
import { clearActiveStage, ActiveStage } from "../modules/active-stage";

function useActiveStageLog(): [ActiveStage, StageLog | null] {
  return useSelector<AppState, [ActiveStage, StageLog | null]>((state) => {
    if (!state.activeStage) {
      return [null, null];
    }
    return [state.activeStage, state.stageLogs[state.activeStage.id]];
  });
}

const useStyles = makeStyles({
  container: {
    overflow: "scroll",
  },
  toolbarLeft: {
    flex: 1,
  },
  toolbarRight: {
    flex: 1,
    justifyContent: "flex-end",
    display: "flex",
  },
  stageName: {
    fontFamily: "Roboto Mono",
  },
});

export const LogViewer: FC = memo(function LogViewer() {
  const classes = useStyles();
  const [activeStage, stageLog] = useActiveStageLog();
  const dispatch = useDispatch();

  const handleOnClickClose = (): void => {
    dispatch(clearActiveStage());
  };

  if (!stageLog || !activeStage) {
    return null;
  }

  return (
    <div className={classes.container}>
      <Divider />
      <Toolbar variant="dense">
        <div className={classes.toolbarLeft}>
          <Typography variant="subtitle2" className={classes.stageName}>
            {activeStage.name}
          </Typography>
        </div>
        <div className={classes.toolbarRight}>
          <IconButton aria-label="close log" onClick={handleOnClickClose}>
            <Close />
          </IconButton>
        </div>
      </Toolbar>
      <Log
        height={400}
        loading={stageLog.completed === false}
        logs={stageLog.logBlocks}
      />
    </div>
  );
});
