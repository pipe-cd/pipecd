import {
  Box,
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  Divider,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
} from "@mui/material";
import {
  Add as AddIcon,
  Close as CloseIcon,
  FilterList as FilterIcon,
  Update as UpgradeIcon,
} from "@mui/icons-material";
import Alert from "@mui/material/Alert";
import { FC, memo, useCallback, useMemo, useState } from "react";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import {
  UI_TEXT_ADD,
  UI_TEXT_CLOSE,
  UI_TEXT_FILTER,
  UI_TEXT_HIDE_FILTER,
  UI_TEXT_UPGRADE,
} from "~/constants/ui-text";
import { REQUEST_PIPED_RESTART_SUCCESS } from "~/constants/toast-text";
import { Piped } from "pipecd/web/model/piped_pb";
import { AddPipedDialog } from "./components/add-piped-dialog";
import { EditPipedDialog } from "./components/edit-piped-dialog";
import { FilterValues, PipedFilter } from "./components/piped-filter";
import { PipedTableRow } from "./components/piped-table-row";
import { UpgradePipedDialog } from "./components/upgrade-dialog";
import { TableCellNoWrap } from "../styles";
import { useGetProject } from "~/queries/project/use-get-project";
import { useToast } from "~/contexts/toast-context";
import { useGetPipeds } from "~/queries/pipeds/use-get-pipeds";
import { useGetReleasedVersions } from "~/queries/pipeds/use-get-released-versions";
import { useGetBreakingChanges } from "~/queries/pipeds/use-get-breaking-changes";
import { useAddNewPipedKey } from "~/queries/pipeds/use-add-new-piped-key";
import { useDisablePiped } from "~/queries/pipeds/use-disable-piped";
import { useEnablePiped } from "~/queries/pipeds/use-enable-piped";
import { useRestartPiped } from "~/queries/pipeds/use-restart-piped";
import BreakingChangeNotes from "./components/breaking-change";

const OLD_KEY_ALERT_MESSAGE =
  "The old key is still there.\nDo not forget to delete it once you update your Piped to use this new key.";

const FETCH_INTERVAL = 30000;

