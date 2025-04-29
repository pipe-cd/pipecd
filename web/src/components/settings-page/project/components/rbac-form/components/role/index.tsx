import {
  Button,
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
import { useProjectSettingStyles } from "~/styles/project-setting";
import { Add as AddIcon } from "@mui/icons-material";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  fetchProject,
  addRBACRole,
  deleteRBACRole,
  updateRBACRole,
} from "~/modules/project";
import { AddRoleDialog } from "../add-role-dialog";
import { DeleteRoleConfirmDialog } from "../delete-role-confirm-dialog";
import { RoleTableRow } from "../role-table-row";
import { EditRoleDialog } from "../edit-role-dialog";
import { FC, memo, useCallback, useState } from "react";
import { UI_TEXT_ADD } from "~/constants/ui-text";
import { addToast } from "~/modules/toasts";
import { parseRBACPolicies } from "~/modules/project";
import {
  ADD_RBAC_ROLE_SUCCESS,
  DELETE_RBAC_ROLE_SUCCESS,
  UPDATE_RBAC_ROLE_SUCCESS,
} from "~/constants/toast-text";

const useStyles = makeStyles((theme) => ({
  title: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    paddingTop: theme.spacing(2),
  },
}));

const SUB_SECTION_TITLE = "Role";

export const RoleTable: FC = memo(function RoleTable() {
  const classes = useStyles();
  const projectSettingClasses = useProjectSettingStyles();
  const rbacRoles = useAppSelector((state) => state.project.rbacRoles);
  const dispatch = useAppDispatch();
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [deleteRole, setDeleteRole] = useState<null | string>(null);
  const [editRole, setEditRole] = useState<null | string>(null);

  const handleSubmit = useCallback(
    (values: { name: string; policies: string }) => {
      const params = {
        name: values.name,
        policies: parseRBACPolicies({ policies: values.policies }),
      };
      dispatch(addRBACRole(params)).then((result) => {
        if (addRBACRole.fulfilled.match(result)) {
          dispatch(fetchProject());
          dispatch(
            addToast({
              message: ADD_RBAC_ROLE_SUCCESS,
              severity: "success",
            })
          );
        }
      });
    },
    [dispatch]
  );

  const handleDelete = useCallback(
    (role: string) => {
      dispatch(deleteRBACRole({ name: role })).then((result) => {
        if (deleteRBACRole.fulfilled.match(result)) {
          dispatch(fetchProject());
          dispatch(
            addToast({
              message: DELETE_RBAC_ROLE_SUCCESS,
              severity: "success",
            })
          );
        }
      });
      setDeleteRole(null);
    },
    [dispatch]
  );

  const handleDeleteClose = useCallback(() => {
    setDeleteRole(null);
  }, [setDeleteRole]);

  const handleUpdate = useCallback(
    (values: { name: string; policies: string }) => {
      const params = {
        name: values.name,
        policies: parseRBACPolicies({ policies: values.policies }),
      };
      dispatch(updateRBACRole(params)).then((result) => {
        if (updateRBACRole.fulfilled.match(result)) {
          dispatch(fetchProject());
          dispatch(
            addToast({
              message: UPDATE_RBAC_ROLE_SUCCESS,
              severity: "success",
            })
          );
        }
      });
      setEditRole(null);
    },
    [dispatch]
  );

  const handleEditClose = useCallback(() => {
    setEditRole(null);
  }, []);

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
              <TableCell>Role</TableCell>
              <TableCell>Policies</TableCell>
              <TableCell align="right" />
            </TableRow>
          </TableHead>
          <TableBody>
            {rbacRoles.map((role) => (
              <RoleTableRow
                key={role.name}
                role={role.name}
                onEdit={(role) => setEditRole(role)}
                onDelete={(role) => setDeleteRole(role)}
              />
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <AddRoleDialog
        open={isOpenAddForm}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />

      <DeleteRoleConfirmDialog
        role={deleteRole}
        onClose={handleDeleteClose}
        onDelete={handleDelete}
      />

      <EditRoleDialog
        role={editRole}
        onClose={handleEditClose}
        onUpdate={handleUpdate}
      />
    </>
  );
});
