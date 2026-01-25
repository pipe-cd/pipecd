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
  Box,
} from "@mui/material";
import { Close, SkipNext } from "@mui/icons-material";
import { FC, memo, useCallback, useMemo, useState } from "react";
import Draggable from "react-draggable";
import { APP_HEADER_HEIGHT } from "~/components/header";
import { PipelineStage, StageStatus } from "pipecd/web/model/deployment_pb";
import { Log } from "./log";
import { ManualOperation } from "~~/model/deployment_pb";
import { StageLog } from "~/types/stage-log";
import { useSkipStage } from "~/queries/deployment/use-skip-stage";
import { isStageRunning } from "~/utils/is-stage-running";
import { useCommand } from "~/contexts/command-context";
import { Command, CommandStatus } from "~~/model/command_pb";
import { ActiveStageInfo } from "..";

const INITIAL_HEIGHT = 400;
const TOOLBAR_HEIGHT = 48;
const ANALYSIS_STAGE_NAME = "ANALYSIS";

type Props = {
  activeStage: PipelineStage.AsObject | null;
  changeActiveStage: (activeStage: ActiveStageInfo | null) => void;
  stageLog?: StageLog | null;
};

export const LogViewer: FC<Props> = memo(function LogViewer({
  activeStage,
  stageLog,
  changeActiveStage,
}) {
  const maxHandlePosY =
    document.body.clientHeight - APP_HEADER_HEIGHT - TOOLBAR_HEIGHT;
  const [handlePosY, setHandlePosY] = useState(maxHandlePosY - INITIAL_HEIGHT);
  const logViewHeight = maxHandlePosY - handlePosY;
  const [isOpenSkipDialog, setOpenSkipDialog] = useState(false);
  const [commandId, setCommandId] = useState<string>();

  const stageId = activeStage ? activeStage.id : "";
  const { mutate: skipStage, isLoading: isSkipInitLoading } = useSkipStage();
  const { fetchedCommands, commandIds } = useCommand();

  const handleOnClickClose = (): void => {
    changeActiveStage(null);
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
    skipStage(
      { deploymentId: deploymentId, stageId: stageId },
      {
        onSuccess: (commandId) => {
          setCommandId(commandId);
        },
      }
    );
    setOpenSkipDialog(false);
  };

  const isSkippable = useMemo(() => {
    const deploymentId = stageLog ? stageLog.deploymentId : "";

    const isSkipCommandRunning = (deploymentId: string): boolean => {
      const stageCommand = Object.values(fetchedCommands).find(
        (item) =>
          item.deploymentId === deploymentId &&
          item.stageId === stageId &&
          item.type === Command.Type.SKIP_STAGE
      );
      if (
        stageCommand &&
        stageCommand?.status !== CommandStatus.COMMAND_FAILED &&
        stageCommand?.status !== CommandStatus.COMMAND_TIMEOUT
      ) {
        return true;
      }

      return false;
    };

    const isSkipCommandInit = commandIds?.has(commandId ?? "");

    return (
      !isSkipInitLoading &&
      !isSkipCommandInit &&
      !isSkipCommandRunning(deploymentId)
    );
  }, [
    commandId,
    commandIds,
    fetchedCommands,
    isSkipInitLoading,
    stageId,
    stageLog,
  ]);

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
        axis="y"
      >
        <Box
          className={"handle"}
          sx={(theme) => ({
            position: "absolute",
            zIndex: 10,
            width: "100%",
            paddingTop: theme.spacing(0.5),
            paddingBottom: theme.spacing(0.5),
            cursor: "ns-resize",
          })}
        />
      </Draggable>
      <Box
        data-testid="log-viewer"
        sx={{
          position: "absolute",
          bottom: "0px",
          width: "100%",
        }}
      >
        <Divider />
        <Toolbar variant="dense" sx={{ backgroundColor: "background.default" }}>
          <Box
            sx={{
              flex: 1,
              display: "flex",
              alignItems: "center",
            }}
          >
            {/* TODO: Remove stageName condition after finishing deployments which are made
                      while the server does not inject availableOperation */}
            {(activeStage.name === ANALYSIS_STAGE_NAME ||
              activeStage.availableOperation ===
                ManualOperation.MANUAL_OPERATION_SKIP) &&
              activeStage.status === StageStatus.STAGE_RUNNING && (
                <Button
                  sx={(theme) => ({
                    color: theme.palette.common.white,
                    background: theme.palette.success.main,
                    marginRight: "10px",
                    "& .MuiButton-endIcon": {
                      marginLeft: 0,
                    },
                    "&:hover": {
                      backgroundColor: theme.palette.success.dark,
                    },
                  })}
                  onClick={() => setOpenSkipDialog(true)}
                  variant="contained"
                  endIcon={<SkipNext />}
                  disabled={!isSkippable}
                >
                  SKIP
                </Button>
              )}
            <Typography
              variant="subtitle2"
              sx={{
                fontFamily: "fontFamilyMono",
              }}
            >
              {activeStage.name}
            </Typography>
            <Typography
              variant="body2"
              sx={{
                color: "text.secondary",
                ml: 2,
              }}
            >
              {activeStage.desc}
            </Typography>
          </Box>
          <Box sx={{ flex: 1, justifyContent: "flex-end", display: "flex" }}>
            <IconButton
              aria-label="close log"
              onClick={handleOnClickClose}
              size="large"
            >
              <Close />
            </IconButton>
          </Box>
        </Toolbar>
        <Box
          sx={{
            overflowY: "scroll",
            height: logViewHeight,
          }}
        >
          <Log
            loading={isStageRunning(activeStage.status)}
            logs={stageLog.logBlocks}
          />
        </Box>
      </Box>
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
