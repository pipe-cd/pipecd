import {
  Divider,
  makeStyles,
  Toolbar,
  IconButton,
  Typography,
} from "@material-ui/core";
import React, { FC } from "react";
import { useSelector, useDispatch } from "react-redux";
import { AppState } from "../modules";
import { StageLog } from "../modules/stage-logs";
import { Log } from "./log";
import { Close } from "@material-ui/icons";
import { clearActiveStage } from "../modules/active-stage";

function useActiveStageLog() {
  return useSelector<AppState, StageLog | null>((state) => {
    if (!state.activeStage) {
      return null;
    }
    return state.stageLogs[state.activeStage];
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
  stageId: {
    fontFamily: "Roboto Mono",
  },
});

interface Props {}

export const LogViewer: FC<Props> = ({}) => {
  const classes = useStyles();
  const stageLog = useActiveStageLog();
  const dispatch = useDispatch();

  const handleOnClickClose = () => {
    dispatch(clearActiveStage());
  };

  if (!stageLog) {
    return null;
  }

  return (
    <div className={classes.container}>
      <Divider />
      <Toolbar variant="dense">
        <div className={classes.toolbarLeft}>
          <Typography variant="subtitle2">
            {/** TODO: Show stage name instead of stage id */}
            <div className={classes.stageId}>{stageLog.stageId}</div>
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
        logs={stageLog.logBlocks.map((block) => block.log)}
      />
    </div>
  );
};
