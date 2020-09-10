import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  ListItem,
  ListItemSecondaryAction,
  ListItemText,
  makeStyles,
  Menu,
  MenuItem,
  TextField,
} from "@material-ui/core";
import { MoreVert as MoreVertIcon } from "@material-ui/icons";
import { EntityId } from "@reduxjs/toolkit";
import React, { FC, memo, useCallback, useState } from "react";
import { useSelector } from "react-redux";
import {
  UI_TEXT_CANCEL,
  UI_TEXT_EDIT,
  UI_TEXT_SAVE,
} from "../constants/ui-text";
import { AppState } from "../modules";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";

const useStyles = makeStyles((theme) => ({
  item: {
    backgroundColor: theme.palette.background.paper,
  },
}));

const ITEM_HEIGHT = 48;
const TEXT_NO_DESCRIPTION = "No description";
const DIALOG_TITLE = "Edit Environment description";

interface Props {
  id: EntityId;
}

export const EnvironmentListItem: FC<Props> = memo(
  function EnvironmentListItem({ id }) {
    const classes = useStyles();
    const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
    const [isEdit, setIsEdit] = useState(false);
    const [desc, setDesc] = useState("");
    const env = useSelector<AppState, Environment | undefined>((state) =>
      selectEnvById(state.environments, id)
    );

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
        <ListItem key={`env-${env.id}`} divider dense className={classes.item}>
          <ListItemText
            primary={env.name}
            secondary={env.desc || TEXT_NO_DESCRIPTION}
          />
          {/** TODO: Remove this style after implemented editing environment's desc API */}
          <ListItemSecondaryAction style={{ display: "none" }}>
            <IconButton
              edge="end"
              aria-label="open menu"
              onClick={handleClickMenu}
            >
              <MoreVertIcon />
            </IconButton>
          </ListItemSecondaryAction>
        </ListItem>

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
