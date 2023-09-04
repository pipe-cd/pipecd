import {
  Box,
  Button,
  Divider,
  IconButton,
  Menu,
  MenuItem,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
  Tooltip,
  Typography,
} from "@material-ui/core";
import { Add as AddIcon, MoreVert as MenuIcon } from "@material-ui/icons";
import Skeleton from "@material-ui/lab/Skeleton";
import dayjs from "dayjs";
import * as React from "react";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { API_KEY_ROLE_TEXT } from "~/constants/api-key-role-text";
import {
  DISABLE_API_KEY_SUCCESS,
  GENERATE_API_KEY_SUCCESS,
} from "~/constants/toast-text";
import { unwrapResult, useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  APIKey,
  disableAPIKey,
  fetchAPIKeys,
  generateAPIKey,
  selectAll as selectAPIKeys,
} from "~/modules/api-keys";
import { addToast } from "~/modules/toasts";
import { DisableAPIKeyConfirmDialog } from "./components/disable-api-key-confirm-dialog";
import { GenerateAPIKeyDialog } from "./components/generate-api-key-dialog";
import { GeneratedAPIKeyDialog } from "./components/generated-api-key-dialog";

const LoadingSkelton = memo(function LoadingSkelton() {
  return (
    <TableRow>
      <TableCell colSpan={2}>
        <Skeleton width={200} height={30} />
      </TableCell>
      <TableCell>
        <Skeleton width={200} height={30} />
      </TableCell>
      <TableCell align="right">
        <Box height={48} width={48} />
      </TableCell>
    </TableRow>
  );
});

const EmptyTableContent = memo(function EmptyTableContent() {
  return (
    <TableRow>
      <TableCell colSpan={3}>
        <Typography>No API Keys</Typography>
      </TableCell>
      <TableCell align="right">
        <Box height={48} width={48} />
      </TableCell>
    </TableRow>
  );
});

const menuStyle = {
  style: {
    width: "25ch",
  },
};

export const APIKeyPage: FC = memo(function APIKeyPage() {
  const dispatch = useAppDispatch();
  const [loading, keys] = useAppSelector<[boolean, APIKey.AsObject[]]>(
    (state) => [state.apiKeys.loading, selectAPIKeys(state.apiKeys)]
  );
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [disableTargetId, setDisableTargetId] = useState<null | string>(null);
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );

  useEffect(() => {
    dispatch(fetchAPIKeys({ enabled: true }));
  }, [dispatch]);

  const unixTimeToString = (unixTime: number): string => {
    const dateTime = new Date(unixTime * 1000);
    return dateTime.toString();
  };

  const handleSubmit = useCallback(
    (values: { name: string; role: APIKey.Role }) => {
      dispatch(generateAPIKey(values))
        .then(unwrapResult)
        .then(() => {
          console.log("handleSubmit.then");
          dispatch(fetchAPIKeys({ enabled: true }));
          dispatch(
            addToast({ message: GENERATE_API_KEY_SUCCESS, severity: "success" })
          );
        })
        .catch(() => {});
    },
    [dispatch]
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

  const handleCancelDisabling = useCallback(() => {
    setDisableTargetId(null);
  }, [setDisableTargetId]);

  const handleDisable = useCallback(
    (id: string) => {
      dispatch(disableAPIKey({ id })).then(() => {
        dispatch(fetchAPIKeys({ enabled: true }));
        dispatch(addToast({ message: DISABLE_API_KEY_SUCCESS }));
      });
      setDisableTargetId(null);
    },
    [dispatch, setDisableTargetId]
  );

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenAddForm(true)}
        >
          ADD
        </Button>
      </Toolbar>
      <Divider />

      <TableContainer component={Paper} square>
        <Table size="small" stickyHeader>
          <TableHead>
            <TableRow>
              <TableCell colSpan={2}>Name</TableCell>
              <TableCell>Role</TableCell>
              <TableCell>CreatedAt</TableCell>
              <TableCell>LastUsedAt</TableCell>
              <TableCell align="right" />
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <LoadingSkelton />
            ) : keys.length === 0 ? (
              <EmptyTableContent />
            ) : (
              keys.map((key) => (
                <TableRow key={key.id}>
                  <TableCell colSpan={2}>{key.name}</TableCell>
                  <TableCell>{API_KEY_ROLE_TEXT[key.role]}</TableCell>
                  <TableCell>
                    <Tooltip
                      placement="top-start"
                      title={unixTimeToString(key.createdAt)}
                    >
                      <span>{dayjs(key.createdAt * 1000).fromNow()}</span>
                    </Tooltip>
                  </TableCell>
                  <TableCell>
                    {key.lastUsedAt === 0 ? (
                      "never used"
                    ) : (
                      <Tooltip
                        placement="top-start"
                        title={unixTimeToString(key.lastUsedAt)}
                      >
                        <span>{dayjs(key.lastUsedAt * 1000).fromNow()}</span>
                      </Tooltip>
                    )}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton data-id={key.id} onClick={handleOpenMenu}>
                      <MenuIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>

      <Menu
        id="api-key-menu"
        open={Boolean(anchorEl)}
        anchorEl={anchorEl}
        onClose={handleCloseMenu}
        PaperProps={menuStyle}
      >
        <MenuItem
          onClick={() => {
            if (anchorEl && anchorEl.dataset.id) {
              setDisableTargetId(anchorEl.dataset.id);
            }
            setAnchorEl(null);
          }}
        >
          Disable API Key
        </MenuItem>
      </Menu>

      <GenerateAPIKeyDialog
        open={isOpenAddForm}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />

      <GeneratedAPIKeyDialog />

      <DisableAPIKeyConfirmDialog
        apiKeyId={disableTargetId}
        onCancel={handleCancelDisabling}
        onDisable={handleDisable}
      />
    </>
  );
});
