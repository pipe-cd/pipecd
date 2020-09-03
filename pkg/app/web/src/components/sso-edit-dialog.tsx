import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from "@material-ui/core";
import React, { FC, useState } from "react";
import { GitHubSSO } from "../modules/project";

interface Props {
  open: boolean;
  onSave: (
    params: Partial<GitHubSSO> & { clientId: string; clientSecret: string }
  ) => void;
  onClose: () => void;
  currentBaseURL: string;
  currentUploadURL: string;
}

export const SSOEditDialog: FC<Props> = ({
  open,
  currentBaseURL,
  currentUploadURL,
  onSave,
  onClose,
}) => {
  const [clientId, setClientID] = useState("");
  const [clientSecret, setClientSecret] = useState("");
  const [baseUrl, setBaseUrl] = useState(currentBaseURL);
  const [uploadUrl, setUploadUrl] = useState(currentUploadURL);

  const handleSave = (): void => {
    onSave({
      clientId,
      clientSecret,
      baseUrl,
      uploadUrl,
    });
    onClose();
  };

  const isValid = Boolean(clientId) && Boolean(clientSecret);

  return (
    <Dialog
      open={open}
      onEnter={() => {
        setBaseUrl(currentBaseURL);
        setUploadUrl(currentUploadURL);
      }}
      onClose={onClose}
    >
      <DialogTitle>Edit GitHub SSO</DialogTitle>
      <DialogContent>
        <TextField
          value={clientId}
          variant="outlined"
          margin="dense"
          label="Client ID"
          fullWidth
          required
          onChange={(e) => setClientID(e.currentTarget.value)}
        />
        <TextField
          value={clientSecret}
          variant="outlined"
          margin="dense"
          label="Client Secret"
          fullWidth
          required
          onChange={(e) => setClientSecret(e.currentTarget.value)}
        />
        <TextField
          value={baseUrl}
          variant="outlined"
          margin="dense"
          label="Base URL"
          fullWidth
          onChange={(e) => setBaseUrl(e.currentTarget.value)}
        />
        <TextField
          value={uploadUrl}
          variant="outlined"
          margin="dense"
          label="Upload URL"
          fullWidth
          onChange={(e) => setUploadUrl(e.currentTarget.value)}
        />
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>CANCEL</Button>
        <Button
          onClick={handleSave}
          type="submit"
          color="primary"
          disabled={isValid === false}
        >
          SAVE
        </Button>
      </DialogActions>
    </Dialog>
  );
};
