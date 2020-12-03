import {
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
  Typography,
} from "@material-ui/core";
import { Add as AddIcon, MoreVert as MenuIcon } from "@material-ui/icons";
import React, { FC, memo, useCallback, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { GenerateAPIKeyDialog } from "../../components/generate-api-key-dialog";
import { GeneratedApiKeyDialog } from "../../components/generated-api-key-dialog";
import { API_KEY_ROLE_TEXT } from "../../constants/api-key-role-text";
import { AppState } from "../../modules";
import {
  APIKey,
  APIKeyModel,
  generateAPIKey,
  fetchAPIKeys,
  clearGeneratedKey,
} from "../../modules/api-keys";

export const APIKeyPage: FC = memo(function APIKeyPage() {
  const dispatch = useDispatch();
  const keys = useSelector<AppState, APIKey[]>((state) => state.apiKeys.items);
  const generatedKey = useSelector<AppState, string | null>(
    (state) => state.apiKeys.generatedKey
  );
  const [isOpenAddForm, setIsOpenAddForm] = useState(false);
  const [anchorEl, setAnchorEl] = React.useState<HTMLButtonElement | null>(
    null
  );

  useEffect(() => {
    dispatch(fetchAPIKeys({ enabled: true }));
  }, [dispatch]);

  const handleSubmit = useCallback(
    (values: { name: string; role: APIKeyModel.Role }) => {
      dispatch(generateAPIKey(values));
    },
    [dispatch]
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

      <TableContainer component={Paper}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell colSpan={2}>Name</TableCell>
              <TableCell>Role</TableCell>
              <TableCell align="right" />
            </TableRow>
          </TableHead>
          <TableBody>
            {keys.length === 0 ? (
              <TableRow>
                <TableCell colSpan={2}>
                  <Typography>No API Keys</Typography>
                </TableCell>
              </TableRow>
            ) : (
              keys.map((key) => (
                <TableRow key={key.id}>
                  <TableCell colSpan={2}>{key.name}</TableCell>
                  <TableCell>{API_KEY_ROLE_TEXT[key.role]}</TableCell>
                  <TableCell align="right">
                    <IconButton
                      onClick={(e) => {
                        setAnchorEl(e.currentTarget);
                      }}
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
        onClose={() => setAnchorEl(null)}
      >
        <MenuItem>Disable API Key</MenuItem>
      </Menu>

      <GenerateAPIKeyDialog
        open={isOpenAddForm}
        onClose={() => setIsOpenAddForm(false)}
        onSubmit={handleSubmit}
      />

      <GeneratedApiKeyDialog
        open={Boolean(generatedKey)}
        generatedKey={generatedKey}
        onClose={() => {
          dispatch(clearGeneratedKey());
        }}
      />
    </>
  );
});
