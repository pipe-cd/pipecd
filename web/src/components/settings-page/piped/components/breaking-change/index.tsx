import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FC, useState } from "react";
import { IGNORE_BREAKING_CHANGE_NOTES_PIPEDS } from "~/constants/localstorage";

type Props = {
  notes?: string | null;
};

const getVersionsIgnoredWarning = (): string[] => {
  try {
    const rawIgnoredNotes =
      localStorage.getItem(IGNORE_BREAKING_CHANGE_NOTES_PIPEDS) || "[]";
    return JSON.parse(rawIgnoredNotes) as string[];
  } catch {
    return [];
  }
};

const shouldIgnoredBreakingChangeNotes = (): boolean => {
  const version = process.env.PIPECD_VERSION;
  if (!version) return false;

  try {
    const ignoredNotes = getVersionsIgnoredWarning();
    return ignoredNotes.includes(version);
  } catch {
    return false;
  }
};

const BreakingChangeNotes: FC<Props> = ({ notes }) => {
  const [showDialog, setShowDialog] = useState(false);
  const [showNotes, setShowNotes] = useState(
    !shouldIgnoredBreakingChangeNotes()
  );

  const onIgnoreWarning = (): void => {
    setShowDialog(false);
    const pipedVersion = process.env.PIPECD_VERSION;
    if (!pipedVersion) return;

    try {
      const ignoredVersions = JSON.parse(
        localStorage.getItem(IGNORE_BREAKING_CHANGE_NOTES_PIPEDS) || "[]"
      );

      if (!ignoredVersions.includes(pipedVersion)) {
        ignoredVersions.push(pipedVersion);
      }

      localStorage.setItem(
        IGNORE_BREAKING_CHANGE_NOTES_PIPEDS,
        JSON.stringify(ignoredVersions)
      );
    } finally {
      setShowNotes(false);
    }
  };

  if (!notes || !showNotes) {
    return null;
  }
  return (
    <>
      <Alert
        severity="warning"
        sx={{ alignItems: "center" }}
        action={
          <Button onClick={() => setShowDialog(true)}>View details</Button>
        }
      >
        <Typography
          sx={{
            overflow: "hidden",
            textOverflow: "ellipsis",
            display: "-webkit-box",
            WebkitLineClamp: "2",
            WebkitBoxOrient: "vertical",
          }}
        >
          {notes}
        </Typography>
      </Alert>

      <Dialog open={showDialog} onClose={() => setShowDialog(false)}>
        <DialogTitle>Breaking Changes</DialogTitle>
        <DialogContent
          sx={{
            whiteSpace: "pre-wrap",
            maxHeight: "60vh",
            overflowY: "auto",
          }}
        >
          {notes}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => onIgnoreWarning()}>Ignore</Button>
          <Button onClick={() => setShowDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default BreakingChangeNotes;
