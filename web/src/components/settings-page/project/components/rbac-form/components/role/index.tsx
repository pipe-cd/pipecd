import {
  Box,
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import { ProjectTitle } from "~/styles/project-setting";
import { Add as AddIcon } from "@mui/icons-material";
import { AddRoleDialog } from "../add-role-dialog";
import { DeleteRoleConfirmDialog } from "../delete-role-confirm-dialog";
import { RoleTableRow } from "../role-table-row";
import { EditRoleDialog } from "../edit-role-dialog";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { UI_TEXT_ADD } from "~/constants/ui-text";

import {
  ADD_RBAC_ROLE_SUCCESS,
  DELETE_RBAC_ROLE_SUCCESS,
  UPDATE_RBAC_ROLE_SUCCESS,
} from "~/constants/toast-text";
import { useGetProject } from "~/queries/project/use-get-project";
import { useAddProjectRBACRole } from "~/queries/project/use-add-project-rbac-role";
import { useToast } from "~/contexts/toast-context";
import { useDeleteProjectRBACRole } from "~/queries/project/use-delete-project-rbac-role";
import { useUpdateProjectRBACRole } from "~/queries/project/use-update-project-rbac-role";
import { parseRBACPolicies } from "~/utils/parse-rbac-policies";
import { ProjectRBACRole } from "pipecd/web/model/project_pb";

interface RoleTableProps {
  isProjectDisabled: boolean;
}

const SUB_SECTION_TITLE = "Role";

export const RoleTable: FC<RoleTableProps> = memo(function RoleTable({
  isProjectDisabled,
}) {
  const { data: projectDetail } = useGetProject();
  const rbacRoles = projectDetail?.rbacRoles || [];

  const { mutateAsync: addRBACRole } = useAddProjectRBACRole();
  const { mutateAsync: deleteRBACRole } = useDeleteProjectRBACRole();
  const { mutateAsync: updateRBACRole } = useUpdateProjectRBACRole();
  const { addToast } = useToast();

  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [deleteRoleName, setDeleteRoleName] = useState<null | string>(null);
  const [editRole, setEditRole] = useState<null | ProjectRBACRole.AsObject>(
    null
  );

  useEffect(() => {
    if (isProjectDisabled) {
      setIsOpenAddForm(false);
      setDeleteRoleName(null);
      setEditRole(null);
    }
  }, [isProjectDisabled]);

  const handleSubmit = useCallback(
    (values: { name: string; policies: string }) => {
      const params = {
        name: values.name,
        policies: parseRBACPolicies({ policies: values.policies }),
      };
      addRBACRole(params).then(() => {
        addToast({
          message: ADD_RBAC_ROLE_SUCCESS,
          severity: "success",
        });
      });
    },
    [addRBACRole, addToast]
  );

  const handleDelete = useCallback(
    (roleName: string) => {
      deleteRBACRole({ name: roleName }).then(() => {
        addToast({
          message: DELETE_RBAC_ROLE_SUCCESS,
          severity: "success",
        });
      });
      setDeleteRoleName(null);
    },
    [addToast, deleteRBACRole]
  );

  const handleDeleteClose = useCallback(() => {
    setDeleteRoleName(null);
  }, [setDeleteRoleName]);

  const handleUpdate = useCallback(
    (values: { name: string; policies: string }) => {
      const params = {
        name: values.name,
        policies: parseRBACPolicies({ policies: values.policies }),
      };
      updateRBACRole(params).then(() => {
        addToast({
          message: UPDATE_RBAC_ROLE_SUCCESS,
          severity: "success",
        });
      });
      setEditRole(null);
    },
    [addToast, updateRBACRole]
  );

  const handleEditClose = useCallback(() => {
    setEditRole(null);
  }, []);

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
          disabled={isProjectDisabled}
        >
          {UI_TEXT_ADD}
        </Button>
      </Box>

      <TableContainer
        component={Paper}
        square
        sx={{ opacity: isProjectDisabled ? 0.6 : 1 }}
      >
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
                role={role}
                onEdit={(role) => setEditRole(role)}
                onDelete={(role) => setDeleteRoleName(role)}
                disabled={isProjectDisabled}
              />
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <AddRoleDialog
        open={isOpenAddForm && !isProjectDisabled}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />

      <DeleteRoleConfirmDialog
        roleName={deleteRoleName}
        onClose={handleDeleteClose}
        onDelete={handleDelete}
      />

      <EditRoleDialog
        role={isProjectDisabled ? null : editRole}
        onClose={handleEditClose}
        onUpdate={handleUpdate}
      />
    </>
  );
});
