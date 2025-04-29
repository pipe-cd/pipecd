import { IconButton, Menu, MenuItem, TableCell, TableRow } from "@mui/material";
import { MoreVert as MoreVertIcon } from "@mui/icons-material";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { formalizePoliciesList } from "~/modules/project";
import { UI_TEXT_EDIT, UI_TEXT_DELETE } from "~/constants/ui-text";
import { useAppSelector } from "~/hooks/redux";

interface Props {
  role: string;
  onEdit: (role: string) => void;
  onDelete: (role: string) => void;
}

const ITEM_HEIGHT = 48;
const menuStyle = {
  style: {
    maxHeight: ITEM_HEIGHT * 4.5,
    width: "15ch",
  },
};

export const RoleTableRow: FC<Props> = memo(function RoleTableRow({
  role,
  onDelete,
  onEdit,
}) {
  const rs = useAppSelector((state) => state.project.rbacRoles);
  const r = rs.filter((r) => r.name == role)[0];
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
    onDelete(role);
  }, [role, onDelete]);

  if (!r) {
    return null;
  }

  return (
    <>
      <TableRow key={role}>
        <TableCell>{r.name}</TableCell>
        <TableCell>
          {formalizePoliciesList({
            policiesList: r.policiesList,
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
            disabled={r.isBuiltin}
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
        PaperProps={menuStyle}
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
