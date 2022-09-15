import {
  Box,
  Button,
  Menu,
  MenuItem,
  IconButton,
  Divider,
  makeStyles,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
} from "@material-ui/core";
import * as React from "react";
import { Add as AddIcon, MoreVert as MenuIcon } from "@material-ui/icons";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { UI_TEXT_ADD } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { fetchProject, addUserGroup, deleteUserGroup } from "~/modules/project";
import { AddUserGroupDialog } from "./components/add-user-group-dialog";
import { addToast } from "~/modules/toasts";
import {
  ADD_USER_GROUP_SUCCESS,
  DELETE_USER_GROUP_SUCCESS,
} from "~/constants/toast-text";

const useStyles = makeStyles(() => ({
  toolbarSpacer: {
    flexGrow: 1,
  },
}));

export const SettingsUserGroupPage: FC = memo(function SettingsUserGroupPage() {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const userGroups = useAppSelector((state) => state.project.userGroups);
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );

  const handleSubmit = useCallback(
    (values: { ssoGroup: string; role: string }) => {
      dispatch(addUserGroup(values)).then(() => {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: ADD_USER_GROUP_SUCCESS,
            severity: "success",
          })
        );
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

  const handleDelete = useCallback(
    (ssoGroup: string) => {
      setAnchorEl(null);
      dispatch(deleteUserGroup({ ssoGroup: ssoGroup })).then(() => {
        dispatch(fetchProject());
        dispatch(
          addToast({
            message: DELETE_USER_GROUP_SUCCESS,
            severity: "success",
          })
        );
      });
    },
    [dispatch]
  );

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenAddForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
        <div className={classes.toolbarSpacer} />
      </Toolbar>
      <Divider />

      <Box display="flex" flex={1} overflow="hidden">
        <TableContainer component={Paper} square>
          <Table aria-label="user group list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell colSpan={2}>Group</TableCell>
                <TableCell>Role</TableCell>
                <TableCell align="right" />
              </TableRow>
            </TableHead>
            <TableBody>
              {userGroups.map((group) => (
                <TableRow key={group.ssoGroup}>
                  <TableCell colSpan={2}>{group.ssoGroup}</TableCell>
                  <TableCell colSpan={2}>{group.role}</TableCell>
                  <TableCell align="right">
                    <IconButton
                      data-id={group.ssoGroup}
                      onClick={handleOpenMenu}
                    >
                      <MenuIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>

      <Menu
        id="user-group-menu"
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        onClose={handleCloseMenu}
      >
        <MenuItem
          onClick={() => {
            if (anchorEl && anchorEl.dataset.id) {
              handleDelete(anchorEl.dataset.id);
            }
          }}
        >
          Delete User Group
        </MenuItem>
      </Menu>

      <AddUserGroupDialog
        open={isOpenAddForm}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />
    </>
  );
});
