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
} from "@material-ui/icons";
import { createSelector } from "@reduxjs/toolkit";
import { FC, memo, useCallback, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AddPipedDrawer } from "../../../components/add-piped-drawer";
import { EditPipedDrawer } from "../../../components/edit-piped-drawer";
import { FilterValues, PipedFilter } from "../../../components/piped-filter";
import { TextWithCopyButton } from "../../../components/text-with-copy-button";
import {
  UI_TEXT_ADD,
  UI_TEXT_CLOSE,
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
} from "../../../constants/ui-text";
import { AppState } from "../../../modules";
import {
  clearRegisteredPipedInfo,
  disablePiped,
  enablePiped,
  fetchPipeds,
  Piped,
  recreatePipedKey,
  RegisteredPiped,
  selectAllPipeds,
} from "../../../modules/pipeds";
import { AppDispatch } from "../../../store";
import { PipedTableRow } from "./piped-table-row";

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

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const classes = useStyles();
  const [openFilter, setOpenFilter] = useState(false);
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [editPipedId, setEditPipedId] = useState<string | null>(null);
  const [filterValues, setFilterValues] = useState<FilterValues>({
    enabled: true,
  });
  const dispatch = useDispatch<AppDispatch>();
  const pipeds = useSelector((state: AppState) =>
    selectFilteredPipeds(state, filterValues.enabled)
  );

  const registeredPiped = useSelector<AppState, RegisteredPiped | null>(
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

  const handleRecreate = useCallback(
    (id: string) => {
      dispatch(recreatePipedKey({ pipedId: id }));
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
          startIcon={openFilter ? <CloseIcon /> : <FilterIcon />}
          onClick={() => setOpenFilter(!openFilter)}
        >
          {openFilter ? UI_TEXT_HIDE_FILTER : UI_TEXT_FILTER}
        </Button>
      </Toolbar>
      <Divider />

      <Box display="flex" height="100%">
        <TableContainer component={Paper} square>
          <Table aria-label="piped list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>ID</TableCell>
                <TableCell>Version</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Started At</TableCell>
                <TableCell align="right" />
              </TableRow>
            </TableHead>
            <TableBody>
              {pipeds.map((piped) => (
                <PipedTableRow
                  key={piped.id}
                  pipedId={piped.id}
                  onEdit={handleEdit}
                  onRecreateKey={handleRecreate}
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

      <Dialog fullWidth open={Boolean(registeredPiped)}>
        <DialogTitle>Piped registered</DialogTitle>
        <DialogContent>
          <TextWithCopyButton
            name="Piped Id"
            label="Copy piped id"
            value={registeredPiped?.id || ""}
          />
          <TextWithCopyButton
            name="Secret Key"
            label="Copy secret key"
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
