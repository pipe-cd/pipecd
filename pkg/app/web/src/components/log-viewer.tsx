import {
  Divider,
  makeStyles,
  Toolbar,
  IconButton,
  Typography,
} from "@material-ui/core";
import React, { FC, memo, useCallback, useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { AppState } from "../modules";
import { StageLog, selectStageLogById } from "../modules/stage-logs";
import { Log } from "./log";
import { Close } from "@material-ui/icons";
import { clearActiveStage } from "../modules/active-stage";
import { selectById, Stage, isStageRunning } from "../modules/deployments";
import Draggable from "react-draggable";
import clsx from "clsx";

const INITIAL_HEIGHT = 400;

function useActiveStageLog(): [Stage | null, StageLog | null] {
  return useSelector<AppState, [Stage | null, StageLog | null]>((state) => {
    if (!state.activeStage) {
      return [null, null];
    }

    const deployment = selectById(
      state.deployments,
      state.activeStage.deploymentId
    );

    if (!deployment) {
      return [null, null];
    }

    const stage = deployment.stagesList.find(
      (s) => s.id === state.activeStage?.stageId
    );

    if (!stage) {
      return [null, null];
    }

    return [stage, selectStageLogById(state.stageLogs, state.activeStage)];
  });
}

const useStyles = makeStyles((theme) => ({
  root: {
    position: "absolute",
    bottom: "0px",
    width: "100%",
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
  logContainer: {
    overflowY: "scroll",
  },
  dividerWrapper: {
    width: "100%",
    paddingTop: theme.spacing(0.5),
    paddingBottom: theme.spacing(0.5),
    cursor: "ns-resize",
  },
  handle: {
    position: "absolute",
    // view height + header
    bottom: `${INITIAL_HEIGHT + 48}px`,
    zIndex: 10,
  },
}));

export const LogViewer: FC = memo(function LogViewer() {
  const classes = useStyles();
  const [activeStage, stageLog] = useActiveStageLog();
  const dispatch = useDispatch();
  const [posY, setPosY] = useState(0);
  const [height, setHeight] = useState(INITIAL_HEIGHT);

  const handleOnClickClose = (): void => {
    dispatch(clearActiveStage());
  };

  const handleOnDrag = useCallback(
    (_, data) => {
      setPosY(data.y);
      setHeight(height - (data.y - data.lastY));
    },
    [setPosY, height, setHeight]
  );

  if (!stageLog || !activeStage) {
    return null;
  }

  return (
    <>
      <Draggable
        onDrag={handleOnDrag}
        handle=".handle"
        position={{ x: 0, y: posY }}
        defaultClassName={classes.handle}
        axis="y"
      >
        <div className={clsx("handle", classes.dividerWrapper)}>
          <Divider />
        </div>
      </Draggable>
      <div className={classes.root}>
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
        <div className={classes.logContainer} style={{ height }}>
          <Log
            loading={isStageRunning(activeStage.status)}
            logs={stageLog.logBlocks}
          />
        </div>
      </div>
    </>
  );
});
