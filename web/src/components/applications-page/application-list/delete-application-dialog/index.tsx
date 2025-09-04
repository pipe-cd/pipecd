import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FC, Fragment, memo, useCallback } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_DELETE } from "~/constants/ui-text";

import Alert from "@mui/material/Alert";
import { DELETE_APPLICATION_SUCCESS } from "~/constants/toast-text";
import { Skeleton } from "@mui/material";
import { red } from "@mui/material/colors";
import { SpinnerIcon } from "~/styles/button";
import { useDeleteApplication } from "~/queries/applications/use-delete-application";
import { useToast } from "~/contexts/toast-context";
import { Application } from "~/types/applications";

const TITLE = "Delete Application";
const ALERT_TEXT = "Are you sure you want to delete the application?";

export interface DeleteApplicationDialogProps {
  open: boolean;
  application?: Application.AsObject | null;
  onDeleted: () => void;
  onCancel: () => void;
}

export const DeleteApplicationDialog: FC<DeleteApplicationDialogProps> = memo(
  function DeleteApplicationDialog({ onDeleted, application, open, onCancel }) {
    const { addToast } = useToast();
    const {
      mutate: deleteApplication,
      isPending: isDeleting,
    } = useDeleteApplication();

    const handleDelete = useCallback(() => {
      if (!application) return;

      deleteApplication(
        { applicationId: application?.id },
        {
          onSuccess: () => {
            onDeleted();
            addToast({
              severity: "success",
              message: DELETE_APPLICATION_SUCCESS,
            });
          },
        }
      );
    }, [addToast, application, deleteApplication, onDeleted]);

    const renderLabels = useCallback(() => {
      if (!application?.labelsMap) return <Skeleton height={24} width={200} />;

      if (application.labelsMap.length === 0) return "-";

      return application.labelsMap.map(([key, value]) => (
        <Fragment key={key}>
          <span>{`${key}: ${value}`}</span>
          <br />
        </Fragment>
      ));
    }, [application?.labelsMap]);

    return (
      <Dialog
        open={Boolean(application) && open}
        onClose={(_event, reason) => {
          if (reason !== "backdropClick" || !isDeleting) {
            onCancel();
          }
        }}
      >
        <DialogTitle>{TITLE}</DialogTitle>
        <DialogContent>
          <Alert severity="error" sx={{ mb: 2 }}>
            {ALERT_TEXT}
          </Alert>
          <Typography variant="caption">Name</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
            {application ? (
              application.name
            ) : (
              <Skeleton height={24} width={200} />
            )}
          </Typography>
          <Box
            sx={{
              height: 24,
            }}
          />
          <Typography variant="caption">Labels</Typography>
          <Typography
            variant="body1"
            sx={(theme) => ({
              color: theme.palette.text.primary,
              fontWeight: theme.typography.fontWeightMedium,
            })}
          >
            {renderLabels()}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={onCancel} disabled={isDeleting}>
            {UI_TEXT_CANCEL}
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={handleDelete}
            sx={(theme) => ({
              color: theme.palette.getContrastText(red[400]),
              backgroundColor: red[800],
              "&:hover": {
                backgroundColor: red[800],
              },
            })}
            disabled={isDeleting}
          >
            {UI_TEXT_DELETE}
            {isDeleting && <SpinnerIcon />}
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
);
