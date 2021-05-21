import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
  TextField,
  Typography,
} from "@material-ui/core";
import { MoreVert as MoreVertIcon } from "@material-ui/icons";
import { EntityId } from "@reduxjs/toolkit";
import { FC, memo, useCallback, useState } from "react";
import * as React from "react";
import {
  UI_TEXT_CANCEL,
  UI_TEXT_EDIT,
  UI_TEXT_SAVE,
} from "../../constants/ui-text";
import { selectEnvById } from "../../modules/environments";
import { CopyIconButton } from "../copy-icon-button";
import { useAppSelector } from "../../hooks/redux";

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
const DIALOG_TITLE = "Edit Environment description";

export interface EnvironmentListItemProps {
  id: EntityId;
}

export const EnvironmentListItem: FC<EnvironmentListItemProps> = memo(
  function EnvironmentListItem({ id }) {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const [isEdit, setIsEdit] = useState(false);
    const [desc, setDesc] = useState("");
    const env = useAppSelector(selectEnvById(id));

    // menu event handler
    const handleClickMenu = useCallback(
      (e: React.MouseEvent<HTMLButtonElement>) => {
        setAnchorEl(e.currentTarget);
      },
      [setAnchorEl]
    );
    const handleCloseMenu = useCallback(() => {
      setAnchorEl(null);
    }, [setAnchorEl]);

    // edit event handler
    const handleEdit = useCallback(() => {
      setIsEdit(true);
      setAnchorEl(null);
    }, [setIsEdit, setAnchorEl]);
    const handleCloseEdit = useCallback(() => {
      setIsEdit(false);
    }, [setIsEdit]);
    const handleSave = useCallback(() => {
      // not implemented yet
    }, []);

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
              style={{ display: "none" }}
            >
              <MoreVertIcon />
            </IconButton>
          </TableCell>
        </TableRow>

        <Menu
          id="env-menu"
          anchorEl={anchorEl}
          open={isOpenMenu}
          onClose={handleCloseMenu}
          PaperProps={{
            style: {
              maxHeight: ITEM_HEIGHT * 4.5,
              width: "20ch",
            },
          }}
        >
          <MenuItem key="env-menu-edit" onClick={handleEdit}>
            {UI_TEXT_EDIT}
          </MenuItem>
        </Menu>

        <Dialog
          open={isEdit}
          onEnter={() => {
            setDesc(env.desc);
          }}
          onClose={handleCloseEdit}
          fullWidth
        >
          <form onSubmit={handleSave}>
            <DialogTitle>{DIALOG_TITLE}</DialogTitle>
            <DialogContent>
              <TextField
                value={desc}
                variant="outlined"
                margin="dense"
                label="Description"
                fullWidth
                required
                autoFocus
                onChange={(e) => setDesc(e.currentTarget.value)}
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseEdit}>{UI_TEXT_CANCEL}</Button>
              <Button
                type="submit"
                color="primary"
                disabled={desc === "" || desc === env.desc}
              >
                {UI_TEXT_SAVE}
              </Button>
            </DialogActions>
          </form>
        </Dialog>
      </>
    );
  }
);
