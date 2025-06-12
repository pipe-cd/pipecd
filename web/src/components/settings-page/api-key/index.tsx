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
} from "@mui/material";
import { Add as AddIcon, MoreVert as MenuIcon } from "@mui/icons-material";
import Skeleton from "@mui/material/Skeleton";
import dayjs from "dayjs";
import * as React from "react";
import { FC, memo, useCallback, useState } from "react";
import { API_KEY_ROLE_TEXT } from "~/constants/api-key-role-text";
import {
  DISABLE_API_KEY_SUCCESS,
  GENERATE_API_KEY_SUCCESS,
} from "~/constants/toast-text";
import { APIKey } from "pipecd/web/model/apikey_pb";
import { DisableAPIKeyConfirmDialog } from "./components/disable-api-key-confirm-dialog";
import { GenerateAPIKeyDialog } from "./components/generate-api-key-dialog";
import { GeneratedAPIKeyDialog } from "./components/generated-api-key-dialog";
import { useGenerateApiKey } from "~/queries/api-keys/use-generate-api-key";
import { useDisableApiKey } from "~/queries/api-keys/use-disable-api-key";
import { useGetApiKeys } from "~/queries/api-keys/use-get-api-keys";
import { useToast } from "~/contexts/toast-context";

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
        <Box
          sx={{
            height: 48,
            width: 48,
          }}
        />
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
        <Box
          sx={{
            height: 48,
            width: 48,
          }}
        />
      </TableCell>
    </TableRow>
  );
});

export const APIKeyPage: FC = memo(function APIKeyPage() {
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [disableApiKey, setDisableApiKey] = useState<APIKey.AsObject | null>(
    null
  );
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );
  const [generatedKey, setGeneratedKey] = useState<string | null>(null);

  const { data: keys = [], isLoading: loading } = useGetApiKeys(
    { enabled: true },
    { retry: false }
  );
  const { addToast } = useToast();

  const { mutateAsync: generateApiKey } = useGenerateApiKey();
  const { mutateAsync: disableAPIKey } = useDisableApiKey();

  const unixTimeToString = (unixTime: number): string => {
    const dateTime = new Date(unixTime * 1000);
    return dateTime.toString();
  };

  const handleGenerateKey = useCallback(
    (values: { name: string; role: APIKey.Role }) => {
      generateApiKey(values)
        .then((result) => {
          setGeneratedKey(result.key);
          addToast({ message: GENERATE_API_KEY_SUCCESS, severity: "success" });
        })
        .catch(() => undefined);
    },
    [addToast, generateApiKey]
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
    setDisableApiKey(null);
  }, [setDisableApiKey]);

  const handleDisable = useCallback(
    (id: string) => {
      disableAPIKey({ id }).then(() => {
        addToast({ message: DISABLE_API_KEY_SUCCESS });
      });
      setDisableApiKey(null);
    },
    [addToast, disableAPIKey]
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
                    <IconButton
                      data-id={key.id}
                      onClick={handleOpenMenu}
                      size="large"
                    >
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
        slotProps={{
          paper: {
            style: {
              width: "25ch",
            },
          },
        }}
      >
        <MenuItem
          onClick={() => {
            if (anchorEl && anchorEl.dataset.id) {
              const apiKey = keys.find((key) => key.id === anchorEl.dataset.id);
              if (apiKey) {
                setDisableApiKey(apiKey);
              }
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
        onSubmit={handleGenerateKey}
      />
      <GeneratedAPIKeyDialog
        generatedKey={generatedKey}
        onClose={() => setGeneratedKey(null)}
      />
      <DisableAPIKeyConfirmDialog
        apiKey={disableApiKey}
        onCancel={handleCancelDisabling}
        onDisable={handleDisable}
      />
    </>
  );
});
