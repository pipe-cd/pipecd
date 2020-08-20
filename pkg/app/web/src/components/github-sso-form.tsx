import {
  Button,
  TextField,
  Typography,
  Dialog,
  DialogContent,
  DialogTitle,
  DialogActions,
} from "@material-ui/core";
import React, { FC, useState } from "react";

export interface GitHubSSOFormParams {
  clientId: string;
  clientSecret: string;
  baseUrl: string;
  uploadUrl: string;
  org: string;
  adminTeam: string;
  editorTeam: string;
  viewerTeam: string;
}

interface Props {
  isSaving: boolean;
  onSave: (params: GitHubSSOFormParams) => Promise<unknown>;
}

function hasEmptyValue(obj: Record<string, string>): boolean {
  return Object.keys(obj).some((key) => obj[key] === "");
}

const initialParams = {
  clientId: "",
  clientSecret: "",
  baseUrl: "",
  uploadUrl: "",
  org: "",
  adminTeam: "",
  editorTeam: "",
  viewerTeam: "",
};

export const GithubSSOForm: FC<Props> = ({ isSaving, onSave }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [params, setParams] = useState(initialParams);

  const hasEmpty = hasEmptyValue(params);

  const handleOnSave = (): void => {
    onSave(params).then(() => {
      setIsOpen(false);
      setParams(initialParams);
    });
  };

  const handleCancel = (): void => {
    setIsOpen(false);
    setParams(initialParams);
  };

  return (
    <div>
      <Typography variant="h6">GitHub</Typography>
      <Button variant="contained" onClick={() => setIsOpen(true)}>
        EDIT
      </Button>

      <Dialog open={isOpen}>
        <DialogTitle>GitHub Single Sign On Setting</DialogTitle>
        <DialogContent>
          <TextField
            label="Client ID"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.clientId}
            onChange={(e) => setParams({ ...params, clientId: e.target.value })}
          />
          <TextField
            label="Client Secret"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.clientSecret}
            onChange={(e) =>
              setParams({ ...params, clientSecret: e.target.value })
            }
          />
          <TextField
            label="Base URL"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.baseUrl}
            onChange={(e) => setParams({ ...params, baseUrl: e.target.value })}
          />
          <TextField
            label="Upload URL"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.uploadUrl}
            onChange={(e) =>
              setParams({ ...params, uploadUrl: e.target.value })
            }
          />
          <TextField
            label="Organization"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.org}
            onChange={(e) => setParams({ ...params, org: e.target.value })}
          />
          <TextField
            label="Admin Team"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.adminTeam}
            onChange={(e) =>
              setParams({ ...params, adminTeam: e.target.value })
            }
          />
          <TextField
            label="Editor Team"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.editorTeam}
            onChange={(e) =>
              setParams({ ...params, editorTeam: e.target.value })
            }
          />
          <TextField
            label="Viewer Team"
            margin="dense"
            variant="outlined"
            fullWidth
            value={params.viewerTeam}
            onChange={(e) =>
              setParams({ ...params, viewerTeam: e.target.value })
            }
          />
          <DialogActions>
            <Button color="primary" onClick={handleCancel} disabled={isSaving}>
              CANCEL
            </Button>
            <Button
              color="primary"
              onClick={handleOnSave}
              disabled={isSaving || hasEmpty}
            >
              SAVE
            </Button>
          </DialogActions>
        </DialogContent>
      </Dialog>
    </div>
  );
};
