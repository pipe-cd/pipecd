import {
  Button,
  Menu,
  MenuItem,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import * as React from "react";
import { Add as AddIcon, MoreVert as MenuIcon } from "@mui/icons-material";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { UI_TEXT_ADD } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchProject, addUserGroup, deleteUserGroup } from "~/modules/project";
import { AddUserGroupDialog } from "../add-user-group-dialog";
import { DeleteUserGroupConfirmDialog } from "../delete-user-group-confirm-dialog";
import { useProjectSettingStyles } from "~/styles/project-setting";
import { addToast } from "~/modules/toasts";
import {
  ADD_USER_GROUP_SUCCESS,
  DELETE_USER_GROUP_SUCCESS,
} from "~/constants/toast-text";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    paddingTop: theme.spacing(2),
  },
}));

const menuStyle = {
  style: {
    width: "15ch",
  },
};

const SUB_SECTION_TITLE = "User Group";

export const UserGroupTable: FC = memo(function UserGroupTable() {
  const classes = useStyles();
  const projectSettingClasses = useProjectSettingStyles();
  const dispatch = useAppDispatch();
  const userGroups = useAppSelector((state) => state.project.userGroups);
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [deleteSSOGroup, setDeleteSSOGroup] = useState<null | string>(null);
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );

  const handleSubmit = useCallback(
    (values: { ssoGroup: string; role: string }) => {
      dispatch(addUserGroup(values)).then((result) => {
        if (addUserGroup.fulfilled.match(result)) {
          dispatch(fetchProject());
          dispatch(
            addToast({
              message: ADD_USER_GROUP_SUCCESS,
              severity: "success",
            })
          );
        }
      });
    },
    [dispatch]
  );

  const handleOpenMenu = useCallback(
    (e: React.MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(e.currentTarget);
    },
    [setAnchorEl]
  );

  const handleCloseMenu = useCallback(() => {
    setAnchorEl(null);
  }, [setAnchorEl]);

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  const handleCancelDeleting = useCallback(() => {
    setDeleteSSOGroup(null);
  }, [setDeleteSSOGroup]);

  const handleDelete = useCallback(
    (ssoGroup: string) => {
      dispatch(deleteUserGroup({ ssoGroup: ssoGroup })).then((result) => {
        if (deleteUserGroup.fulfilled.match(result)) {
          dispatch(fetchProject());
          dispatch(
            addToast({
              message: DELETE_USER_GROUP_SUCCESS,
              severity: "success",
            })
          );
        }
      });
      setDeleteSSOGroup(null);
    },
    [dispatch]
  );

  return (
    <>
      <div className={classes.title}>
        <Typography
          variant="h6"
          className={projectSettingClasses.titleWithIcon}
        >
          {SUB_SECTION_TITLE}
        </Typography>

        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenAddForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
      </div>
      <TableContainer component={Paper} square>
        <Table size="small" stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell>Team/Group</TableCell>
              <TableCell>Role</TableCell>
              <TableCell align="right" />
            </TableRow>
          </TableHead>
          <TableBody>
            {userGroups.map((group) => (
              <TableRow key={group.ssoGroup}>
                <TableCell>{group.ssoGroup}</TableCell>
                <TableCell>{group.role}</TableCell>
                <TableCell align="right">
                  <IconButton
                    data-id={group.ssoGroup}
                    onClick={handleOpenMenu}
                    size="large"
                  >
                    <MenuIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
      <Menu
        id="user-group-menu"
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        onClose={handleCloseMenu}
        PaperProps={menuStyle}
      >
        <MenuItem
          onClick={() => {
            if (anchorEl && anchorEl.dataset.id) {
              setDeleteSSOGroup(anchorEl.dataset.id);
            }
            setAnchorEl(null);
          }}
        >
          Delete
        </MenuItem>
      </Menu>
      <AddUserGroupDialog
        open={isOpenAddForm}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />
      <DeleteUserGroupConfirmDialog
        ssoGroup={deleteSSOGroup}
        onCancel={handleCancelDeleting}
        onDelete={handleDelete}
      />
    </>
  );
});
