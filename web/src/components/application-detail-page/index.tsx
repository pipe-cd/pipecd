import { IconButton, Menu, MenuItem } from "@mui/material";
import DehazeIcon from "@mui/icons-material/Dehaze";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { UI_ENCRYPT_SECRET, UI_TEXT_DELETE } from "~/constants/ui-text";
import { useInterval } from "~/hooks/use-interval";
import {
  Application,
  fetchApplication,
  enableApplication,
  selectById,
} from "~/modules/applications";
import { setDeletingAppId } from "~/modules/delete-application";
import { DeleteApplicationDialog } from "../applications-page/application-list/delete-application-dialog";
import { DisableApplicationDialog } from "../applications-page/application-list/disable-application-dialog";
import { SealedSecretDialog } from "../applications-page/application-list/sealed-secret-dialog";
import { ApplicationDetail } from "./application-detail";
import { ApplicationStateView } from "./application-state-view";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId ?? "");
  const [hasFetchApplicationError] = useAppSelector<[boolean]>((state) => [
    state.applications.fetchApplicationError !== null,
  ]);

  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const [openDisableDialog, setOpenDisableDialog] = useState(false);
  const [openEncryptSecretDialog, setOpenEncryptSecretDialog] = useState(false);

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplication(applicationId));
    }
  }, [applicationId, dispatch]);

  useInterval(
    () => {
      if (applicationId) {
        dispatch(fetchApplication(applicationId));
      }
    },
    applicationId && hasFetchApplicationError === false ? FETCH_INTERVAL : null
  );

  const app = useAppSelector<Application.AsObject | undefined>((state) =>
    selectById(state.applications, applicationId)
  );

  const handleEncryptSecretClick = (): void => {
    setOpenEncryptSecretDialog(true);
    setAnchorEl(null);
  };

  const handleEnableClick = useCallback(async () => {
    await dispatch(enableApplication({ applicationId: applicationId }));
    setAnchorEl(null);
  }, [dispatch, applicationId]);

  const handleDisableClick = (): void => {
    setOpenDisableDialog(true);
    setAnchorEl(null);
  };

  const handleDeleteClick = useCallback(() => {
    dispatch(setDeletingAppId(applicationId));
    setAnchorEl(null);
  }, [dispatch, applicationId]);

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
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
        applicationId={applicationId}
        onClose={() => setOpenEncryptSecretDialog(false)}
      />
      <DisableApplicationDialog
        open={openDisableDialog}
        applicationId={applicationId}
        onDisable={() => setOpenDisableDialog(false)}
        onCancel={() => setOpenDisableDialog(false)}
      />
      <DeleteApplicationDialog
        onDeleted={() => navigate(PAGE_PATH_APPLICATIONS)}
      />
    </>
  );
});
