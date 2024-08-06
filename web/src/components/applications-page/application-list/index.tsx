import {
  makeStyles,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
} from "@material-ui/core";
import { FC, memo, useCallback, useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  Application,
  enableApplication,
  selectAll,
} from "~/modules/applications";
import { setDeletingAppId } from "~/modules/delete-application";
import { setUpdateTargetId } from "~/modules/update-application";
import { ApplicationListItem } from "./application-list-item";
import { DeleteApplicationDialog } from "./delete-application-dialog";
import { DisableApplicationDialog } from "./disable-application-dialog";
import { SealedSecretDialog } from "./sealed-secret-dialog";

const useStyles = makeStyles(() => ({
  container: {
    flex: 1,
  },
  tooltip: {
    paddingLeft: 2,
    marginBottom: -4,
  },
}));

const PAGER_ROWS_PER_PAGE = [20, 50, { label: "All", value: -1 }];
const SMALL_SCREEN_SIZE = 1440;

export interface ApplicationListProps {
  currentPage: number;
  onPageChange?: (page: number) => void;
  onRefresh?: () => void;
}

export const ApplicationList: FC<ApplicationListProps> = memo(
  function ApplicationList({
    currentPage,
    onPageChange = () => null,
    onRefresh = () => null,
  }) {
    const classes = useStyles();
    const dispatch = useAppDispatch();
    const [actionTarget, setActionTarget] = useState<string | null>(null);
    const [dialogState, setDialogState] = useState({
      disabling: false,
      generateSecret: false,
    });
    const [rowsPerPage, setRowsPerPage] = useState(20);
    const page = currentPage - 1;

    const applications = useAppSelector<Application.AsObject[]>((state) =>
      selectAll(state.applications)
    );

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
        ...dialogState,
        disabling: false,
      });
      onRefresh();
    };

    // Menu item event handler

    const handleEditClick = useCallback(
      (id: string) => {
        closeMenu();
        dispatch(setUpdateTargetId(id));
      },
      [dispatch, closeMenu]
    );

    const handleDisableClick = useCallback(
      (id: string) => {
        setActionTarget(id);
        setDialogState({
          ...dialogState,
          disabling: true,
        });
      },
      [dialogState]
    );

    const handleEnableClick = useCallback(
      async (id: string) => {
        await dispatch(enableApplication({ applicationId: id }));
        onRefresh();
        closeMenu();
      },
      [dispatch, closeMenu, onRefresh]
    );

    const handleDeleteClick = useCallback(
      (id: string) => {
        dispatch(setDeletingAppId(id));
        closeMenu();
      },
      [dispatch, closeMenu]
    );

    const handleEncryptSecretClick = useCallback(
      (id: string) => {
        setActionTarget(id);
        setDialogState({
          ...dialogState,
          generateSecret: true,
        });
      },
      [dialogState]
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
        <TableContainer component={Paper} className={classes.container} square>
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
                  applicationId={app.id}
                  displayAllProperties={!isSmallScreen}
                  onEdit={handleEditClick}
                  onDisable={handleDisableClick}
                  onEnable={handleEnableClick}
                  onDelete={handleDeleteClick}
                  onEncryptSecret={handleEncryptSecretClick}
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
          open={dialogState.disabling}
          applicationId={actionTarget}
          onDisable={handleCloseDialog}
          onCancel={handleCloseDialog}
        />

        <SealedSecretDialog
          open={Boolean(actionTarget) && dialogState.generateSecret}
          applicationId={actionTarget}
          onClose={handleOnCloseGenerateDialog}
        />

        <DeleteApplicationDialog onDeleted={onRefresh} />
      </>
    );
  }
);
