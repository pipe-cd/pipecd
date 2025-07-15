import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
} from "@mui/material";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { ApplicationListItem } from "./application-list-item";
import { DeleteApplicationDialog } from "./delete-application-dialog";
import { DisableApplicationDialog } from "./disable-application-dialog";
import { SealedSecretDialog } from "./sealed-secret-dialog";
import { useEnableApplication } from "~/queries/applications/use-enable-application";
import { Application } from "~/types/applications";
import EditApplicationDrawer from "../edit-application-drawer";

const PAGER_ROWS_PER_PAGE = [20, 50, { label: "All", value: -1 }];
const SMALL_SCREEN_SIZE = 1440;

export interface ApplicationListProps {
  applications: Application.AsObject[];
  currentPage: number;
  onPageChange?: (page: number) => void;
  onRefresh?: () => void;
}

export const ApplicationList: FC<ApplicationListProps> = memo(
  function ApplicationList({
    applications,
    currentPage,
    onPageChange = () => null,
  }) {
    const [
      actionTarget,
      setActionTarget,
    ] = useState<Application.AsObject | null>(null);
    const [dialogState, setDialogState] = useState({
      edit: false,
      disabling: false,
      generateSecret: false,
      delete: false,
    });
    const [rowsPerPage, setRowsPerPage] = useState(20);
    const page = currentPage - 1;

    const { mutate: enableApplication } = useEnableApplication();

    const closeMenu = useCallback(() => {
      setActionTarget(null);
    }, []);

    const handleOnCloseGenerateDialog = (): void => {
      closeMenu();
      setDialogState({
        ...dialogState,
        generateSecret: false,
      });
    };

    const handleCloseDialog = (): void => {
      closeMenu();
      setDialogState({
        edit: false,
        delete: false,
        generateSecret: false,
        disabling: false,
      });
    };

    const handleEditClick = useCallback((app: Application.AsObject) => {
      setActionTarget(app);
      setDialogState((p) => ({
        ...p,
        edit: true,
      }));
    }, []);

    const handleDisableClick = useCallback(
      (app: Application.AsObject) => {
        setActionTarget(app);
        setDialogState((p) => ({
          ...p,
          disabling: true,
        }));
      },
      [setActionTarget, setDialogState]
    );

    const handleEnableClick = useCallback(
      (app: Application.AsObject) => {
        enableApplication(
          { applicationId: app.id },
          { onSuccess: () => closeMenu() }
        );
      },
      [enableApplication, closeMenu]
    );

    const handleDeleteClick = useCallback((app: Application.AsObject) => {
      setActionTarget(app);
      setDialogState((p) => ({
        ...p,
        delete: true,
      }));
    }, []);

    const handleEncryptSecretClick = useCallback(
      (app: Application.AsObject) => {
        setActionTarget(app);
        setDialogState((p) => ({
          ...p,
          generateSecret: true,
        }));
      },
      []
    );

    const [isSmallScreen, setIsSmallScreen] = useState(
      window.innerWidth <= SMALL_SCREEN_SIZE
    );

    const checkSmallScreen = (): void => {
      setIsSmallScreen(window.innerWidth <= SMALL_SCREEN_SIZE);
    };

    useEffect(() => {
      window.addEventListener("resize", checkSmallScreen);
      return () => window.removeEventListener("resize", checkSmallScreen);
    });

    return (
      <>
        <TableContainer component={Paper} sx={{ flex: 1 }} square>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Status</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Kind</TableCell>
                <TableCell>Labels</TableCell>
                <TableCell>Artifact Versions</TableCell>
                {!isSmallScreen && <TableCell>Running Commit</TableCell>}
                {!isSmallScreen && <TableCell>Deployed By</TableCell>}
                <TableCell>Updated At</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>
              {(rowsPerPage > 0
                ? applications.slice(
                    page * rowsPerPage,
                    page * rowsPerPage + rowsPerPage
                  )
                : applications
              ).map((app) => (
                <ApplicationListItem
                  key={`app-${app.id}`}
                  app={app}
                  displayAllProperties={!isSmallScreen}
                  onEdit={() => handleEditClick(app)}
                  onDisable={() => handleDisableClick(app)}
                  onEnable={() => handleEnableClick(app)}
                  onDelete={() => handleDeleteClick(app)}
                  onEncryptSecret={() => handleEncryptSecretClick(app)}
                />
              ))}
            </TableBody>
            <TableFooter>
              <TableRow>
                <TablePagination
                  rowsPerPageOptions={PAGER_ROWS_PER_PAGE}
                  count={applications.length}
                  rowsPerPage={rowsPerPage}
                  page={page}
                  colSpan={9}
                  onPageChange={(_, newPage) => {
                    onPageChange(newPage + 1);
                  }}
                  onRowsPerPageChange={(e) => {
                    setRowsPerPage(parseInt(e.target.value, 10));
                    onPageChange(1);
                  }}
                />
              </TableRow>
            </TableFooter>
          </Table>
        </TableContainer>

        <DisableApplicationDialog
          open={Boolean(actionTarget) && dialogState.disabling}
          application={actionTarget}
          onDisable={handleCloseDialog}
          onCancel={handleCloseDialog}
        />

        <SealedSecretDialog
          open={Boolean(actionTarget) && dialogState.generateSecret}
          application={actionTarget}
          onClose={handleOnCloseGenerateDialog}
        />

        <DeleteApplicationDialog
          open={Boolean(actionTarget) && dialogState.delete}
          application={actionTarget}
          onDeleted={handleCloseDialog}
          onCancel={handleCloseDialog}
        />

        <EditApplicationDrawer
          onUpdated={handleCloseDialog}
          open={Boolean(actionTarget) && dialogState.edit}
          application={actionTarget || undefined}
          onClose={handleCloseDialog}
        />
      </>
    );
  }
);
