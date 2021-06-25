import {
  Box,
  IconButton,
  makeStyles,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
  Typography,
} from "@material-ui/core";
import { MoreVert as MoreVertIcon } from "@material-ui/icons";
import { EntityId } from "@reduxjs/toolkit";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { CopyIconButton } from "~/components/copy-icon-button";
import { UI_TEXT_DELETE, UI_TEXT_EDIT } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { setTargetEnv } from "~/modules/deleting-env";
import { selectEnvById } from "~/modules/environments";
import { EditEnvironmentDialog } from "./edit-dialog";

const useStyles = makeStyles((theme) => ({
  item: {
    backgroundColor: theme.palette.background.paper,
  },
  idCell: {
    "& button": {
      visibility: "hidden",
    },
    "&:hover button": {
      visibility: "visible",
    },
  },
}));

const ITEM_HEIGHT = 48;
const TEXT_NO_DESCRIPTION = "No description";

export interface EnvironmentListItemProps {
  id: EntityId;
}

export const EnvironmentListItem: FC<EnvironmentListItemProps> = memo(
  function EnvironmentListItem({ id }) {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const [isEdit, setIsEdit] = useState(false);
    const env = useAppSelector(selectEnvById(id));
    const dispatch = useAppDispatch();

    // menu event handler
    const handleClickMenu = useCallback(
      (e: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(e.currentTarget);
      },
      []
    );
    const closeMenu = useCallback(() => {
      setAnchorEl(null);
    }, []);

    // edit event handler
    const handleEdit = useCallback(() => {
      setIsEdit(true);
      setAnchorEl(null);
    }, [setIsEdit, setAnchorEl]);
    const handleCloseEdit = useCallback(() => {
      setIsEdit(false);
    }, [setIsEdit]);

    const handleDeleteClick = useCallback(() => {
      closeMenu();
      if (env) {
        dispatch(setTargetEnv({ id: env.id, name: env.name }));
      }
    }, [dispatch, env, closeMenu]);

    if (!env) {
      return null;
    }

    const isOpenMenu = Boolean(anchorEl);

    return (
      <>
        <TableRow key={`env-${env.id}`} className={classes.item}>
          <TableCell>
            <Typography variant="subtitle2" component="span">
              {env.name}
            </Typography>
          </TableCell>
          <TableCell colSpan={2}>{env.desc || TEXT_NO_DESCRIPTION}</TableCell>
          <TableCell className={classes.idCell}>
            <Box display="flex" alignItems="center" fontFamily="fontFamilyMono">
              {env.id}
              <CopyIconButton name="Environment ID" value={env.id} />
            </Box>
          </TableCell>
          <TableCell align="right" style={{ height: 61 }}>
            <IconButton
              edge="end"
              aria-label="open menu"
              onClick={handleClickMenu}
            >
              <MoreVertIcon />
            </IconButton>
          </TableCell>
        </TableRow>

        <Menu
          id="env-menu"
          anchorEl={anchorEl}
          open={isOpenMenu}
          onClose={closeMenu}
          PaperProps={{
            style: {
              maxHeight: ITEM_HEIGHT * 4.5,
              width: "20ch",
            },
          }}
        >
          <MenuItem
            key="env-menu-edit"
            onClick={handleEdit}
            style={{ display: "none" }}
          >
            {UI_TEXT_EDIT}
          </MenuItem>
          <MenuItem key="env-menu-delete" onClick={handleDeleteClick}>
            {UI_TEXT_DELETE}
          </MenuItem>
        </Menu>

        <EditEnvironmentDialog
          description={env.desc}
          open={isEdit}
          onClose={handleCloseEdit}
        />
      </>
    );
  }
);
