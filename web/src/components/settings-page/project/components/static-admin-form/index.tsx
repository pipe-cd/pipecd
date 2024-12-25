import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  Switch,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import Skeleton from "@material-ui/lab/Skeleton/Skeleton";
import clsx from "clsx";
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
import { useProjectSettingStyles } from "~/styles/project-setting";
import { ProjectSettingLabeledText } from "../project-setting-labeled-text";

const useStyles = makeStyles(() => ({
  disabled: {
    opacity: 0.5,
  },
}));

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
  const classes = useStyles();
  const projectSettingClasses = useProjectSettingStyles();
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
      <div className={projectSettingClasses.title}>
        <Typography
          variant="h5"
          className={projectSettingClasses.titleWithIcon}
        >
          {SECTION_TITLE}
        </Typography>

        <Switch
          checked={isEnabled}
          color="primary"
          onClick={handleToggleAvailability}
          disabled={currentUsername === null}
        />
      </div>

      <Typography
        variant="body1"
        color="textSecondary"
        className={projectSettingClasses.description}
      >
        {STATIC_ADMIN_DESCRIPTION}
      </Typography>

      <div
        className={clsx(projectSettingClasses.valuesWrapper, {
          [classes.disabled]: isEnabled === false,
        })}
      >
        {currentUsername ? (
          <>
            <div className={projectSettingClasses.values}>
              <ProjectSettingLabeledText
                label="Username"
                value={currentUsername}
              />
              <ProjectSettingLabeledText label="Password" value="********" />
            </div>
            <div>
              <IconButton
                aria-label="edit static admin user"
                onClick={() => setIsEdit(true)}
                disabled={isEnabled === false}
              >
                <EditIcon />
              </IconButton>
            </div>
          </>
        ) : (
          <div className={projectSettingClasses.values}>
            <Skeleton width={200} height={28} />
            <Skeleton width={200} height={28} />
          </div>
        )}
      </div>
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
