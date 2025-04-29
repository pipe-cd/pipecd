import {
  Divider,
  IconButton,
  Toolbar,
  Typography,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Button,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { Close, SkipNext } from "@mui/icons-material";
import clsx from "clsx";
import { FC, memo, useCallback, useState } from "react";
import Draggable from "react-draggable";
import { APP_HEADER_HEIGHT } from "~/components/header";
import {
  useAppDispatch,
  useShallowEqualSelector,
  useAppSelector,
} from "~/hooks/redux";
import { clearActiveStage } from "~/modules/active-stage";
import {
  isStageRunning,
  selectById,
  Stage,
  StageStatus,
  skipStage,
  selectDeploymentStageIsSkippable,
  updateSkippableState,
} from "~/modules/deployments";
import { selectStageLogById, StageLog } from "~/modules/stage-logs";
import { Log } from "./log";

const INITIAL_HEIGHT = 400;
const TOOLBAR_HEIGHT = 48;
const ANALYSIS_STAGE_NAME = "ANALYSIS";

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
  skipButton: {
    color: theme.palette.common.white,
    background: theme.palette.success.main,
    marginRight: "10px",
    "& .MuiButton-endIcon": {
      marginLeft: 0,
    },
    "&:hover": {
      backgroundColor: theme.palette.success.dark,
    },
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
  const [isOpenSkipDialog, setOpenSkipDialog] = useState(false);
  const stageId = activeStage ? activeStage.id : "";

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

  const handleSkip = (): void => {
    const deploymentId = stageLog ? stageLog.deploymentId : "";
    dispatch(skipStage({ deploymentId: deploymentId, stageId: stageId }));
    dispatch(updateSkippableState({ stageId: stageId, skippable: false }));
    setOpenSkipDialog(false);
  };

  const isSkippable = useAppSelector(selectDeploymentStageIsSkippable(stageId));

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
            {activeStage.name === ANALYSIS_STAGE_NAME &&
              activeStage.status === StageStatus.STAGE_RUNNING && (
                <Button
                  className={classes.skipButton}
                  onClick={() => setOpenSkipDialog(true)}
                  variant="contained"
                  endIcon={<SkipNext />}
                  disabled={!isSkippable}
                >
                  SKIP
                </Button>
              )}
            <Typography variant="subtitle2" className={classes.stageName}>
              {activeStage.name}
            </Typography>
            <Typography variant="body2" className={classes.stageDescription}>
              {activeStage.desc}
            </Typography>
          </div>
          <div className={classes.toolbarRight}>
            <IconButton
              aria-label="close log"
              onClick={handleOnClickClose}
              size="large"
            >
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
      <Dialog open={isOpenSkipDialog} onClose={() => setOpenSkipDialog(false)}>
        <DialogTitle>Skip stage</DialogTitle>
        <DialogContent>
          <DialogContentText>
            {`To skip this stage, click "SKIP".`}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenSkipDialog(false)}>CANCEL</Button>
          <Button color="primary" onClick={handleSkip}>
            SKIP
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
});
