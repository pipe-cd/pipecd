import { IconButton, makeStyles, Menu, MenuItem } from "@material-ui/core";
import DehazeIcon from "@material-ui/icons/Dehaze";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
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

const useStyles = makeStyles(() => ({
  actionsMenuBtn: {
    backgroundColor: "#283778",
    position: "absolute",
    bottom: "30px",
    right: "30px",
    "&:hover": {
      backgroundColor: "grey",
    },
  },
  warning: {
    color: "red",
  },
}));

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const classes = useStyles();
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
        className={classes.actionsMenuBtn}
        aria-label="Open menu"
        onClick={(e) => {
          setAnchorEl(e.currentTarget);
        }}
      >
        <DehazeIcon fontSize="large" htmlColor="#fff" />
      </IconButton>

      <Menu
        id="action-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={() => setAnchorEl(null)}
        PaperProps={{
          style: {
            width: "20ch",
            transform: "translateX(-50%) translateY(-20%)",
          },
        }}
      >
        {app && app.disabled ? (
          <MenuItem onClick={handleEnableClick}>Enable</MenuItem>
        ) : (
          <div>
            <MenuItem onClick={handleEncryptSecretClick}>
              Encrypt Secret
            </MenuItem>
            <MenuItem onClick={handleDisableClick}>Disable</MenuItem>
          </div>
        )}
        <MenuItem className={classes.warning} onClick={handleDeleteClick}>
          Delete
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