export const SettingsPipedPage: FC = memo(function SettingsPipedPage() {
  const [openFilter, setOpenFilter] = useState(false);
  const [isOpenForm, setIsOpenForm] = useState(false);
  const [editPiped, setEditPiped] = useState<Piped.AsObject | null>(null);
  const [isUpgradeDialogOpen, setIsUpgradeDialogOpen] = useState(false);
  const [registeredPiped, setRegisteredPiped] = useState<{
    id: string;
    key: string;
    isNewKey: boolean;
  } | null>(null);
  const [filterValues, setFilterValues] = useState<FilterValues>({
    enabled: true,
  });

  const { addToast } = useToast();

  const { data: projectDetail } = useGetProject();

  const { mutateAsync: addNewPipedKey } = useAddNewPipedKey();

  const { mutateAsync: restartPiped } = useRestartPiped();

  const { mutateAsync: disablePiped } = useDisablePiped();

  const { mutateAsync: enablePiped } = useEnablePiped();

  const { data: allPipeds } = useGetPipeds(
    { withStatus: true },
    { refetchInterval: FETCH_INTERVAL }
  );

  const { data: releasedVersions = [] } = useGetReleasedVersions();

  const { data: breakingChangesNote } = useGetBreakingChanges(
    { projectId: projectDetail?.id ?? "" },
    { enabled: !!projectDetail?.id }
  );

  const pipeds = useMemo(() => {
    return (
      allPipeds?.filter((piped) => {
        if (filterValues.enabled === true) return !piped.disabled;
        if (filterValues.enabled === false) return piped.disabled;
        return true;
      }) ?? []
    );
  }, [allPipeds, filterValues.enabled]);

  // TODO: Remove this console.log
  console.log("[DEBUG]", breakingChangesNote);

  const handleUpgradeDialogClose = (): void => {
    setIsUpgradeDialogOpen(false);
  };

  const handleDisable = useCallback(
    (id: string) => {
      disablePiped({ pipedId: id });
    },
    [disablePiped]
  );

  const handleEnable = useCallback(
    (id: string) => {
      enablePiped({ pipedId: id });
    },
    [enablePiped]
  );

  const handleRestart = useCallback(
    (id: string) => {
      restartPiped({ pipedId: id }).then(() => {
        addToast({
          message: REQUEST_PIPED_RESTART_SUCCESS,
          severity: "success",
        });
      });
    },
    [addToast, restartPiped]
  );

  const handleEdit = useCallback((piped: Piped.AsObject) => {
    setEditPiped(piped);
  }, []);

  const handleClose = useCallback(() => {
    setIsOpenForm(false);
  }, []);

  const handleClosePipedInfo = useCallback(() => {
    setRegisteredPiped(null);
  }, []);

  const handleEditClose = useCallback(() => {
    setEditPiped(null);
  }, []);

  const handleSuccessAddedPiped = (data: { id: string; key: string }): void => {
    setIsOpenForm(false);
    setRegisteredPiped({ ...data, isNewKey: false });
  };

  const handleAddNewKey = useCallback(
    (piped: Piped.AsObject) => {
      addNewPipedKey({ pipedId: piped.id }).then((key) => {
        setRegisteredPiped({
          id: piped.id,
          key,
          isNewKey: true,
        });
      });
    },
    [addNewPipedKey]
  );

  return (
    <>
      <BreakingChangeNotes notes={breakingChangesNote} />
      <Toolbar variant="dense">
        <Button
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => setIsOpenForm(true)}
        >
          {UI_TEXT_ADD}
        </Button>
        <Box
          sx={{
            flexGrow: 1,
          }}
        />
        <Button
          color="primary"
          startIcon={<UpgradeIcon />}
          onClick={() => setIsUpgradeDialogOpen(true)}
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
      <Box
        sx={{
          display: "flex",
          flex: 1,
          overflow: "hidden",
        }}
      >
        <TableContainer component={Paper} square>
          <Table aria-label="piped list" size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCellNoWrap>Name</TableCellNoWrap>
                <TableCellNoWrap>ID</TableCellNoWrap>
                <TableCellNoWrap>Version</TableCellNoWrap>
                <TableCellNoWrap>Description</TableCellNoWrap>
                <TableCellNoWrap>Started At</TableCellNoWrap>
                <TableCell align="right" />
              </TableRow>
            </TableHead>
            <TableBody>
              {pipeds.map((piped) => (
                <PipedTableRow
                  key={piped.id}
                  piped={piped}
                  onAddNewKey={handleAddNewKey}
                  onEdit={handleEdit}
                  onDisable={handleDisable}
                  onEnable={handleEnable}
                  onRestart={handleRestart}
                />
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        {openFilter && (
          <PipedFilter values={filterValues} onChange={setFilterValues} />
        )}
      </Box>
      <AddPipedDialog
        open={isOpenForm}
        onSuccess={handleSuccessAddedPiped}
        onClose={handleClose}
      />

      <EditPipedDialog piped={editPiped} onClose={handleEditClose} />

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
            value={registeredPiped?.id ?? ""}
          />
          <TextWithCopyButton
            name="Piped Key"
            value={registeredPiped?.key ?? ""}
          />
          <TextWithCopyButton
            name="Base64 Encoded Piped Key"
            value={
              registeredPiped?.key !== undefined
                ? btoa(registeredPiped?.key)
                : ""
            }
          />
          <Box
            sx={{
              display: "flex",
              justifyContent: "flex-end",
              m: 1,
              mt: 2,
            }}
          >
            <Button color="primary" onClick={handleClosePipedInfo}>
              {UI_TEXT_CLOSE}
            </Button>
          </Box>
        </DialogContent>
      </Dialog>
    </>
  );
});
