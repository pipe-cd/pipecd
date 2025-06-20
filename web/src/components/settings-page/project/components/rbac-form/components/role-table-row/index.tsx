import { IconButton, Menu, MenuItem, TableCell, TableRow } from "@mui/material";
import { MoreVert as MoreVertIcon } from "@mui/icons-material";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { UI_TEXT_EDIT, UI_TEXT_DELETE } from "~/constants/ui-text";
import { formalizePoliciesList } from "~/utils/formalize-policies-list";
import { ProjectRBACRole } from "pipecd/web/model/project_pb";

interface Props {
  role: ProjectRBACRole.AsObject;
  onEdit: (role: ProjectRBACRole.AsObject) => void;
  onDelete: (roleName: string) => void;
}

const ITEM_HEIGHT = 48;

export const RoleTableRow: FC<Props> = memo(function RoleTableRow({
  role,
  onDelete,
  onEdit,
}) {
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);

  const handleMenuOpen = useCallback(
    (event: React.MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(event.currentTarget);
    },
    []
  );

  const handleMenuClose = useCallback(() => {
    setAnchorEl(null);
  }, []);

  const handleEdit = useCallback(() => {
    setAnchorEl(null);
    onEdit(role);
  }, [role, onEdit]);

  const handleDelete = useCallback(() => {
    setAnchorEl(null);
    onDelete(role.name);
  }, [role, onDelete]);

  return (
    <>
      <TableRow key={role.name}>
        <TableCell>{role.name}</TableCell>
        <TableCell>
          {formalizePoliciesList({
            policiesList: role.policiesList,
          })
            .split("\n\n")
            .map((policy, i) => (
              <p key={i}>{policy}</p>
            ))}
        </TableCell>
        <TableCell align="right">
          <IconButton
            edge="end"
            aria-label="open menu"
            onClick={handleMenuOpen}
            disabled={role.isBuiltin}
            size="large"
          >
            <MoreVertIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      <Menu
        id="role-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        slotProps={{
          paper: {
            style: {
              maxHeight: ITEM_HEIGHT * 4.5,
              width: "15ch",
            },
          },
        }}
      >
        <MenuItem key="role-menu-edit" onClick={handleEdit}>
          {UI_TEXT_EDIT}
        </MenuItem>
        <MenuItem key="role-menu-delete" onClick={handleDelete}>
          {UI_TEXT_DELETE}
        </MenuItem>
      </Menu>
    </>
  );
});
