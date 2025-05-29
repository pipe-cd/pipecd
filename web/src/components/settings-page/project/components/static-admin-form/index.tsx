import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Skeleton,
  Switch,
  TextField,
} from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";

import { useFormik } from "formik";
import { FC, memo, useState } from "react";
import * as yup from "yup";
import { STATIC_ADMIN_DESCRIPTION } from "~/constants/text";
import { UPDATE_STATIC_ADMIN_INFO_SUCCESS } from "~/constants/toast-text";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  fetchProject,
  toggleAvailability,
  updateStaticAdmin,
} from "~/modules/project";
import { addToast } from "~/modules/toasts";
import {
  ProjectDescription,
  ProjectTitleWrap,
  ProjectTitle,
  ProjectValues,
  ProjectValuesWrapper,
} from "~/styles/project-setting";
import { ProjectSettingLabeledText } from "../project-setting-labeled-text";

const SECTION_TITLE = "Static Admin";
const DIALOG_TITLE = `Edit ${SECTION_TITLE}`;

const validationSchema = yup.object().shape({
  username: yup.string().min(1).required(),
  password: yup.string().min(1).required(),
});

const StaticAdminDialog: FC<{
  open: boolean;
  currentUsername: string;
  onClose: () => void;
  onSubmit: (values: { username: string; password: string }) => void;
}> = ({ open, currentUsername, onClose, onSubmit }) => {
  const formik = useFormik({
    initialValues: {
      username: currentUsername,
      password: "",
    },
    validationSchema,
    onSubmit,
  });

  return (
    <Dialog
      open={open}
      TransitionProps={{
        onExited: () => {
          formik.resetForm();
        },
      }}
      onClose={onClose}
    >
      <form onSubmit={formik.handleSubmit}>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <TextField
            id="username"
            name="username"
            value={formik.values.username}
            variant="outlined"
            margin="dense"
            label="Username"
            fullWidth
            required
            autoFocus
            onChange={formik.handleChange}
          />
          <TextField
            id="password"
            name="password"
            value={formik.values.password}
            autoComplete="new-password"
            variant="outlined"
            margin="dense"
            label="Password"
            type="password"
            fullWidth
            required
            onChange={formik.handleChange}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>{UI_TEXT_CANCEL}</Button>
          <Button
            type="submit"
            color="primary"
            disabled={formik.isValid === false}
          >
            {UI_TEXT_SAVE}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export const StaticAdminForm: FC = memo(function StaticAdminForm() {
  const dispatch = useAppDispatch();
  const [isEnabled, currentUsername] = useAppSelector<[boolean, string | null]>(
    (state) => [
      state.project.staticAdminDisabled === false,
      state.project.username,
    ]
  );
  const [isEdit, setIsEdit] = useState(false);

  const handleSubmit = (values: {
    username: string;
    password: string;
  }): void => {
    dispatch(updateStaticAdmin(values)).then((result) => {
      if (updateStaticAdmin.fulfilled.match(result)) {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: UPDATE_STATIC_ADMIN_INFO_SUCCESS,
            severity: "success",
          })
        );
      }
    });
    setIsEdit(false);
  };

  const handleClose = (): void => {
    setIsEdit(false);
  };

  const handleToggleAvailability = (): void => {
    dispatch(toggleAvailability()).then(() => {
      dispatch(fetchProject());
    });
  };

  return (
    <>
      <ProjectTitleWrap>
        <ProjectTitle variant="h5">{SECTION_TITLE}</ProjectTitle>

        <Switch
          checked={isEnabled}
          color="primary"
          onClick={handleToggleAvailability}
          disabled={currentUsername === null}
        />
      </ProjectTitleWrap>
      <ProjectDescription variant="body1" color="textSecondary">
        {STATIC_ADMIN_DESCRIPTION}
      </ProjectDescription>
      <ProjectValuesWrapper sx={{ opacity: isEnabled === false ? 0.5 : 1 }}>
        {currentUsername ? (
          <>
            <ProjectValues>
              <ProjectSettingLabeledText
                label="Username"
                value={currentUsername}
              />
              <ProjectSettingLabeledText label="Password" value="********" />
            </ProjectValues>
            <div>
              <IconButton
                aria-label="edit static admin user"
                onClick={() => setIsEdit(true)}
                disabled={isEnabled === false}
                size="large"
              >
                <EditIcon />
              </IconButton>
            </div>
          </>
        ) : (
          <ProjectValues>
            <Skeleton width={200} height={28} />
            <Skeleton width={200} height={28} />
          </ProjectValues>
        )}
      </ProjectValuesWrapper>
      {currentUsername && (
        <StaticAdminDialog
          open={isEdit}
          currentUsername={currentUsername}
          onClose={handleClose}
          onSubmit={handleSubmit}
        />
      )}
    </>
  );
});
