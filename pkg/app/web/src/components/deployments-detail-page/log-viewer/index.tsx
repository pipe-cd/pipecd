import {
  Divider,
  IconButton,
  makeStyles,
  Toolbar,
  Typography,
} from "@material-ui/core";
import { Close } from "@material-ui/icons";
import clsx from "clsx";
import { FC, memo, useCallback, useState } from "react";
import Draggable from "react-draggable";
import { APP_HEADER_HEIGHT } from "~/components/header";
import { useAppDispatch, useShallowEqualSelector } from "~/hooks/redux";
import { clearActiveStage } from "~/modules/active-stage";
import { isStageRunning, selectById, Stage } from "~/modules/deployments";
import { selectStageLogById, StageLog } from "~/modules/stage-logs";
import { Log } from "./log";

const INITIAL_HEIGHT = 400;
const TOOLBAR_HEIGHT = 48;

function useActiveStageLog(): [Stage | null, StageLog | null] {
  return useShallowEqualSelector<[Stage | null, StageLog | null]>((state) => {
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
  toolbar: {
    background: theme.palette.background.default,
  },
  toolbarLeft: {
    flex: 1,
    display: "flex",
    alignItems: "center",
  },
  toolbarRight: {
    flex: 1,
    justifyContent: "flex-end",
    display: "flex",
  },
  stageName: {
    fontFamily: theme.typography.fontFamilyMono,
  },
  stageDescription: {
    marginLeft: theme.spacing(2),
    color: theme.palette.text.secondary,
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
    zIndex: 10,
  },
}));

export const LogViewer: FC = memo(function LogViewer() {
  const maxHandlePosY =
    document.body.clientHeight - APP_HEADER_HEIGHT - TOOLBAR_HEIGHT;
  const classes = useStyles();
  const [activeStage, stageLog] = useActiveStageLog();
  const dispatch = useAppDispatch();
  const [handlePosY, setHandlePosY] = useState(maxHandlePosY - INITIAL_HEIGHT);
  const logViewHeight = maxHandlePosY - handlePosY;

  const handleOnClickClose = (): void => {
    dispatch(clearActiveStage());
  };

  const handleDrag = useCallback(
    (_, data) => {
      if (data.y < 0) {
        setHandlePosY(0);
      } else if (data.y > maxHandlePosY) {
        setHandlePosY(maxHandlePosY);
      } else {
        setHandlePosY(data.y);
      }
    },
    [setHandlePosY, maxHandlePosY]
  );

  if (!stageLog || !activeStage) {
    return null;
  }

  return (
    <>
      <Draggable
        onDrag={handleDrag}
        onStop={handleDrag}
        handle=".handle"
        position={{ x: 0, y: handlePosY }}
        defaultClassName={classes.handle}
        axis="y"
      >
        <div className={clsx("handle", classes.dividerWrapper)} />
      </Draggable>

      <div className={classes.root} data-testid="log-viewer">
        <Divider />
        <Toolbar variant="dense" className={classes.toolbar}>
          <div className={classes.toolbarLeft}>
            <Typography variant="subtitle2" className={classes.stageName}>
              {activeStage.name}
            </Typography>
            <Typography variant="body2" className={classes.stageDescription}>
              {activeStage.desc}
            </Typography>
          </div>
          <div className={classes.toolbarRight}>
            <IconButton aria-label="close log" onClick={handleOnClickClose}>
              <Close />
            </IconButton>
          </div>
        </Toolbar>
        <div className={classes.logContainer} style={{ height: logViewHeight }}>
          <Log
            loading={isStageRunning(activeStage.status)}
            logs={stageLog.logBlocks}
          />
        </div>
      </div>
    </>
  );
});
