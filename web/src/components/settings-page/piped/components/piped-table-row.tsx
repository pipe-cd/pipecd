import {
  Box,
  Button,
  Badge,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Menu,
  MenuItem,
  TableCell,
  TableRow,
  Typography,
  Tooltip,
} from "@mui/material";
import { MoreVert as MoreVertIcon } from "@mui/icons-material";
import dayjs from "dayjs";
import { Highlight, themes } from "prism-react-renderer";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { CopyIconButton } from "~/components/copy-icon-button";
import DialogConfirm from "~/components/dialog-confirm";
import { PIPED_CONNECTION_STATUS_TEXT } from "~/constants/piped-connection-status-text";
import { DELETE_OLD_PIPED_KEY_SUCCESS } from "~/constants/toast-text";
import {
  UI_TEXT_ADD_NEW_KEY,
  UI_TEXT_DELETE_OLD_KEY,
  UI_TEXT_VIEW_THE_CONFIGURATION,
  UI_TEXT_DISABLE,
  UI_TEXT_EDIT,
  UI_TEXT_ENABLE,
  UI_TEXT_RESTART,
  UI_TEXT_CANCEL,
} from "~/constants/ui-text";
import { Piped } from "pipecd/web/model/piped_pb";
import { useToast } from "~/contexts/toast-context";
import { useDeleteOldPipedKey } from "~/queries/pipeds/use-delete-old-piped-key";

interface Props {
  piped: Piped.AsObject;
  onEdit: (piped: Piped.AsObject) => void;
  onAddNewKey: (piped: Piped.AsObject) => void;
  onDisable: (id: string) => void;
  onEnable: (id: string) => void;
  onRestart: (id: string) => void;
}

const ITEM_HEIGHT = 48;

