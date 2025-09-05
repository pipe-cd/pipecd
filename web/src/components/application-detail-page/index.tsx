import { IconButton, Menu, MenuItem } from "@mui/material";
import DehazeIcon from "@mui/icons-material/Dehaze";
import { FC, memo, useCallback, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { DeleteApplicationDialog } from "../applications-page/application-list/delete-application-dialog";
import { DisableApplicationDialog } from "../applications-page/application-list/disable-application-dialog";
import { SealedSecretDialog } from "../applications-page/application-list/sealed-secret-dialog";
import { ApplicationDetail } from "./application-detail";
import { ApplicationStateView } from "./application-state-view";
import { useGetApplicationStateById } from "~/queries/application-live-state/use-get-application-state-by-id";
import { useGetApplicationDetail } from "~/queries/applications/use-get-application-detail";
import { useEnableApplication } from "~/queries/applications/use-enable-application";
import { checkPipedAppVersion } from "~/utils/common";
import { ApplicationKind } from "~~/model/common_pb";
import { PIPED_VERSION } from "~/types/piped";
import { Application } from "~/types/applications";
import { UI_ENCRYPT_SECRET, UI_TEXT_DELETE } from "~/constants/ui-text";

const FETCH_INTERVAL = 4000;

const isDisplayLiveState = (app: Application.AsObject | undefined): boolean => {
  const result = checkPipedAppVersion(app);
  if (result[PIPED_VERSION.V1]) return true;

  return (
    app?.kind === ApplicationKind.KUBERNETES ||
    app?.kind === ApplicationKind.CLOUDRUN ||
    app?.kind === ApplicationKind.ECS ||
    app?.kind === ApplicationKind.LAMBDA
  );
};

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const navigate = useNavigate();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId ?? "");
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const [openDisableDialog, setOpenDisableDialog] = useState(false);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [openEncryptSecretDialog, setOpenEncryptSecretDialog] = useState(false);

  const {
    data: app,
    isError: isAppError,
    refetch: refetchApp,
  } = useGetApplicationDetail(applicationId, {
    enabled: !!applicationId,
    refetchInterval: (_data, query) =>
      !query.state.error ? FETCH_INTERVAL : false,
  });

  const {
    data: liveState,
    isError: isLiveStateError,
    isLoading: isLiveStateLoading,
    refetch: refetchLiveState,
  } = useGetApplicationStateById(applicationId, {
    enabled: isDisplayLiveState(app),
    retry: false,
    refetchInterval: (_data, query) =>
      !query.state.error ? FETCH_INTERVAL : false,
  });

  const { mutate: enableApplication } = useEnableApplication();

  const handleEncryptSecretClick = (): void => {
    setOpenEncryptSecretDialog(true);
    setAnchorEl(null);
  };

  const handleEnableClick = useCallback(() => {
    enableApplication(
      { applicationId },
      { onSuccess: () => setAnchorEl(null) }
    );
  }, [enableApplication, applicationId]);

  const handleDisableClick = (): void => {
    setOpenDisableDialog(true);
    setAnchorEl(null);
  };

  const handleDeleteClick = useCallback(() => {
    setOpenDeleteDialog(true);
    setAnchorEl(null);
  }, []);

  return (
    <>
      <ApplicationDetail
        app={app}
        hasError={isAppError}
        refetchApp={refetchApp}
        liveState={liveState}
        liveStateLoading={isLiveStateLoading}
      />
      <ApplicationStateView
        app={app}
        hasError={isLiveStateError}
        liveState={liveState}
        refetchLiveState={refetchLiveState}
      />
      <IconButton
        aria-label="Open menu"
        onClick={(e) => {
          setAnchorEl(e.currentTarget);
        }}
        sx={{
          backgroundColor: "#283778",
          position: "absolute",
          bottom: 30,
          right: 30,
          "&:hover": {
            backgroundColor: "grey",
          },
        }}
        size="large"
      >
        <DehazeIcon fontSize="large" htmlColor="#fff" />
      </IconButton>
      <Menu
        id="action-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
        transformOrigin={{
          vertical: "bottom",
          horizontal: 150,
        }}
      >
        {app && app.disabled ? (
          <MenuItem onClick={handleEnableClick}>Enable</MenuItem>
        ) : (
          <div>
            <MenuItem onClick={handleEncryptSecretClick}>
              {UI_ENCRYPT_SECRET}
            </MenuItem>
            <MenuItem onClick={handleDisableClick}>Disable</MenuItem>
          </div>
        )}
        <MenuItem onClick={handleDeleteClick} sx={{ color: "red" }}>
          {UI_TEXT_DELETE}
        </MenuItem>
      </Menu>
      <SealedSecretDialog
        open={openEncryptSecretDialog}
        application={app}
        onClose={() => setOpenEncryptSecretDialog(false)}
      />
      <DisableApplicationDialog
        open={openDisableDialog}
        application={app}
        onDisable={() => setOpenDisableDialog(false)}
        onCancel={() => setOpenDisableDialog(false)}
      />
      <DeleteApplicationDialog
        open={openDeleteDialog}
        application={app}
        onDeleted={() => navigate(PAGE_PATH_APPLICATIONS)}
        onCancel={() => setOpenDeleteDialog(false)}
      />
    </>
  );
});
