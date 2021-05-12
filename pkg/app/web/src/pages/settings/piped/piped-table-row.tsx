import {
  Box,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
  Typography,
  DialogActions,
  Button,
  DialogContentText,
} from "@material-ui/core";
import {
  FileCopyOutlined as CopyIcon,
  MoreVert as MoreVertIcon,
} from "@material-ui/icons";
import clsx from "clsx";
import copy from "copy-to-clipboard";
import dayjs from "dayjs";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppDispatch } from "../../../store";
import {
  COPY_PIPED_ID,
  DELETE_OLD_PIPED_KEY_SUCCESS,
} from "../../../constants/toast-text";
import {
  UI_TEXT_ADD_NEW_KEY,
  UI_TEXT_DELETE_OLD_KEY,
  UI_TEXT_DISABLE,
  UI_TEXT_EDIT,
  UI_TEXT_ENABLE,
} from "../../../constants/ui-text";
import {
  addNewPipedKey,
  deleteOldKey,
  fetchPipeds,
  selectPipedById,
} from "../../../modules/pipeds";
import { addToast } from "../../../modules/toasts";

const useStyles = makeStyles((theme) => ({
  disabledItem: {
    background: theme.palette.grey[200],
  },
  copyButton: {
    marginLeft: theme.spacing(1),
    visibility: "hidden",
    "tr:hover &": {
      visibility: "visible",
    },
  },
}));

interface Props {
  pipedId: string;
  onEdit: (id: string) => void;
  onDisable: (id: string) => void;
  onEnable: (id: string) => void;
}

const ITEM_HEIGHT = 48;
const menuStyle = {
  style: {
    maxHeight: ITEM_HEIGHT * 4.5,
    width: "20ch",
  },
};

export const PipedTableRow: FC<Props> = memo(function PipedTableRow({
  pipedId,
  onEnable,
  onDisable,
  onEdit,
}) {
  const classes = useStyles();
  const piped = useSelector(selectPipedById(pipedId));
  const dispatch = useDispatch<AppDispatch>();
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const hasOldKey = piped ? piped.keysList.length > 1 : false;
  const [openOldKeyAlert, setOpenOldKeyAlert] = useState(false);

  const handleMenuOpen = useCallback(
    (event: React.MouseEvent<HTMLButtonElement>) => {
      setAnchorEl(event.currentTarget);
    },
    []
  );

  const handleMenuClose = useCallback(() => {
    setAnchorEl(null);
  }, []);

  const handleAlertClose = useCallback(() => {
    setOpenOldKeyAlert(false);
  }, []);

  const handleEdit = useCallback(() => {
    setAnchorEl(null);
    onEdit(pipedId);
  }, [pipedId, onEdit]);

  const handleAddNewKey = useCallback(() => {
    setAnchorEl(null);
    if (hasOldKey) {
      setOpenOldKeyAlert(true);
    } else {
      dispatch(addNewPipedKey({ pipedId }));
    }
  }, [dispatch, pipedId, hasOldKey]);

  const handleDeleteOldKey = useCallback(() => {
    setAnchorEl(null);
    dispatch(deleteOldKey({ pipedId })).then(() => {
      dispatch(fetchPipeds(true));
      dispatch(
        addToast({
          message: DELETE_OLD_PIPED_KEY_SUCCESS,
          severity: "success",
        })
      );
    });
  }, [pipedId, dispatch]);

  const handleEnable = useCallback(() => {
    setAnchorEl(null);
    onEnable(pipedId);
  }, [pipedId, onEnable]);

  const handleDisable = useCallback(() => {
    setAnchorEl(null);
    onDisable(pipedId);
  }, [pipedId, onDisable]);

  const handleCopy = useCallback(() => {
    if (piped) {
      copy(piped.id);
      dispatch(addToast({ message: COPY_PIPED_ID }));
    }
  }, [dispatch, piped]);

  if (!piped) {
    return null;
  }

  return (
    <>
      <TableRow
        key={`pipe-${piped.id}`}
        className={clsx({ [classes.disabledItem]: piped.disabled })}
      >
        <TableCell>
          <Typography variant="subtitle2">{piped.name}</Typography>
        </TableCell>
        <TableCell title={piped.id}>
          <Box display="flex" alignItems="center" fontFamily="fontFamilyMono">
            {piped.id}
            <IconButton
              className={classes.copyButton}
              aria-label="Copy piped id"
              onClick={handleCopy}
            >
              <CopyIcon />
            </IconButton>
          </Box>
        </TableCell>
        <TableCell>{piped.version}</TableCell>
        <TableCell>
          <Typography variant="body2" color="textSecondary">
            {piped.desc}
          </Typography>
        </TableCell>
        <TableCell>
          {piped.startedAt === 0
            ? "Not Yet Started"
            : dayjs(piped.startedAt * 1000).fromNow()}
        </TableCell>
        <TableCell align="right">
          <IconButton
            edge="end"
            aria-label="open menu"
            onClick={handleMenuOpen}
          >
            <MoreVertIcon />
          </IconButton>
        </TableCell>
      </TableRow>

      <Menu
        id="piped-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        PaperProps={menuStyle}
      >
        {piped.disabled ? (
          <MenuItem onClick={handleEnable}>{UI_TEXT_ENABLE}</MenuItem>
        ) : (
          [
            <MenuItem key="piped-menu-edit" onClick={handleEdit}>
              {UI_TEXT_EDIT}
            </MenuItem>,
            <MenuItem key="piped-menu-add-new-key" onClick={handleAddNewKey}>
              {UI_TEXT_ADD_NEW_KEY}
            </MenuItem>,
            <MenuItem
              disabled={hasOldKey === false}
              key="piped-menu-delete-old-key"
              onClick={handleDeleteOldKey}
            >
              {UI_TEXT_DELETE_OLD_KEY}
            </MenuItem>,
            <MenuItem key="piped-menu-disable" onClick={handleDisable}>
              {UI_TEXT_DISABLE}
            </MenuItem>,
          ]
        )}
      </Menu>

      <Dialog open={openOldKeyAlert} onClose={handleAlertClose}>
        <DialogTitle>There are already 2 keys for this piped</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Before adding a new key, you need to delete the old one. <br />
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button color="primary" autoFocus onClick={handleAlertClose}>
            OK
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
});
