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

type Props = {
  notes?: string | null;
};

const BreakingChangeNotes: FC<Props> = ({ notes }) => {
  const [showDialog, setShowDialog] = useState(false);

  if (!notes) {
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
          <Button onClick={() => setShowDialog(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default BreakingChangeNotes;
