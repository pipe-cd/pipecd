import {
  Box,
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  Divider,
  makeStyles,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
} from "@material-ui/core";
import {
  Add as AddIcon,
  Close as CloseIcon,
  FilterList as FilterIcon,
  Update as UpgradeIcon,
} from "@material-ui/icons";
import Alert from "@material-ui/lab/Alert";
import { createSelector } from "@reduxjs/toolkit";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import {
  UI_TEXT_ADD,
  UI_TEXT_CLOSE,
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_UPGRADE,
} from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import {
  clearRegisteredPipedInfo,
  disablePiped,
  enablePiped,
  fetchPipeds,
  fetchReleasedVersions,
  Piped,
  RegisteredPiped,
  selectAllPipeds,
} from "~/modules/pipeds";
import { AppState } from "~/store";
import { useSettingsStyles } from "../styles";
import { AddPipedDrawer } from "./components/add-piped-drawer";
import { EditPipedDrawer } from "./components/edit-piped-drawer";
import { FilterValues, PipedFilter } from "./components/piped-filter";
import { PipedTableRow } from "./components/piped-table-row";
import { UpgradePipedDialog } from "./components/upgrade-dialog";

const useStyles = makeStyles(() => ({
  toolbarSpacer: {
    flexGrow: 1,
  },
}));

const selectFilteredPipeds = createSelector<
  AppState,
  boolean | undefined,
  Piped.AsObject[],
  boolean | undefined,
  Piped.AsObject[]
>(
  selectAllPipeds,
  (_, enabled) => enabled,
  (pipeds, enabled) => {
    switch (enabled) {
      case true:
        return pipeds.filter((piped) => piped.disabled === false);
      case false:
        return pipeds.filter((piped) => piped.disabled);
      default:
        return pipeds;
    }
  }
);

const OLD_KEY_ALERT_MESSAGE =
  "The old key is still there.\nDo not forget to delete it once you update your Piped to use this new key.";

const FETCH_INTERVAL = 30000;

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const classes = useStyles();
  const settingsClasses = useSettingsStyles();
  const [openFilter, setOpenFilter] = useState(false);
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [editPipedId, setEditPipedId] = useState<string | null>(null);
  const [filterValues, setFilterValues] = useState<FilterValues>({
    enabled: true,
  });
  const dispatch = useAppDispatch();
  const pipeds = useAppSelector((state) =>
    selectFilteredPipeds(state, filterValues.enabled)
  );

  useEffect(() => {
    dispatch(fetchReleasedVersions());
  }, [dispatch]);

  const releasedVersions = useAppSelector<string[]>(
    (state) => state.pipeds.releasedVersions
  );

  const [isUpgradeDialogOpen, setUpgradeDialogOpen] = useState(false);
  const handleUpgradeDialogClose = (): void => setUpgradeDialogOpen(false);

  const registeredPiped = useAppSelector<RegisteredPiped | null>(
    (state) => state.pipeds.registeredPiped
  );

  const handleDisable = useCallback(
    (id: string) => {
      dispatch(disablePiped({ pipedId: id })).then(() => {
        dispatch(fetchPipeds(true));
      });
    },
    [dispatch]
  );
  const handleEnable = useCallback(
    (id: string) => {
      dispatch(enablePiped({ pipedId: id })).then(() => {
        dispatch(fetchPipeds(true));
      });
    },
    [dispatch]
  );

  const handleEdit = useCallback((id: string) => {
    setEditPipedId(id);
  }, []);

  const handleClose = useCallback(() => {
    setIsOpenForm(false);
  }, []);

  const handleClosePipedInfo = useCallback(() => {
    dispatch(clearRegisteredPipedInfo());
    dispatch(fetchPipeds(true));
  }, [dispatch]);

  const handleEditClose = useCallback(() => {
    setEditPipedId(null);
  }, []);

  useInterval(() => {
    dispatch(fetchPipeds(true));
  }, FETCH_INTERVAL);

  return (
    <>
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
        <div className={classes.toolbarSpacer} />
        <Button
          color="primary"
          startIcon={<UpgradeIcon />}
          onClick={() => setUpgradeDialogOpen(true)}
        >
          {UI_TEXT_UPGRADE}
        </Button>
        <Button
          color="primary"
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />

      <Box display="flex" flex={1} overflow="hidden">
        <TableContainer component={Paper} square>
          <Table aria-label="piped list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell className={settingsClasses.tableCell}>
                  Name
                </TableCell>
                <TableCell className={settingsClasses.tableCell}>ID</TableCell>
                <TableCell className={settingsClasses.tableCell}>
                  Version
                </TableCell>
                <TableCell className={settingsClasses.tableCell}>
                  Description
                </TableCell>
                <TableCell className={settingsClasses.tableCell}>
                  Started At
                </TableCell>
                <TableCell align="right" />
              </TableRow>
            </TableHead>
            <TableBody>
              {pipeds.map((piped) => (
                <PipedTableRow
                  key={piped.id}
                  pipedId={piped.id}
                  onEdit={handleEdit}
                  onDisable={handleDisable}
                  onEnable={handleEnable}
                />
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {openFilter && (
          <PipedFilter values={filterValues} onChange={setFilterValues} />
        )}
      </Box>

      <AddPipedDrawer open={isOpenForm} onClose={handleClose} />
      <EditPipedDrawer pipedId={editPipedId} onClose={handleEditClose} />
      <UpgradePipedDialog
        open={isUpgradeDialogOpen}
        pipeds={pipeds}
        releasedVersions={releasedVersions}
        onClose={handleUpgradeDialogClose}
      />

      <Dialog fullWidth open={Boolean(registeredPiped)}>
        <DialogTitle>
          {registeredPiped?.isNewKey
            ? "Added a new piped key"
            : "Piped registered"}
        </DialogTitle>
        {registeredPiped?.isNewKey ? (
          <Alert severity="info">{OLD_KEY_ALERT_MESSAGE}</Alert>
        ) : null}
        <DialogContent>
          <TextWithCopyButton
            name="Piped Id"
            value={registeredPiped?.id || ""}
          />
          <TextWithCopyButton
            name="Piped Key"
            value={registeredPiped?.key || ""}
          />
          <Box display="flex" justifyContent="flex-end" m={1} mt={2}>
            <Button color="primary" onClick={handleClosePipedInfo}>
              {UI_TEXT_CLOSE}
            </Button>
          </Box>
        </DialogContent>
      </Dialog>
    </>
  );
});