export const PipedTableRow: FC<Props> = memo(function PipedTableRow({
  piped,
  onEnable,
  onDisable,
  onEdit,
  onRestart,
  onAddNewKey,
}) {
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);
  const hasOldKey = piped ? piped.keysList.length > 1 : false;
  const [openOldKeyAlert, setOpenOldKeyAlert] = useState(false);
  const [openConfigAlert, setOpenConfigAlert] = useState(false);
  const [openConfirmAddKey, setOpenConfirmAddKey] = useState(false);

  const { addToast } = useToast();

  const { mutateAsync: deleteOldPipedKey } = useDeleteOldPipedKey();

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
    setOpenConfigAlert(false);
  }, []);

  const handleEdit = useCallback(() => {
    setAnchorEl(null);
    onEdit(piped);
  }, [piped, onEdit]);

  const handleAddNewKey = useCallback(() => {
    setAnchorEl(null);
    if (hasOldKey) {
      setOpenOldKeyAlert(true);
    } else {
      setOpenConfirmAddKey(true);
    }
  }, [hasOldKey]);

  const handleConfirmAddKey = useCallback(() => {
    setOpenConfirmAddKey(false);
    onAddNewKey(piped);
  }, [onAddNewKey, piped]);

  const handleCancelAddKey = useCallback(() => {
    setOpenConfirmAddKey(false);
  }, []);

  const handleDeleteOldKey = useCallback(() => {
    setAnchorEl(null);
    deleteOldPipedKey({ pipedId: piped.id }).then(() => {
      addToast({
        message: DELETE_OLD_PIPED_KEY_SUCCESS,
        severity: "success",
      });
    });
  }, [addToast, deleteOldPipedKey, piped.id]);

  const handleOpenPipedConfig = useCallback(() => {
    setAnchorEl(null);
    setOpenConfigAlert(true);
  }, []);

  const handleEnable = useCallback(() => {
    setAnchorEl(null);
    onEnable(piped.id);
  }, [piped.id, onEnable]);

  const handleDisable = useCallback(() => {
    setAnchorEl(null);
    onDisable(piped.id);
  }, [piped.id, onDisable]);

  const handleRestart = useCallback(() => {
    setAnchorEl(null);
    onRestart(piped.id);
  }, [piped.id, onRestart]);

  if (!piped) {
    return null;
  }

  const badgeColor = {
    [Piped.ConnectionStatus.ONLINE]: "green",
    [Piped.ConnectionStatus.OFFLINE]: "red",
    [Piped.ConnectionStatus.UNKNOWN]: "grey",
  };

  return (
    <>
      <TableRow
        key={`pipe-${piped.id}`}
        sx={{ bgcolor: piped.disabled ? "grey.200" : undefined }}
      >
        <TableCell>
          <Typography variant="subtitle2">
            {piped.name}
            <Tooltip
              placement="top"
              title={PIPED_CONNECTION_STATUS_TEXT[piped.status]}
            >
              <Badge
                variant="dot"
                overlap="rectangular"
                sx={{
                  paddingLeft: 1.5,
                  "& span": {
                    backgroundColor: badgeColor[piped.status],
                  },
                }}
              />
            </Tooltip>
          </Typography>
        </TableCell>
        <TableCell
          title={piped.id}
          sx={{
            "& button": { visibility: "hidden" },
            "&:hover button": { visibility: "visible" },
          }}
        >
          <Box
            sx={{
              display: "flex",
              alignItems: "center",
              fontFamily: "fontFamilyMono",
            }}
          >
            {piped.id}
            <CopyIconButton name="Piped ID" value={piped.id} />
          </Box>
        </TableCell>
        <TableCell>
          {piped.version}
          {piped.desiredVersion && piped.desiredVersion !== piped.version
            ? ` (upgrading to ${piped.desiredVersion})`
            : ""}
        </TableCell>
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
            size="large"
          >
            <MoreVertIcon />
          </IconButton>
        </TableCell>
      </TableRow>
      <Menu
        id="piped-menu"
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
        slotProps={{
          paper: {
            style: {
              maxHeight: ITEM_HEIGHT * 5.5,
              width: "25ch",
            },
          },
        }}
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
            <MenuItem
              key="piped-menu-open-piped-config"
              onClick={handleOpenPipedConfig}
              disabled={piped.config.length === 0}
            >
              {UI_TEXT_VIEW_THE_CONFIGURATION}
            </MenuItem>,
            <MenuItem
              key="piped-menu-restart"
              onClick={handleRestart}
              disabled={piped.status !== Piped.ConnectionStatus.ONLINE}
            >
              {UI_TEXT_RESTART}
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
      <DialogConfirm
        open={openConfirmAddKey}
        onCancel={handleCancelAddKey}
        onConfirm={handleConfirmAddKey}
        title="Add new piped key"
        description={`This piped has one key. Are you sure you want to generate a new key?\nAfter adding a new key and selecting 'Delete old key', the existing key will no longer be valid.`}
        confirmText={UI_TEXT_ADD_NEW_KEY}
        cancelText={UI_TEXT_CANCEL}
      />
      <Dialog
        fullWidth
        maxWidth="md"
        open={openConfigAlert}
        onClose={handleAlertClose}
      >
        <DialogTitle>
          <CopyIconButton name="Piped config" value={`${piped.config}`} />
          Piped configuration
        </DialogTitle>
        <DialogContent>
          <Highlight theme={themes.github} code={piped.config} language="yaml">
            {({ style, tokens, getLineProps, getTokenProps }) => (
              <Box
                component={"pre"}
                style={style}
                sx={{
                  padding: 2,
                  overflow: "auto",
                }}
              >
                {tokens.map((line, i) => (
                  <div key={i} {...getLineProps({ line })}>
                    {line.map((token, key) => (
                      <span key={key} {...getTokenProps({ token })} />
                    ))}
                  </div>
                ))}
              </Box>
            )}
          </Highlight>
        </DialogContent>
      </Dialog>
    </>
  );
});
