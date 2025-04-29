import {
  Box,
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Paper,
  Table,
  TableBody,
  TableContainer,
  TableCell,
  TableHead,
  TableRow,
  Typography,
  TextField,
} from "@mui/material";
import { FC, memo, useCallback, useState, FormEvent } from "react";
import { UPGRADE_PIPEDS_SUCCESS } from "~/constants/toast-text";
import { useAppDispatch } from "~/hooks/redux";
import { UI_TEXT_CANCEL, UI_TEXT_UPGRADE } from "~/constants/ui-text";
import { Piped, updatePipedDesiredVersion } from "~/modules/pipeds";
import { addToast } from "~/modules/toasts";
import { useSettingsStyles } from "../../../styles";
import { Autocomplete } from "@mui/material";

export interface UpgradePipedProps {
  open: boolean;
  pipeds: Piped.AsObject[];
  releasedVersions: string[];
  onClose: () => void;
}

export const UpgradePipedDialog: FC<UpgradePipedProps> = memo(
  function UpgradePipedDialog({ open, pipeds, releasedVersions, onClose }) {
    const settingsClasses = useSettingsStyles();
    const dispatch = useAppDispatch();

    const [upgradeVersion, setUpgradeVersion] = useState("");
    const [upgradePipedIds, setUpgradePipedIds] = useState<Array<string>>([]);

    // This function will be triggered when a checkbox changes its state.
    const selectUpgradePiped = (
      e: React.ChangeEvent<HTMLInputElement>
    ): void => {
      const selectedId = e.target.value;
      // Check if "upgradePipedIds" contains "selectedIds".
      if (upgradePipedIds.includes(selectedId)) {
        // If true, this checkbox is already checked and we have to remote it from the list.
        setUpgradePipedIds((ids) => ids.filter((id) => id !== selectedId));
      } else {
        // Otherwise, it is not selected yet and we have to add it into the list.
        setUpgradePipedIds((upgradePipedIds) => [
          ...upgradePipedIds,
          selectedId,
        ]);
      }
    };

    const handleClose = useCallback(() => {
      onClose();
      setUpgradeVersion("");
      setUpgradePipedIds([]);
    }, [onClose]);

    const handleSubmit = async (e: FormEvent): Promise<void> => {
      e.preventDefault();

      await dispatch(
        updatePipedDesiredVersion({
          version: upgradeVersion,
          pipedIds: upgradePipedIds,
        })
      ).then(() => {
        dispatch(
          addToast({ message: UPGRADE_PIPEDS_SUCCESS, severity: "success" })
        );
        onClose();
        setUpgradeVersion("");
        setUpgradePipedIds([]);
      });
    };

    return (
      <Dialog open={open} onClose={handleClose}>
        <form onSubmit={handleSubmit}>
          <DialogTitle>Upgrade pipeds to a new version</DialogTitle>
          <DialogContent>
            <Box mb={3}>
              <Typography>1. Input your desired version</Typography>
              <Autocomplete
                id="version"
                freeSolo
                autoSelect
                options={releasedVersions.slice(0, 6)}
                onChange={(_, value) => {
                  setUpgradeVersion(value || "");
                }}
                renderInput={(params) => (
                  <TextField {...params} variant="outlined" />
                )}
              />
            </Box>

            <Typography>2. Select pipeds to upgrade</Typography>
            <Box display="flex" flex={1} overflow="hidden" mt={1}>
              <TableContainer component={Paper} square>
                <Table aria-label="piped list" size="small" stickyHeader>
                  <TableHead>
                    <TableRow>
                      <TableCell className={settingsClasses.tableCell}>
                        Name
                      </TableCell>
                      <TableCell className={settingsClasses.tableCell}>
                        Running Version
                      </TableCell>
                      <TableCell align="right">Select</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {pipeds.map((piped) => (
                      <TableRow key={piped.id}>
                        <TableCell>
                          <Typography variant="subtitle2">
                            {piped.name}
                          </Typography>
                        </TableCell>
                        <TableCell>{piped.version}</TableCell>
                        <TableCell>
                          <Checkbox
                            checked={
                              upgradePipedIds.includes(piped.id) ? true : false
                            }
                            value={piped.id}
                            onChange={selectUpgradePiped}
                          />
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={onClose}>{UI_TEXT_CANCEL}</Button>
            <Button
              type="submit"
              color="primary"
              disabled={
                upgradePipedIds.length === 0 || upgradeVersion.length === 0
              }
            >
              {UI_TEXT_UPGRADE}
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    );
  }
);
