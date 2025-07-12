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
  Box,
} from "@mui/material";
import * as React from "react";
import { Add as AddIcon, MoreVert as MenuIcon } from "@mui/icons-material";
import { FC, memo, useCallback, useState } from "react";
import { UI_TEXT_ADD } from "~/constants/ui-text";
import { AddUserGroupDialog } from "../add-user-group-dialog";
import { DeleteUserGroupConfirmDialog } from "../delete-user-group-confirm-dialog";
import { ProjectTitle } from "~/styles/project-setting";
import {
  ADD_USER_GROUP_SUCCESS,
  DELETE_USER_GROUP_SUCCESS,
} from "~/constants/toast-text";
import { useToast } from "~/contexts/toast-context";
import { useGetProject } from "~/queries/project/use-get-project";
import { useAddUserGroup } from "~/queries/project/use-add-user-group";
import { useDeleteUserGroup } from "~/queries/project/use-delete-user-group";

const SUB_SECTION_TITLE = "User Group";

export const UserGroupTable: FC = memo(function UserGroupTable() {
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [deleteSSOGroup, setDeleteSSOGroup] = useState<null | string>(null);
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );
  const { addToast } = useToast();
  const { data: projectDetail } = useGetProject();
  const { mutateAsync: addUserGroup } = useAddUserGroup();
  const { mutateAsync: deleteUserGroup } = useDeleteUserGroup();
  const userGroups = projectDetail?.userGroups || [];

  const handleSubmit = useCallback(
    (values: { ssoGroup: string; role: string }) => {
      addUserGroup(values).then(() => {
        addToast({
          message: ADD_USER_GROUP_SUCCESS,
          severity: "success",
        });
      });
    },
    [addToast, addUserGroup]
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

  const handleCancelDeleting = useCallback(() => {
    setDeleteSSOGroup(null);
  }, [setDeleteSSOGroup]);

  const handleDelete = useCallback(
    (ssoGroup: string) => {
      deleteUserGroup({ ssoGroup: ssoGroup }).then(() => {
        addToast({
          message: DELETE_USER_GROUP_SUCCESS,
          severity: "success",
        });
      });
      setDeleteSSOGroup(null);
    },
    [addToast, deleteUserGroup]
  );

  return (
    <>
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between",
          paddingTop: 2,
        }}
      >
        <ProjectTitle variant="h6">{SUB_SECTION_TITLE}</ProjectTitle>

        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenAddForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
      </Box>
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
        slotProps={{
          paper: {
            sx: { width: "15ch" },
          },
        }}
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
